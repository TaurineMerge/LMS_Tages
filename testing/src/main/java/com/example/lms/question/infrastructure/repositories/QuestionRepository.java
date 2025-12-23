package com.example.lms.question.infrastructure.repositories;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Types;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.question.domain.model.QuestionModel;
import com.example.lms.question.domain.repository.QuestionRepositoryInterface;

public class QuestionRepository implements QuestionRepositoryInterface {
    private static final Logger logger = LoggerFactory.getLogger(QuestionRepository.class);
    private final DatabaseConfig dbConfig;

    // Обновленные SQL запросы
    private static final String INSERT_SQL = """
        INSERT INTO testing.question_d (test_id, draft_id, text_of_question, "order")
        VALUES (?, ?, ?, ?)
        RETURNING id
        """.trim();

    private static final String UPDATE_SQL = """
        UPDATE testing.question_d
        SET test_id = ?, draft_id = ?, text_of_question = ?, "order" = ?
        WHERE id = ?
        """.trim();

    private static final String SELECT_BY_ID = """
        SELECT id, test_id, draft_id, text_of_question, "order"
        FROM testing.question_d
        WHERE id = ?
        """.trim();

    private static final String SELECT_ALL = """
        SELECT id, test_id, draft_id, text_of_question, "order"
        FROM testing.question_d
        ORDER BY COALESCE(test_id, draft_id), "order"
        """.trim();

    private static final String SELECT_BY_TEST = """
        SELECT id, test_id, draft_id, text_of_question, "order"
        FROM testing.question_d
        WHERE test_id = ?
        ORDER BY "order"
        """.trim();

    private static final String SELECT_BY_DRAFT = """
        SELECT id, test_id, draft_id, text_of_question, "order"
        FROM testing.question_d
        WHERE draft_id = ?
        ORDER BY "order"
        """.trim();

    private static final String DELETE_BY_ID = "DELETE FROM testing.question_d WHERE id = ?";
    
    private static final String DELETE_BY_DRAFT = "DELETE FROM testing.question_d WHERE draft_id = ?";

    private static final String COUNT_BY_TEST = """
        SELECT COUNT(*)
        FROM testing.question_d
        WHERE test_id = ?
        """.trim();

    private static final String COUNT_BY_DRAFT = """
        SELECT COUNT(*)
        FROM testing.question_d
        WHERE draft_id = ?
        """.trim();

    private static final String SEARCH_BY_TEXT = """
        SELECT id, test_id, draft_id, text_of_question, "order"
        FROM testing.question_d
        WHERE LOWER(text_of_question) LIKE LOWER(?)
        ORDER BY COALESCE(test_id, draft_id), "order"
        """.trim();

    private static final String MAX_ORDER_BY_TEST = """
        SELECT COALESCE(MAX("order"), -1)
        FROM testing.question_d
        WHERE test_id = ?
        """.trim();

    private static final String MAX_ORDER_BY_DRAFT = """
        SELECT COALESCE(MAX("order"), -1)
        FROM testing.question_d
        WHERE draft_id = ?
        """.trim();

    private static final String SHIFT_ORDER_SQL = """
        UPDATE testing.question_d
        SET "order" = "order" + ?
        WHERE test_id = ? AND "order" >= ?
        """.trim();

    private static final String SHIFT_DRAFT_ORDER_SQL = """
        UPDATE testing.question_d
        SET "order" = "order" + ?
        WHERE draft_id = ? AND "order" >= ?
        """.trim();

    public QuestionRepository(DatabaseConfig dbConfig) {
        this.dbConfig = dbConfig;
    }

    private Connection getConnection() throws SQLException {
        return DriverManager.getConnection(
                dbConfig.getUrl(),
                dbConfig.getUser(),
                dbConfig.getPassword());
    }

    static {
        try {
            Class.forName("org.postgresql.Driver");
            logger.info("PostgreSQL драйвер зарегистрирован");
        } catch (ClassNotFoundException e) {
            logger.error("Не удалось зарегистрировать драйвер PostgreSQL", e);
            throw new RuntimeException("Драйвер БД не найден", e);
        }
    }

    @Override
    public QuestionModel save(QuestionModel question) {
        logger.info("Сохранение нового вопроса");
        question.validate();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {

            // test_id может быть null
            if (question.getTestId() != null) {
                stmt.setObject(1, question.getTestId());
            } else {
                stmt.setNull(1, Types.OTHER);
            }
            
            // draft_id может быть null
            if (question.getDraftId() != null) {
                stmt.setObject(2, question.getDraftId());
            } else {
                stmt.setNull(2, Types.OTHER);
            }
            
            stmt.setString(3, question.getTextOfQuestion());
            stmt.setInt(4, question.getOrder());

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                UUID id = rs.getObject("id", UUID.class);
                question.setId(id);
                return question;
            }

            throw new RuntimeException("Не удалось сохранить вопрос");

        } catch (SQLException e) {
            logger.error("Ошибка при сохранении вопроса", e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public QuestionModel update(QuestionModel question) {
        logger.info("Обновление вопроса id={}", question.getId());

        if (question.getId() == null) {
            throw new IllegalArgumentException("Вопрос должен иметь ID для обновления");
        }

        question.validate();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {

            // test_id может быть null
            if (question.getTestId() != null) {
                stmt.setObject(1, question.getTestId());
            } else {
                stmt.setNull(1, Types.OTHER);
            }
            
            // draft_id может быть null
            if (question.getDraftId() != null) {
                stmt.setObject(2, question.getDraftId());
            } else {
                stmt.setNull(2, Types.OTHER);
            }
            
            stmt.setString(3, question.getTextOfQuestion());
            stmt.setInt(4, question.getOrder());
            stmt.setObject(5, question.getId());

            int updated = stmt.executeUpdate();
            if (updated == 0) {
                throw new RuntimeException("Вопрос с ID " + question.getId() + " не найден");
            }

            return question;

        } catch (SQLException e) {
            logger.error("Ошибка при обновлении вопроса", e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public Optional<QuestionModel> findById(UUID id) {
        logger.debug("Поиск вопроса id={}", id);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return Optional.of(mapRowToQuestion(rs));
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка поиска вопроса id={}", id, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public List<QuestionModel> findAll() {
        logger.debug("Получение всех вопросов");

        List<QuestionModel> result = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_ALL);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                result.add(mapRowToQuestion(rs));
            }

            return result;

        } catch (SQLException e) {
            logger.error("Ошибка при получении всех вопросов", e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public List<QuestionModel> findByTestId(UUID testId) {
        logger.debug("Поиск вопросов теста id={}", testId);

        List<QuestionModel> result = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                result.add(mapRowToQuestion(rs));
            }

            return result;

        } catch (SQLException e) {
            logger.error("Ошибка поиска вопросов теста id={}", testId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public List<QuestionModel> findByDraftId(UUID draftId) {
        logger.debug("Поиск вопросов черновика id={}", draftId);

        List<QuestionModel> result = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_DRAFT)) {

            stmt.setObject(1, draftId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                result.add(mapRowToQuestion(rs));
            }

            return result;

        } catch (SQLException e) {
            logger.error("Ошибка поиска вопросов черновика id={}", draftId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public boolean deleteById(UUID id) {
        logger.info("Удаление вопроса id={}", id);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(DELETE_BY_ID)) {

            stmt.setObject(1, id);
            int rows = stmt.executeUpdate();
            return rows > 0;

        } catch (SQLException e) {
            logger.error("Ошибка удаления вопроса id={}", id, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public boolean deleteByDraftId(UUID draftId) {
        logger.info("Удаление вопросов черновика id={}", draftId);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(DELETE_BY_DRAFT)) {

            stmt.setObject(1, draftId);
            int rows = stmt.executeUpdate();
            return rows > 0;

        } catch (SQLException e) {
            logger.error("Ошибка удаления вопросов черновика id={}", draftId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public int countByTestId(UUID testId) {
        logger.debug("Подсчёт вопросов теста id={}", testId);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(COUNT_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            return rs.next() ? rs.getInt(1) : 0;

        } catch (SQLException e) {
            logger.error("Ошибка countByTestId id={}", testId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public int countByDraftId(UUID draftId) {
        logger.debug("Подсчёт вопросов черновика id={}", draftId);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(COUNT_BY_DRAFT)) {

            stmt.setObject(1, draftId);
            ResultSet rs = stmt.executeQuery();

            return rs.next() ? rs.getInt(1) : 0;

        } catch (SQLException e) {
            logger.error("Ошибка countByDraftId id={}", draftId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public List<QuestionModel> findByTextContaining(String text) {
        logger.debug("Поиск вопросов содержащих '{}'", text);

        List<QuestionModel> result = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SEARCH_BY_TEXT)) {

            stmt.setString(1, "%" + text + "%");
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                result.add(mapRowToQuestion(rs));
            }

            return result;

        } catch (SQLException e) {
            logger.error("Ошибка поиска по тексту '{}'", text, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public int getNextOrderForTest(UUID testId) {
        logger.debug("Получение следующего порядка для теста id={}", testId);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(MAX_ORDER_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            return rs.next() ? rs.getInt(1) + 1 : 0;

        } catch (SQLException e) {
            logger.error("Ошибка getNextOrderForTest id={}", testId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public int getNextOrderForDraft(UUID draftId) {
        logger.debug("Получение следующего порядка для черновика id={}", draftId);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(MAX_ORDER_BY_DRAFT)) {

            stmt.setObject(1, draftId);
            ResultSet rs = stmt.executeQuery();

            return rs.next() ? rs.getInt(1) + 1 : 0;

        } catch (SQLException e) {
            logger.error("Ошибка getNextOrderForDraft id={}", draftId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public int shiftQuestionsOrder(UUID testId, int fromOrder, int shiftBy) {
        logger.info("Сдвиг порядка вопросов теста id={} с {} на {}", testId, fromOrder, shiftBy);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SHIFT_ORDER_SQL)) {

            stmt.setInt(1, shiftBy);
            stmt.setObject(2, testId);
            stmt.setInt(3, fromOrder);

            return stmt.executeUpdate();

        } catch (SQLException e) {
            logger.error("Ошибка shiftQuestionsOrder testId={}", testId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    @Override
    public int shiftDraftQuestionsOrder(UUID draftId, int fromOrder, int shiftBy) {
        logger.info("Сдвиг порядка вопросов черновика id={} с {} на {}", draftId, fromOrder, shiftBy);

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SHIFT_DRAFT_ORDER_SQL)) {

            stmt.setInt(1, shiftBy);
            stmt.setObject(2, draftId);
            stmt.setInt(3, fromOrder);

            return stmt.executeUpdate();

        } catch (SQLException e) {
            logger.error("Ошибка shiftDraftQuestionsOrder draftId={}", draftId, e);
            throw new RuntimeException("Ошибка базы данных", e);
        }
    }

    private QuestionModel mapRowToQuestion(ResultSet rs) throws SQLException {
        return new QuestionModel(
                rs.getObject("id", UUID.class),
                rs.getObject("test_id", UUID.class),
                rs.getObject("draft_id", UUID.class),
                rs.getString("text_of_question"),
                rs.getInt("order"));
    }
}