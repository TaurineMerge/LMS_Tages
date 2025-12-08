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
 * Реализация репозитория для работы с таблицей answer_d
 */
public class AnswerRepository implements AnswerRepositoryInterface {
    private static final Logger logger = LoggerFactory.getLogger(AnswerRepository.class);
    private final DataSource dataSource;
    
    // SQL запросы для таблицы answer_d
    private static final String INSERT_SQL = """
        INSERT INTO answer_d (question_id, answer_text, is_correct, display_order, explanation)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id
        """;
    
    private static final String UPDATE_SQL = """
        UPDATE answer_d 
        SET question_id = ?, answer_text = ?, is_correct = ?, display_order = ?, explanation = ?
        WHERE id = ?
        """;
    
    private static final String SELECT_BY_ID = """
        SELECT id, question_id, answer_text, is_correct, display_order, explanation
        FROM answer_d 
        WHERE id = ?
        """;
    
    private static final String SELECT_ALL = """
        SELECT id, question_id, answer_text, is_correct, display_order, explanation
        FROM answer_d 
        ORDER BY display_order NULLS FIRST, answer_text
        """;
    
    private static final String SELECT_BY_QUESTION = """
        SELECT id, question_id, answer_text, is_correct, display_order, explanation
        FROM answer_d 
        WHERE question_id = ?
        ORDER BY display_order NULLS FIRST, answer_text
        """;
    
    private static final String SELECT_CORRECT_BY_QUESTION = """
        SELECT id, question_id, answer_text, is_correct, display_order, explanation
        FROM answer_d 
        WHERE question_id = ? AND is_correct = true
        ORDER BY display_order NULLS FIRST
        """;
    
    private static final String DELETE_BY_ID = "DELETE FROM answer_d WHERE id = ?";
    
    private static final String DELETE_BY_QUESTION = "DELETE FROM answer_d WHERE question_id = ?";
    
    private static final String EXISTS_BY_ID = "SELECT 1 FROM answer_d WHERE id = ?";
    
    private static final String EXISTS_BY_QUESTION_AND_TEXT = """
        SELECT 1 FROM answer_d 
        WHERE question_id = ? AND LOWER(answer_text) = LOWER(?)
        """;
    
    private static final String COUNT_BY_QUESTION = """
        SELECT COUNT(*) 
        FROM answer_d 
        WHERE question_id = ?
        """;
    
    private static final String COUNT_CORRECT_BY_QUESTION = """
        SELECT COUNT(*) 
        FROM answer_d 
        WHERE question_id = ? AND is_correct = true
        """;
    
    public AnswerRepository(DataSource dataSource) {
        this.dataSource = dataSource;
    }
    
    @Override
    public AnswerModel save(AnswerModel answer) {
        logger.info("Сохранение нового ответа для вопроса: {}", answer.getQuestionId());
        
        // Валидация перед сохранением
        answer.validate();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {
            
            // Устанавливаем параметры
            stmt.setObject(1, answer.getQuestionId());
            stmt.setString(2, answer.getAnswerText());
            stmt.setBoolean(3, answer.getIsCorrect());
            
            if (answer.getDisplayOrder() != null) {
                stmt.setInt(4, answer.getDisplayOrder());
            } else {
                stmt.setNull(4, Types.INTEGER);
            }
            
            if (answer.getExplanation() != null && !answer.getExplanation().isEmpty()) {
                stmt.setString(5, answer.getExplanation());
            } else {
                stmt.setNull(5, Types.VARCHAR);
            }
            
            // Выполняем запрос и получаем сгенерированный ID
            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                UUID generatedId = rs.getObject("id", UUID.class);
                answer.setId(generatedId);
                logger.info("Ответ сохранен с ID: {}", generatedId);
                return answer;
            }
            
            throw new RuntimeException("Не удалось сохранить ответ");
            
        } catch (SQLException e) {
            logger.error("Ошибка при сохранении ответа", e);
            throw new RuntimeException("Ошибка базы данных при сохранении ответа", e);
        }
    }
    
    @Override
    public AnswerModel update(AnswerModel answer) {
        logger.info("Обновление ответа с ID: {}", answer.getId());
        
        if (answer.getId() == null) {
            throw new IllegalArgumentException("Ответ не имеет ID");
        }
        
        // Валидация перед обновлением
        answer.validate();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {
            
            stmt.setObject(1, answer.getQuestionId());
            stmt.setString(2, answer.getAnswerText());
            stmt.setBoolean(3, answer.getIsCorrect());
            
            if (answer.getDisplayOrder() != null) {
                stmt.setInt(4, answer.getDisplayOrder());
            } else {
                stmt.setNull(4, Types.INTEGER);
            }
            
            if (answer.getExplanation() != null && !answer.getExplanation().isEmpty()) {
                stmt.setString(5, answer.getExplanation());
            } else {
                stmt.setNull(5, Types.VARCHAR);
            }
            
            stmt.setObject(6, answer.getId());
            
            int updatedRows = stmt.executeUpdate();
            if (updatedRows == 0) {
                throw new RuntimeException("Ответ с ID " + answer.getId() + " не найден");
            }
            
            logger.info("Ответ обновлен: {}", answer.getId());
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
                AnswerModel answer = mapRowToAnswer(rs);
                return Optional.of(answer);
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
            logger.info("Ответ с ID {} удален: {}", id, deleted);
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
    public boolean existsByQuestionIdAndText(UUID questionId, String answerText) {
        logger.debug("Проверка существования ответа '{}' для вопроса {}", answerText, questionId);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_QUESTION_AND_TEXT)) {
            
            stmt.setObject(1, questionId);
            stmt.setString(2, answerText);
            ResultSet rs = stmt.executeQuery();
            
            boolean exists = rs.next();
            logger.debug("Ответ '{}' для вопроса {} существует: {}", answerText, questionId, exists);
            return exists;
            
        } catch (SQLException e) {
            logger.error("Ошибка при проверке существования ответа для вопроса", e);
            throw new RuntimeException("Ошибка базы данных при проверке ответа", e);
        }
    }
    
    @Override
    public int countByQuestionId(UUID questionId) {
        logger.debug("Подсчет ответов для вопроса: {}", questionId);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_BY_QUESTION)) {
            
            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();
            
            if (rs.next()) {
                return rs.getInt(1);
            }
            
            return 0;
            
        } catch (SQLException e) {
            logger.error("Ошибка при подсчете ответов для вопроса: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при подсчете ответов", e);
        }
    }
    
    @Override
    public int countCorrectAnswersByQuestionId(UUID questionId) {
        logger.debug("Подсчет правильных ответов для вопроса: {}", questionId);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_CORRECT_BY_QUESTION)) {
            
            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();
            
            if (rs.next()) {
                return rs.getInt(1);
            }
            
            return 0;
            
        } catch (SQLException e) {
            logger.error("Ошибка при подсчете правильных ответов для вопроса: {}", questionId, e);
            throw new RuntimeException("Ошибка базы данных при подсчете правильных ответов", e);
        }
    }
    
    /**
     * Преобразование строки ResultSet в объект AnswerModel
     */
    private AnswerModel mapRowToAnswer(ResultSet rs) throws SQLException {
        return new AnswerModel(
            rs.getObject("id", UUID.class),
            rs.getObject("question_id", UUID.class),
            rs.getString("answer_text"),
            rs.getBoolean("is_correct"),
            rs.getObject("display_order", Integer.class),
            rs.getString("explanation")
        );
    }
}
