package com.example.lms.test.infrastructure.repositories;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.test.domain.model.TestModel;
import com.example.lms.test.domain.repository.TestRepositoryInterface;

import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с таблицей test_d.
 */
public class TestRepository implements TestRepositoryInterface {

    private final DatabaseConfig dbConfig;

    public TestRepository(DatabaseConfig dbConfig) {
        this.dbConfig = dbConfig;
    }

    private Connection getConnection() throws SQLException {
        return DriverManager.getConnection(
                dbConfig.getUrl(),
                dbConfig.getUser(),
                dbConfig.getPassword()
        );
    }

    @Override
    public TestModel save(TestModel test) {
        test.validate();

        String sql = """
            INSERT INTO test_d (course_id, title, min_point, description)
            VALUES (?, ?, ?, ?)
            RETURNING id
            """;

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, test.getCourseId());
            stmt.setString(2, test.getTitle());

            if (test.getMinPoint() != null)
                stmt.setInt(3, test.getMinPoint());
            else
                stmt.setNull(3, Types.INTEGER);

            stmt.setString(4, test.getDescription());

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                test.setId(rs.getObject("id", UUID.class));
                return test;
            }

            throw new RuntimeException("Не удалось сохранить тест");

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при сохранении теста", e);
        }
    }

    @Override
    public TestModel update(TestModel test) {
        if (test.getId() == null) {
            throw new IllegalArgumentException("Тест должен иметь ID");
        }

        test.validate();

        String sql = """
            UPDATE test_d
            SET course_id = ?, title = ?, min_point = ?, description = ?
            WHERE id = ?
            """;

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, test.getCourseId());
            stmt.setString(2, test.getTitle());

            if (test.getMinPoint() != null)
                stmt.setInt(3, test.getMinPoint());
            else
                stmt.setNull(3, Types.INTEGER);

            stmt.setString(4, test.getDescription());
            stmt.setObject(5, test.getId());

            int updated = stmt.executeUpdate();
            if (updated == 0)
                throw new RuntimeException("Тест с ID " + test.getId() + " не найден");

            return test;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при обновлении теста", e);
        }
    }

    @Override
    public Optional<TestModel> findById(UUID id) {
        String sql = """
            SELECT id, course_id, title, min_point, description
            FROM test_d
            WHERE id = ?
            """;

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next())
                return Optional.of(mapRowToTest(rs));

            return Optional.empty();

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске теста по ID", e);
        }
    }

    @Override
    public List<TestModel> findAll() {
        String sql = """
            SELECT id, course_id, title, min_point, description
            FROM test_d
            ORDER BY title
            """;

        List<TestModel> tests = new ArrayList<>();

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next())
                tests.add(mapRowToTest(rs));

            return tests;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при получении всех тестов", e);
        }
    }

    @Override
    public boolean deleteById(UUID id) {
        String sql = "DELETE FROM test_d WHERE id = ?";

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            return stmt.executeUpdate() > 0;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при удалении теста", e);
        }
    }

    @Override
    public boolean existsById(UUID id) {
        String sql = "SELECT 1 FROM test_d WHERE id = ?";

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            return stmt.executeQuery().next();

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при проверке существования теста", e);
        }
    }

    @Override
    public List<TestModel> findByTitleContaining(String title) {
        String sql = """
            SELECT id, course_id, title, min_point, description
            FROM test_d
            WHERE LOWER(title) LIKE LOWER(?)
            ORDER BY title
            """;

        List<TestModel> tests = new ArrayList<>();

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setString(1, "%" + title + "%");
            ResultSet rs = stmt.executeQuery();

            while (rs.next())
                tests.add(mapRowToTest(rs));

            return tests;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске тестов по названию", e);
        }
    }

    @Override
    public int countByCourseId(UUID courseId) {
        String sql = "SELECT COUNT(*) FROM test_d WHERE course_id = ?";

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, courseId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next())
                return rs.getInt(1);

            return 0;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при подсчёте тестов для курса", e);
        }
    }

    @Override
    public List<TestModel> findByCourseId(UUID courseId) {
        String sql = """
            SELECT id, course_id, title, min_point, description
            FROM test_d
            WHERE course_id = ?
            ORDER BY title
            """;

        List<TestModel> tests = new ArrayList<>();

        try (Connection conn = getConnection();
             PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, courseId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next())
                tests.add(mapRowToTest(rs));

            return tests;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске тестов по course_id", e);
        }
    }

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