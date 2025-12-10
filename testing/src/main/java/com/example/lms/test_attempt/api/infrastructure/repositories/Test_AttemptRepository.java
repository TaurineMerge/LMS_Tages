package com.example.lms.test_attempt.api.infrastructure.repositories;

import com.example.lms.test_attempt.api.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.api.domain.repository.Test_AttemptRepositoryInterface;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.sql.DataSource;
import java.sql.*;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Реализация репозитория попыток тестов с использованием JDBC.
 * Работает с таблицей test_attempt_b в PostgreSQL.
 */
public class Test_AttemptRepository implements Test_AttemptRepositoryInterface {

    private static final Logger logger = LoggerFactory.getLogger(Test_AttemptRepository.class);
    private final DataSource dataSource;

    // SQL-запросы для таблицы test_attempt_b
    private static final String INSERT_SQL = """
        INSERT INTO test_attempt_b (student_id, test_id, date_of_attempt,
                                   point, certificate_id, attempt_version)
        VALUES (?, ?, ?, ?, ?, ?)
        RETURNING id
        """;

    private static final String UPDATE_SQL = """
        UPDATE test_attempt_b
        SET student_id = ?, test_id = ?, date_of_attempt = ?,
            point = ?, certificate_id = ?, attempt_version = ?
        WHERE id = ?
        """;

    private static final String SELECT_BY_ID = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE id = ?
        """;

    private static final String SELECT_ALL = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        ORDER BY date_of_attempt DESC
        """;

    private static final String SELECT_BY_STUDENT = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE student_id = ?
        ORDER BY date_of_attempt DESC
        """;

    private static final String SELECT_BY_TEST = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE test_id = ?
        ORDER BY date_of_attempt DESC
        """;

    private static final String SELECT_BY_STUDENT_AND_TEST = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE student_id = ? AND test_id = ?
        ORDER BY date_of_attempt DESC
        """;

    private static final String DELETE_BY_ID =
        "DELETE FROM test_attempt_b WHERE id = ?";

    private static final String EXISTS_BY_ID =
        "SELECT 1 FROM test_attempt_b WHERE id = ?";

    private static final String COUNT_BY_STUDENT_AND_TEST = """
        SELECT COUNT(*)
        FROM test_attempt_b
        WHERE student_id = ? AND test_id = ?
        """;

    private static final String SELECT_BY_DATE = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE date_of_attempt = ?
        ORDER BY date_of_attempt DESC
        """;

    private static final String SELECT_COMPLETED = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE point IS NOT NULL
        ORDER BY date_of_attempt DESC
        """;

    private static final String SELECT_INCOMPLETE = """
        SELECT id, student_id, test_id, date_of_attempt,
               point, certificate_id, attempt_version
        FROM test_attempt_b
        WHERE point IS NULL
        ORDER BY date_of_attempt DESC
        """;

    public Test_AttemptRepository(DataSource dataSource) {
        this.dataSource = dataSource;
    }

    @Override
    public TestAttemptModel save(TestAttemptModel attempt) {
        logger.info("Сохранение новой попытки теста для студента: {}", attempt.getStudentId());

        attempt.validate();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {

            stmt.setObject(1, attempt.getStudentId());
            stmt.setObject(2, attempt.getTestId());
            stmt.setDate(3, Date.valueOf(attempt.getDateOfAttempt()));

            if (attempt.getPoint() != null) {
                stmt.setInt(4, attempt.getPoint());
            } else {
                stmt.setNull(4, Types.INTEGER);
            }

            if (attempt.getCertificateId() != null) {
                stmt.setObject(5, attempt.getCertificateId());
            } else {
                // для UUID в Postgres обычно Types.OTHER
                stmt.setNull(5, Types.OTHER);
            }

            stmt.setString(6, attempt.getAttemptVersion());

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                UUID generatedId = rs.getObject("id", UUID.class);
                attempt.setId(generatedId);
                logger.info("Попытка сохранена с ID: {}", generatedId);
                return attempt;
            }

            throw new RuntimeException("Не удалось сохранить попытку");

        } catch (SQLException e) {
            logger.error("Ошибка при сохранении попытки", e);
            throw new RuntimeException("Ошибка базы данных при сохранении попытки", e);
        }
    }

    @Override
    public TestAttemptModel update(TestAttemptModel attempt) {
        logger.info("Обновление попытки с ID: {}", attempt.getId());

        if (attempt.getId() == null) {
            throw new IllegalArgumentException("Попытка не имеет ID");
        }

        attempt.validate();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {

            stmt.setObject(1, attempt.getStudentId());
            stmt.setObject(2, attempt.getTestId());
            stmt.setDate(3, Date.valueOf(attempt.getDateOfAttempt()));

            if (attempt.getPoint() != null) {
                stmt.setInt(4, attempt.getPoint());
            } else {
                stmt.setNull(4, Types.INTEGER);
            }

            if (attempt.getCertificateId() != null) {
                stmt.setObject(5, attempt.getCertificateId());
            } else {
                stmt.setNull(5, Types.OTHER);
            }

            stmt.setString(6, attempt.getAttemptVersion());
            stmt.setObject(7, attempt.getId());

            int updatedRows = stmt.executeUpdate();
            if (updatedRows == 0) {
                throw new RuntimeException("Попытка с ID " + attempt.getId() + " не найдена");
            }

            logger.info("Попытка обновлена: {}", attempt.getId());
            return attempt;

        } catch (SQLException e) {
            logger.error("Ошибка при обновлении попытки", e);
            throw new RuntimeException("Ошибка базы данных при обновлении попытки", e);
        }
    }

    @Override
    public Optional<TestAttemptModel> findById(UUID id) {
        logger.debug("Поиск попытки по ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                TestAttemptModel attempt = mapRowToTestAttempt(rs);
                return Optional.of(attempt);
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка при поиске попытки по ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при поиске попытки", e);
        }
    }

    @Override
    public List<TestAttemptModel> findAll() {
        logger.debug("Получение всех попыток");
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_ALL);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} попыток", attempts.size());
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при получении всех попыток", e);
            throw new RuntimeException("Ошибка базы данных при получении попыток", e);
        }
    }

    @Override
    public List<TestAttemptModel> findByStudentId(UUID studentId) {
        logger.debug("Поиск попыток по студенту: {}", studentId);
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_STUDENT)) {

            stmt.setObject(1, studentId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} попыток для студента {}", attempts.size(), studentId);
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске попыток по студенту: {}", studentId, e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток по студенту", e);
        }
    }

    @Override
    public List<TestAttemptModel> findByTestId(UUID testId) {
        logger.debug("Поиск попыток по тесту: {}", testId);
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_TEST)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} попыток для теста {}", attempts.size(), testId);
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске попыток по тесту: {}", testId, e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток по тесту", e);
        }
    }

    @Override
    public List<TestAttemptModel> findByStudentAndTest(UUID studentId, UUID testId) {
        logger.debug("Поиск попыток студента {} по тесту {}", studentId, testId);
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_STUDENT_AND_TEST)) {

            stmt.setObject(1, studentId);
            stmt.setObject(2, testId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} попыток", attempts.size());
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске попыток студента по тесту", e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток", e);
        }
    }

    @Override
    public boolean deleteById(UUID id) {
        logger.info("Удаление попытки с ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(DELETE_BY_ID)) {

            stmt.setObject(1, id);
            int deletedRows = stmt.executeUpdate();

            boolean deleted = deletedRows > 0;
            logger.info("Попытка с ID {} удалена: {}", id, deleted);
            return deleted;

        } catch (SQLException e) {
            logger.error("Ошибка при удалении попытки с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при удалении попытки", e);
        }
    }

    @Override
    public boolean existsById(UUID id) {
        logger.debug("Проверка существования попытки с ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            boolean exists = rs.next();
            logger.debug("Попытка с ID {} существует: {}", id, exists);
            return exists;

        } catch (SQLException e) {
            logger.error("Ошибка при проверке существования попытки с ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при проверке попытки", e);
        }
    }

    @Override
    public int countAttemptsByStudentAndTest(UUID studentId, UUID testId) {
        logger.debug("Подсчёт попыток студента {} по тесту {}", studentId, testId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(COUNT_BY_STUDENT_AND_TEST)) {

            stmt.setObject(1, studentId);
            stmt.setObject(2, testId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return rs.getInt(1);
            }

            return 0;

        } catch (SQLException e) {
            logger.error("Ошибка при подсчёте попыток", e);
            throw new RuntimeException("Ошибка базы данных при подсчёте попыток", e);
        }
    }

    @Override
    public List<TestAttemptModel> findByDate(LocalDate date) {
        logger.debug("Поиск попыток за дату: {}", date);
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_DATE)) {

            stmt.setDate(1, Date.valueOf(date));
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} попыток за {}", attempts.size(), date);
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске попыток по дате: {}", date, e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток по дате", e);
        }
    }

    @Override
    public List<TestAttemptModel> findCompletedAttempts() {
        logger.debug("Поиск завершённых попыток");
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_COMPLETED);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} завершённых попыток", attempts.size());
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске завершённых попыток", e);
            throw new RuntimeException("Ошибка базы данных при поиске завершённых попыток", e);
        }
    }

    @Override
    public List<TestAttemptModel> findIncompleteAttempts() {
        logger.debug("Поиск незавершённых попыток");
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_INCOMPLETE);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} незавершённых попыток", attempts.size());
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске незавершённых попыток", e);
            throw new RuntimeException("Ошибка базы данных при поиске незавершённых попыток", e);
        }
    }

    /**
     * Преобразование строки ResultSet в объект TestAttemptModel.
     */
    private TestAttemptModel mapRowToTestAttempt(ResultSet rs) throws SQLException {
        return new TestAttemptModel(
                rs.getObject("id", UUID.class),
                rs.getObject("student_id", UUID.class),
                rs.getObject("test_id", UUID.class),
                rs.getDate("date_of_attempt").toLocalDate(),
                rs.getObject("point", Integer.class),
                rs.getObject("certificate_id", UUID.class),
                rs.getString("attempt_version")
        );
    }
}