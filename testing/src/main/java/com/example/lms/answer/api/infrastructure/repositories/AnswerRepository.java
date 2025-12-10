package com.example.lms.answer.api.infrastructure.repositories;

import com.example.lms.answer.api.domain.model.AnswerModel;
import com.example.lms.answer.api.domain.repository.AnswerRepositoryInterface;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.sql.DataSource;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Реализация репозитория для работы с таблицей answer_d.
 *
 * Структура таблицы:
 *  - id          UUID (PK)
 *  - question_id UUID (FK → question_d.id)
 *  - text        TEXT  (текст ответа)
 *  - score       INT   (кол-во баллов за ответ)
 *
 * Договорённость:
 *  - "правильный" ответ — это ответ, у которого score > 0.
 */
public class AnswerRepository implements AnswerRepositoryInterface {

    private static final Logger logger = LoggerFactory.getLogger(AnswerRepository.class);

    private final DataSource dataSource;

    // ---------------------- SQL ЗАПРОСЫ ----------------------

    /** Вставка нового ответа и возврат сгенерированного ID. */
    private static final String INSERT_SQL = """
        INSERT INTO answer_d (question_id, text, score)
        VALUES (?, ?, ?)
        RETURNING id
        """;

    /** Обновление существующего ответа по его ID. */
    private static final String UPDATE_SQL = """
        UPDATE answer_d
        SET question_id = ?, text = ?, score = ?
        WHERE id = ?
        """;

    /** Поиск ответа по ID. */
    private static final String SELECT_BY_ID = """
        SELECT id, question_id, text, score
        FROM answer_d
        WHERE id = ?
        """;

    /** Получение всех ответов. */
    private static final String SELECT_ALL = """
        SELECT id, question_id, text, score
        FROM answer_d
        ORDER BY question_id, text
        """;

    /** Получение всех ответов для конкретного вопроса. */
    private static final String SELECT_BY_QUESTION = """
        SELECT id, question_id, text, score
        FROM answer_d
        WHERE question_id = ?
        ORDER BY text
        """;

    /** Получение правильных ответов для конкретного вопроса (score > 0). */
    private static final String SELECT_CORRECT_BY_QUESTION = """
        SELECT id, question_id, text, score
        FROM answer_d
        WHERE question_id = ? AND score > 0
        ORDER BY text
        """;

    /** Удаление ответа по ID. */
    private static final String DELETE_BY_ID =
        "DELETE FROM answer_d WHERE id = ?";

    /** Удаление всех ответов по ID вопроса. */
    private static final String DELETE_BY_QUESTION =
        "DELETE FROM answer_d WHERE question_id = ?";

    /** Проверка существования ответа по ID. */
    private static final String EXISTS_BY_ID =
        "SELECT 1 FROM answer_d WHERE id = ?";

    /** Проверка наличия ответа с таким же текстом для вопроса. */
    private static final String EXISTS_BY_QUESTION_AND_TEXT = """
        SELECT 1
        FROM answer_d
        WHERE question_id = ? AND LOWER(text) = LOWER(?)
        """;

    /** Подсчёт ответов для вопроса. */
    private static final String COUNT_BY_QUESTION = """
        SELECT COUNT(*)
        FROM answer_d
        WHERE question_id = ?
        """;

    /** Подсчёт правильных ответов (score > 0) для вопроса. */
    private static final String COUNT_CORRECT_BY_QUESTION = """
        SELECT COUNT(*)
        FROM answer_d
        WHERE question_id = ? AND score > 0
        """;


    public AnswerRepository(DataSource dataSource) {
        this.dataSource = dataSource;
    }

    // ---------------------- CRUD ----------------------

    @Override
    public AnswerModel save(AnswerModel answer) {
        logger.info("Сохранение нового ответа для вопроса: {}", answer.getQuestionId());

        // Валидация доменной модели перед сохранением
        answer.validate();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {

            stmt.setObject(1, answer.getQuestionId());
            stmt.setString(2, answer.getText());
            stmt.setInt(3, answer.getScore());

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                UUID generatedId = rs.getObject("id", UUID.class);
                answer.setId(generatedId);
                logger.info("Ответ сохранён с ID: {}", generatedId);
                return answer;
            }

            throw new RuntimeException("Не удалось сохранить ответ: не получен ID");

        } catch (SQLException e) {
            logger.error("Ошибка при сохранении ответа", e);
            throw new RuntimeException("Ошибка базы данных при сохранении ответа", e);
        }
    }

    @Override
    public AnswerModel update(AnswerModel answer) {
        logger.info("Обновление ответа с ID: {}", answer.getId());

        if (answer.getId() == null) {
            throw new IllegalArgumentException("Ответ не имеет ID (null)");
        }

        // Валидация перед обновлением
        answer.validate();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {

            stmt.setObject(1, answer.getQuestionId());
            stmt.setString(2, answer.getText());
            stmt.setInt(3, answer.getScore());
            stmt.setObject(4, answer.getId());

            int updatedRows = stmt.executeUpdate();
            if (updatedRows == 0) {
                throw new RuntimeException("Ответ с ID " + answer.getId() + " не найден для обновления");
            }

            logger.info("Ответ обновлён: {}", answer.getId());
            return answer;

        } catch (SQLException e) {
            logger.error("Ошибка при обновлении ответа", e);
            throw new RuntimeException("Ошибка базы данных при обновлении ответа", e);
        }
    }

    @Override
    public Optional<AnswerModel> findById(UUID id) {
        logger.debug("Поиск ответа по ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return Optional.of(mapRowToAnswer(rs));
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка при поиске ответа по ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при поиске ответа", e);
        }
    }

    @Override
    public List<AnswerModel> findAll() {
        logger.debug("Получение всех ответов");
        List<AnswerModel> answers = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_ALL);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                answers.add(mapRowToAnswer(rs));
            }

            logger.debug("Найдено {} ответов", answers.size());
            return answers;

        } catch (SQLException e) {
            logger.error("Ошибка при получении всех ответов", e);
            throw new RuntimeException("Ошибка базы данных при получении ответов", e);
        }
    }

    @Override
    public List<AnswerModel> findByQuestionId(UUID questionId) {
        logger.debug("Поиск ответов по вопросу: {}", questionId);
        List<AnswerModel> answers = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_QUESTION)) {

            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                answers.add(mapRowToAnswer(rs));
            }

            logger.debug("Найдено {} ответов для вопроса {}", answers.size(), questionId);
            return answers;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске ответов по вопросу: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при поиске ответов по вопросу", e);
        }
    }

    @Override
    public List<AnswerModel> findCorrectAnswersByQuestionId(UUID questionId) {
        logger.debug("Поиск правильных ответов по вопросу: {}", questionId);
        List<AnswerModel> answers = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_CORRECT_BY_QUESTION)) {

            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                answers.add(mapRowToAnswer(rs));
            }

            logger.debug("Найдено {} правильных ответов для вопроса {}", answers.size(), questionId);
            return answers;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске правильных ответов по вопросу: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при поиске правильных ответов", e);
        }
    }

    @Override
    public boolean deleteById(UUID id) {
        logger.info("Удаление ответа с ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(DELETE_BY_ID)) {

            stmt.setObject(1, id);
            int deletedRows = stmt.executeUpdate();

            boolean deleted = deletedRows > 0;
            logger.info("Ответ с ID {} удалён: {}", id, deleted);
            return deleted;

        } catch (SQLException e) {
            logger.error("Ошибка при удалении ответа с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при удалении ответа", e);
        }
    }

    @Override
    public int deleteByQuestionId(UUID questionId) {
        logger.info("Удаление всех ответов для вопроса: {}", questionId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(DELETE_BY_QUESTION)) {

            stmt.setObject(1, questionId);
            int deletedRows = stmt.executeUpdate();

            logger.info("Удалено {} ответов для вопроса {}", deletedRows, questionId);
            return deletedRows;

        } catch (SQLException e) {
            logger.error("Ошибка при удалении ответов для вопроса: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при удалении ответов по вопросу", e);
        }
    }

    @Override
    public boolean existsById(UUID id) {
        logger.debug("Проверка существования ответа с ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            boolean exists = rs.next();
            logger.debug("Ответ с ID {} существует: {}", id, exists);
            return exists;

        } catch (SQLException e) {
            logger.error("Ошибка при проверке существования ответа с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при проверке ответа", e);
        }
    }

    @Override
    public boolean existsByQuestionIdAndText(UUID questionId, String text) {
        logger.debug("Проверка существования ответа '{}' для вопроса {}", text, questionId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_QUESTION_AND_TEXT)) {

            stmt.setObject(1, questionId);
            stmt.setString(2, text);
            ResultSet rs = stmt.executeQuery();

            boolean exists = rs.next();
            logger.debug("Ответ '{}' для вопроса {} существует: {}", text, questionId, exists);
            return exists;

        } catch (SQLException e) {
            logger.error("Ошибка при проверке существования ответа для вопроса", e);
            throw new RuntimeException("Ошибка базы данных при проверке ответа", e);
        }
    }

    @Override
    public int countByQuestionId(UUID questionId) {
        logger.debug("Подсчёт ответов для вопроса: {}", questionId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_BY_QUESTION)) {

            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return rs.getInt(1);
            }

            return 0;

        } catch (SQLException e) {
            logger.error("Ошибка при подсчёте ответов для вопроса: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при подсчёте ответов", e);
        }
    }

    @Override
    public int countCorrectAnswersByQuestionId(UUID questionId) {
        logger.debug("Подсчёт правильных ответов для вопроса: {}", questionId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_CORRECT_BY_QUESTION)) {

            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return rs.getInt(1);
            }

            return 0;

        } catch (SQLException e) {
            logger.error("Ошибка при подсчёте правильных ответов для вопроса: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при подсчёте правильных ответов", e);
        }
    }

    // ---------------------- МАППИНГ ----------------------

    /**
     * Преобразование строки ResultSet в доменную модель AnswerModel.
     */
    private AnswerModel mapRowToAnswer(ResultSet rs) throws SQLException {
        return new AnswerModel(
            rs.getObject("id", UUID.class),
            rs.getObject("question_id", UUID.class),
            rs.getString("text"),
            rs.getInt("score")
        );
    }
}
