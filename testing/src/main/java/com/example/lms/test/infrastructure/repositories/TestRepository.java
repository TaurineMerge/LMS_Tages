package com.example.lms.test.infrastructure.repositories;

import com.example.lms.test.domain.model.TestModel;
import com.example.lms.test.domain.repository.TestRepositoryInterface;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.sql.DataSource;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Реализация репозитория тестов с использованием JDBC
 * Работает с таблицей TEST_D в PostgreSQL
 */
public class TestRepository implements TestRepositoryInterface {
    private static final Logger logger = LoggerFactory.getLogger(TestRepository.class);
    private final DataSource dataSource;
    
    // SQL запросы для таблицы TEST_D
    private static final String INSERT_SQL = """
        INSERT INTO test_d (course_id, title, min_point, description)
        VALUES (?, ?, ?, ?)
        RETURNING id
        """;
    
    private static final String UPDATE_SQL = """
        UPDATE test_d 
        SET course_id = ?, title = ?, min_point = ?, description = ?
        WHERE id = ?
        """;
    
    private static final String SELECT_BY_ID = """
        SELECT id, course_id, title, min_point, description 
        FROM test_d 
        WHERE id = ?
        """;
    
    private static final String SELECT_ALL = """
        SELECT id, course_id, title, min_point, description 
        FROM test_d 
        ORDER BY title
        """;
    
    private static final String SELECT_BY_COURSE = """
        SELECT id, course_id, title, min_point, description 
        FROM test_d 
        WHERE course_id = ?
        ORDER BY title
        """;
    
    private static final String DELETE_BY_ID = "DELETE FROM test_d WHERE id = ?";
    
    private static final String EXISTS_BY_ID = "SELECT 1 FROM test_d WHERE id = ?";
    
    private static final String SEARCH_BY_TITLE = """
        SELECT id, course_id, title, min_point, description 
        FROM test_d 
        WHERE LOWER(title) LIKE LOWER(?)
        ORDER BY title
        """;
    
    private static final String COUNT_BY_COURSE = """
        SELECT COUNT(*) 
        FROM test_d 
        WHERE course_id = ?
        """;
    
    public TestRepository(DataSource dataSource) {
        this.dataSource = dataSource;
    }
    
    @Override
    public TestModel save(TestModel test) {
        logger.info("Сохранение нового теста: {}", test.getTitle());
        
        // Валидация перед сохранением
        test.validate();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {
            
            // Устанавливаем параметры
            stmt.setObject(1, test.getCourseId());
            stmt.setString(2, test.getTitle());
            
            if (test.getMinPoint() != null) {
                stmt.setInt(3, test.getMinPoint());
            } else {
                stmt.setNull(3, Types.INTEGER);
            }
            
            stmt.setString(4, test.getDescription());
            
            // Выполняем запрос и получаем сгенерированный ID
            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                UUID generatedId = rs.getObject("id", UUID.class);
                test.setId(generatedId);
                logger.info("Тест сохранен с ID: {}", generatedId);
                return test;
            }
            
            throw new RuntimeException("Не удалось сохранить тест");
            
        } catch (SQLException e) {
            logger.error("Ошибка при сохранении теста", e);
            throw new RuntimeException("Ошибка базы данных при сохранении теста", e);
        }
    }
    
    @Override
    public TestModel update(TestModel test) {
        logger.info("Обновление теста с ID: {}", test.getId());
        
        if (test.getId() == null) {
            throw new IllegalArgumentException("Тест не имеет ID");
        }
        
        // Валидация перед обновлением
        test.validate();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {
            
            stmt.setObject(1, test.getCourseId());
            stmt.setString(2, test.getTitle());
            
            if (test.getMinPoint() != null) {
                stmt.setInt(3, test.getMinPoint());
            } else {
                stmt.setNull(3, Types.INTEGER);
            }
            
            stmt.setString(4, test.getDescription());
            stmt.setObject(5, test.getId());
            
            int updatedRows = stmt.executeUpdate();
            if (updatedRows == 0) {
                throw new RuntimeException("Тест с ID " + test.getId() + " не найден");
            }
            
            logger.info("Тест обновлен: {}", test.getId());
            return test;
            
        } catch (SQLException e) {
            logger.error("Ошибка при обновлении теста", e);
            throw new RuntimeException("Ошибка базы данных при обновлении теста", e);
        }
    }
    
    @Override
    public Optional<TestModel> findById(UUID id) {
        logger.debug("Поиск теста по ID: {}", id);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {
            
            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();
            
            if (rs.next()) {
                TestModel test = mapRowToTest(rs);
                return Optional.of(test);
            }
            
            return Optional.empty();
            
        } catch (SQLException e) {
            logger.error("Ошибка при поиске теста по ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при поиске теста", e);
        }
    }
    
    @Override
    public List<TestModel> findAll() {
        logger.debug("Получение всех тестов");
        List<TestModel> tests = new ArrayList<>();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_ALL);
             ResultSet rs = stmt.executeQuery()) {
            
            while (rs.next()) {
                tests.add(mapRowToTest(rs));
            }
            
            logger.debug("Найдено {} тестов", tests.size());
            return tests;
            
        } catch (SQLException e) {
            logger.error("Ошибка при получении всех тестов", e);
            throw new RuntimeException("Ошибка базы данных при получении тестов", e);
        }
    }
    
    @Override
    public List<TestModel> findByCourseId(UUID courseId) {
        logger.debug("Поиск тестов по курсу: {}", courseId);
        List<TestModel> tests = new ArrayList<>();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_COURSE)) {
            
            stmt.setObject(1, courseId);
            ResultSet rs = stmt.executeQuery();
            
            while (rs.next()) {
                tests.add(mapRowToTest(rs));
            }
            
            logger.debug("Найдено {} тестов для курса {}", tests.size(), courseId);
            return tests;
            
        } catch (SQLException e) {
            logger.error("Ошибка при поиске тестов по курсу: {}", courseId, e);
            throw new RuntimeException("Ошибка базы данных при поиске тестов по курсу", e);
        }
    }
    
    @Override
    public boolean deleteById(UUID id) {
        logger.info("Удаление теста с ID: {}", id);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(DELETE_BY_ID)) {
            
            stmt.setObject(1, id);
            int deletedRows = stmt.executeUpdate();
            
            boolean deleted = deletedRows > 0;
            logger.info("Тест с ID {} удален: {}", id, deleted);
            return deleted;
            
        } catch (SQLException e) {
            logger.error("Ошибка при удалении теста с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при удалении теста", e);
        }
    }
    
    @Override
    public boolean existsById(UUID id) {
        logger.debug("Проверка существования теста с ID: {}", id);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_ID)) {
            
            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();
            
            boolean exists = rs.next();
            logger.debug("Тест с ID {} существует: {}", id, exists);
            return exists;
            
        } catch (SQLException e) {
            logger.error("Ошибка при проверке существования теста с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при проверке теста", e);
        }
    }
    
    @Override
    public List<TestModel> findByTitleContaining(String title) {
        logger.debug("Поиск тестов по названию: {}", title);
        List<TestModel> tests = new ArrayList<>();
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SEARCH_BY_TITLE)) {
            
            stmt.setString(1, "%" + title + "%");
            ResultSet rs = stmt.executeQuery();
            
            while (rs.next()) {
                tests.add(mapRowToTest(rs));
            }
            
            logger.debug("Найдено {} тестов по названию '{}'", tests.size(), title);
            return tests;
            
        } catch (SQLException e) {
            logger.error("Ошибка при поиске тестов по названию: {}", title, e);
            throw new RuntimeException("Ошибка базы данных при поиске тестов", e);
        }
    }
    
    @Override
    public int countByCourseId(UUID courseId) {
        logger.debug("Подсчет тестов для курса: {}", courseId);
        
        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_BY_COURSE)) {
            
            stmt.setObject(1, courseId);
            ResultSet rs = stmt.executeQuery();
            
            if (rs.next()) {
                return rs.getInt(1);
            }
            
            return 0;
            
        } catch (SQLException e) {
            logger.error("Ошибка при подсчете тестов для курса: {}", courseId, e);
            throw new RuntimeException("Ошибка базы данных при подсчете тестов", e);
        }
    }
    
    /**
     * Преобразование строки ResultSet в объект TestModel
     */
    private TestModel mapRowToTest(ResultSet rs) throws SQLException {
        return new TestModel(
            rs.getObject("id", UUID.class),
            rs.getObject("course_id", UUID.class),
            rs.getString("title"),
            rs.getObject("min_point", Integer.class),
            rs.getString("description")
        );
    }
}