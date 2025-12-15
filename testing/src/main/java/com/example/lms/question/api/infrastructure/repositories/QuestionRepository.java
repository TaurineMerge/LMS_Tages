package com.example.lms.question.api.infrastructure.repositories;

import com.example.lms.question.api.domain.model.QuestionModel;
import com.example.lms.question.api.domain.repository.QuestionRepositoryInterface;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.sql.DataSource;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Реализация репозитория вопросов с использованием JDBC.
 * Работает с таблицей question_d в PostgreSQL.
 *
 * Структура QUESTION_D:
 *  - id              UUID  (PK, not null)
 *  - test_id         UUID  (FK → test_d.id, not null)
 *  - text_of_question TEXT
 *  - "order"         INT   (может быть NULL)
 */
public class QuestionRepository implements QuestionRepositoryInterface {

    private static final Logger logger = LoggerFactory.getLogger(QuestionRepository.class);
    private final DataSource dataSource;

    // ---------------------- SQL ЗАПРОСЫ ----------------------

    private static final String INSERT_SQL = """
        INSERT INTO question_d (test_id, text_of_question, "order")
        VALUES (?, ?, ?)
        RETURNING id
        """;

    private static final String UPDATE_SQL = """
        UPDATE question_d
        SET test_id = ?, text_of_question = ?, "order" = ?
        WHERE id = ?
        """;

    private static final String SELECT_BY_ID = """
        SELECT id, test_id, text_of_question, "order"
        FROM question_d
        WHERE id = ?
        """;

    private static final String SELECT_ALL = """
        SELECT id, test_id, text_of_question, "order"
        FROM question_d
        ORDER BY test_id, "order"
        """;

    private static final String SELECT_BY_TEST = """
        SELECT id, test_id, text_of_question, "order"
        FROM question_d
        WHERE test_id = ?
        ORDER BY "order"
        """;

    private static final String DELETE_BY_ID =
        "DELETE FROM question_d WHERE id = ?";

    private static final String DELETE_BY_TEST =
        "DELETE FROM question_d WHERE test_id = ?";

    private static final String EXISTS_BY_ID =
        "SELECT 1 FROM question_d WHERE id = ?";

    private static final String COUNT_BY_TEST = """
        SELECT COUNT(*)
        FROM question_d
        WHERE test_id = ?
        """;

    private static final String SEARCH_BY_TEXT = """
        SELECT id, test_id, text_of_question, "order"
        FROM question_d
        WHERE LOWER(text_of_question) LIKE LOWER(?)
        ORDER BY test_id, "order"
        """;

    /**
     * COALESCE(MAX("order"), -1) → если вопросов нет ИЛИ все order = NULL,
     * вернётся -1, и следующий порядок будет 0.
     */
    private static final String MAX_ORDER_BY_TEST = """
        SELECT COALESCE(MAX("order"), -1)
        FROM question_d
        WHERE test_id = ?
        """;

    private static final String SHIFT_ORDER_SQL = """
        UPDATE question_d
        SET "order" = "order" + ?
        WHERE test_id = ? AND "order" >= ?
        """;

    public QuestionRepository(DataSource dataSource) {
        this.dataSource = dataSource;
    }

    // ---------------------- CRUD ----------------------

    @Override
    public QuestionModel save(QuestionModel question) {
        logger.info("Сохранение нового вопроса для теста: {}", question.getTestId());

        question.validate();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {

            stmt.setObject(1, question.getTestId());
            stmt.setString(2, question.getTextOfQuestion());

            // order может быть null
            if (question.getOrder() != null) {
                stmt.setInt(3, question.getOrder());
            } else {
                stmt.setNull(3, Types.INTEGER);
            }

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                UUID generatedId = rs.getObject("id", UUID.class);
                question.setId(generatedId);
                logger.info("Вопрос сохранён с ID: {}, порядок: {}", generatedId, question.getOrder());
                return question;
            }

            throw new RuntimeException("Не удалось сохранить вопрос");

        } catch (SQLException e) {
            logger.error("Ошибка при сохранении вопроса", e);
            throw new RuntimeException("Ошибка базы данных при сохранении вопроса", e);
        }
    }

    @Override
    public QuestionModel update(QuestionModel question) {
        logger.info("Обновление вопроса с ID: {}", question.getId());

        if (question.getId() == null) {
            throw new IllegalArgumentException("Вопрос не имеет ID");
        }

        question.validate();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {

            stmt.setObject(1, question.getTestId());
            stmt.setString(2, question.getTextOfQuestion());

            if (question.getOrder() != null) {
                stmt.setInt(3, question.getOrder());
            } else {
                stmt.setNull(3, Types.INTEGER);
            }

            stmt.setObject(4, question.getId());

            int updatedRows = stmt.executeUpdate();
            if (updatedRows == 0) {
                throw new RuntimeException("Вопрос с ID " + question.getId() + " не найден");
            }

            logger.info("Вопрос обновлён: {}", question.getId());
            return question;

        } catch (SQLException e) {
            logger.error("Ошибка при обновлении вопроса", e);
            throw new RuntimeException("Ошибка базы данных при обновлении вопроса", e);
        }
    }

    @Override
    public Optional<QuestionModel> findById(UUID id) {
        logger.debug("Поиск вопроса по ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                QuestionModel question = mapRowToQuestion(rs);
                return Optional.of(question);
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка при поиске вопроса по ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при поиске вопроса", e);
        }
    }

    @Override
    public List<QuestionModel> findAll() {
        logger.debug("Получение всех вопросов");
        List<QuestionModel> questions = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_ALL);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                questions.add(mapRowToQuestion(rs));
            }

            logger.debug("Найдено {} вопросов", questions.size());
            return questions;

        } catch (SQLException e) {
            logger.error("Ошибка при получении всех вопросов", e);
            throw new RuntimeException("Ошибка базы данных при получении вопросов", e);
        }
    }

    @Override
    public List<QuestionModel> findByTestId(UUID testId) {
        logger.debug("Поиск вопросов по тесту: {}", testId);
        List<QuestionModel> questions = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                questions.add(mapRowToQuestion(rs));
            }

            logger.debug("Найдено {} вопросов для теста {}", questions.size(), testId);
            return questions;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске вопросов по тесту: {}", testId, e);
            throw new RuntimeException("Ошибка базы данных при поиске вопросов по тесту", e);
        }
    }

    @Override
    public boolean deleteById(UUID id) {
        logger.info("Удаление вопроса с ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(DELETE_BY_ID)) {

            stmt.setObject(1, id);
            int deletedRows = stmt.executeUpdate();

            boolean deleted = deletedRows > 0;
            logger.info("Вопрос с ID {} удалён: {}", id, deleted);
            return deleted;

        } catch (SQLException e) {
            logger.error("Ошибка при удалении вопроса с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при удалении вопроса", e);
        }
    }

    @Override
    public int deleteByTestId(UUID testId) {
        logger.info("Удаление всех вопросов теста: {}", testId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(DELETE_BY_TEST)) {

            stmt.setObject(1, testId);
            int deletedRows = stmt.executeUpdate();

            logger.info("Удалено {} вопросов теста {}", deletedRows, testId);
            return deletedRows;

        } catch (SQLException e) {
            logger.error("Ошибка при удалении вопросов теста: {}", testId, e);
            throw new RuntimeException("Ошибка базы данных при удалении вопросов теста", e);
        }
    }

    @Override
    public boolean existsById(UUID id) {
        logger.debug("Проверка существования вопроса с ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            boolean exists = rs.next();
            logger.debug("Вопрос с ID {} существует: {}", id, exists);
            return exists;

        } catch (SQLException e) {
            logger.error("Ошибка при проверке существования вопроса с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при проверке вопроса", e);
        }
    }

    @Override
    public int countByTestId(UUID testId) {
        logger.debug("Подсчёт вопросов теста: {}", testId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return rs.getInt(1);
            }

            return 0;

        } catch (SQLException e) {
            logger.error("Ошибка при подсчёте вопросов теста: {}", testId, e);
            throw new RuntimeException("Ошибка базы данных при подсчёте вопросов", e);
        }
    }

    @Override
    public List<QuestionModel> findByTextContaining(String text) {
        logger.debug("Поиск вопросов по тексту: {}", text);
        List<QuestionModel> questions = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SEARCH_BY_TEXT)) {

            stmt.setString(1, "%" + text + "%");
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                questions.add(mapRowToQuestion(rs));
            }

            logger.debug("Найдено {} вопросов с текстом '{}'", questions.size(), text);
            return questions;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске вопросов по тексту: {}", text, e);
            throw new RuntimeException("Ошибка базы данных при поиске вопросов", e);
        }
    }

    @Override
    public int getNextOrderForTest(UUID testId) {
        logger.debug("Получение следующего порядкового номера для теста: {}", testId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(MAX_ORDER_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                int maxOrder = rs.getInt(1); // тут null не будет из-за COALESCE
                int nextOrder = maxOrder + 1;
                logger.debug("Следующий порядковый номер для теста {}: {}", testId, nextOrder);
                return nextOrder;
            }

            // теоретически сюда не попадём, но на всякий:
            return 0;

        } catch (SQLException e) {
            logger.error("Ошибка при получении следующего порядкового номера для теста: {}", testId, e);
            throw new RuntimeException("Ошибка базы данных при получении порядкового номера", e);
        }
    }

    @Override
    public int shiftQuestionsOrder(UUID testId, int fromOrder, int shiftBy) {
        logger.info("Сдвиг порядка вопросов теста {} с позиции {} на {}", testId, fromOrder, shiftBy);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SHIFT_ORDER_SQL)) {

            stmt.setInt(1, shiftBy);
            stmt.setObject(2, testId);
            stmt.setInt(3, fromOrder);

            int updatedRows = stmt.executeUpdate();
            logger.info("Обновлено {} вопросов при сдвиге порядка", updatedRows);
            return updatedRows;

        } catch (SQLException e) {
            logger.error("Ошибка при сдвиге порядка вопросов", e);
            throw new RuntimeException("Ошибка базы данных при сдвиге порядка вопросов", e);
        }
    }

    // ---------------------- МАППИНГ ----------------------

    /**
     * Преобразование строки ResultSet в объект QuestionModel.
     */
    private QuestionModel mapRowToQuestion(ResultSet rs) throws SQLException {
        return new QuestionModel(
            rs.getObject("id", UUID.class),
            rs.getObject("test_id", UUID.class),
            rs.getString("text_of_question"),
            // "order" может быть NULL → берём как Integer
            rs.getObject("order", Integer.class)
        );
    }
}