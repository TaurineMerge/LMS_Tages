package com.example.lms.test_attempt.infrastructure.repositories;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Connection;
import java.sql.Date;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Types;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с попытками тестов в базе данных.
 * Реализует {@link TestAttemptRepositoryInterface}.
 *
 * Таблица: testing.test_attempt_b
 * Колонки:
 * - id UUID (PK)
 * - student_id UUID
 * - test_id UUID
 * - date_of_attempt DATE
 * - point INTEGER
 * - certificate_id UUID
 * - attempt_version JSON
 * - attempt_snapshot VARCHAR(256)
 * - completed BOOLEAN
 */
public class TestAttemptRepository implements TestAttemptRepositoryInterface {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptRepository.class);

    private final DatabaseConfig dbConfig;

    // ---------------------------------------------------------------------
    // SQL
    // ---------------------------------------------------------------------

    private static final String INSERT_SQL = """
            INSERT INTO testing.test_attempt_b
                (student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed)
            VALUES
                (?, ?, ?, ?, ?, CAST(? AS json), ?, ?)
            RETURNING id
            """;

    private static final String UPDATE_SQL = """
            UPDATE testing.test_attempt_b
            SET student_id = ?,
                test_id = ?,
                date_of_attempt = ?,
                point = ?,
                certificate_id = ?,
                attempt_version = CAST(? AS json),
                attempt_snapshot = ?,
                completed = ?
            WHERE id = ?
            """;

    private static final String SELECT_BY_ID = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            WHERE id = ?
            """;

    private static final String SELECT_ALL = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            ORDER BY date_of_attempt DESC NULLS LAST
            """;

    private static final String SELECT_BY_STUDENT = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            WHERE student_id = ?
            ORDER BY date_of_attempt DESC NULLS LAST
            """;

    private static final String SELECT_BY_TEST = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            WHERE test_id = ?
            ORDER BY date_of_attempt DESC NULLS LAST
            """;

    private static final String SELECT_BY_STUDENT_AND_TEST = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            WHERE student_id = ? AND test_id = ?
            ORDER BY date_of_attempt DESC NULLS LAST
            """;

    private static final String DELETE_BY_ID = "DELETE FROM testing.test_attempt_b WHERE id = ?";

    private static final String EXISTS_BY_ID = "SELECT 1 FROM testing.test_attempt_b WHERE id = ?";

    private static final String COUNT_BY_STUDENT = "SELECT COUNT(*) FROM testing.test_attempt_b WHERE student_id = ?";

    private static final String COUNT_BY_TEST = "SELECT COUNT(*) FROM testing.test_attempt_b WHERE test_id = ?";

    private static final String SELECT_COMPLETED = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            WHERE COALESCE(completed, FALSE) = TRUE
               OR point IS NOT NULL
            ORDER BY date_of_attempt DESC NULLS LAST
            """;

    private static final String SELECT_INCOMPLETE = """
            SELECT id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed
            FROM testing.test_attempt_b
            WHERE COALESCE(completed, FALSE) = FALSE
              AND point IS NULL
            ORDER BY date_of_attempt DESC NULLS LAST
            """;

    // --- UI: attempt_version ---
    private static final String SELECT_ATTEMPT_VERSION = """
            SELECT attempt_version
            FROM testing.test_attempt_b
            WHERE student_id = ? AND test_id = ? AND date_of_attempt = ?
            """;

    private static final String UPSERT_ATTEMPT_VERSION = """
            INSERT INTO testing.test_attempt_b
                (id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed)
            VALUES
                (?, ?, ?, ?, NULL, NULL, CAST(? AS json), NULL, FALSE)
            ON CONFLICT (student_id, test_id, date_of_attempt)
            DO UPDATE SET attempt_version = EXCLUDED.attempt_version
            """;

    public TestAttemptRepository(DatabaseConfig dbConfig) {
        this.dbConfig = dbConfig;
    }

    private Connection getConnection() throws SQLException {
        return DriverManager.getConnection(dbConfig.getUrl(), dbConfig.getUser(), dbConfig.getPassword());
    }

    // ---------------------------------------------------------------------
    // CRUD
    // ---------------------------------------------------------------------

    /**
     * Сохраняет новую попытку теста в базе данных.
     *
     * @param testAttempt объект TestAttemptModel для сохранения
     * @return сохраненный объект TestAttemptModel с присвоенным идентификатором
     * @throws RuntimeException если не удалось сохранить попытку или произошла
     *                          SQL-ошибка
     */

    @Override
    public TestAttemptModel save(TestAttemptModel testAttempt) {
        testAttempt.validate();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {

            stmt.setObject(1, testAttempt.getStudentId());
            stmt.setObject(2, testAttempt.getTestId());

            if (testAttempt.getDateOfAttempt() != null) {
                stmt.setDate(3, Date.valueOf(testAttempt.getDateOfAttempt()));
            } else {
                stmt.setNull(3, Types.DATE);
            }

            if (testAttempt.getPoint() != null) {
                stmt.setInt(4, testAttempt.getPoint());
            } else {
                stmt.setNull(4, Types.INTEGER);
            }

            if (testAttempt.getCertificateId() != null) {
                stmt.setObject(5, testAttempt.getCertificateId());
            } else {
                stmt.setNull(5, Types.OTHER);
            }

            if (testAttempt.getAttemptVersion() != null && !testAttempt.getAttemptVersion().isBlank()) {
                stmt.setObject(6, testAttempt.getAttemptVersion(), Types.OTHER);
            } else {
                stmt.setNull(6, Types.OTHER);
            }

            if (testAttempt.getAttemptSnapshot() != null) {
                stmt.setString(7, testAttempt.getAttemptSnapshot());
            } else {
                stmt.setNull(7, Types.VARCHAR);
            }

            if (testAttempt.getCompleted() != null) {
                stmt.setBoolean(8, testAttempt.getCompleted());
            } else {
                stmt.setNull(8, Types.BOOLEAN);
            }

            try (ResultSet rs = stmt.executeQuery()) {
                if (rs.next()) {
                    UUID id = rs.getObject("id", UUID.class);
                    testAttempt.setId(id);
                    return testAttempt;
                }
            }

            throw new RuntimeException("Не удалось сохранить попытку теста");

        } catch (SQLException e) {
            logger.error("Ошибка БД при сохранении попытки", e);
            throw new RuntimeException("Ошибка базы данных при сохранении попытки", e);
        }
    }

    /**
     * Обновляет существующую попытку теста в базе данных.
     *
     * @param testAttempt объект TestAttemptModel с обновленными данными
     * @return обновленный объект TestAttemptModel
     * @throws IllegalArgumentException если попытка не имеет идентификатора
     * @throws RuntimeException         если попытка не найдена или произошла
     *                                  SQL-ошибка
     */
    @Override
    public TestAttemptModel update(TestAttemptModel testAttempt) {
        if (testAttempt.getId() == null) {
            throw new IllegalArgumentException("TestAttempt ID cannot be null for update");
        }

        testAttempt.validate();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {

            stmt.setObject(1, testAttempt.getStudentId());
            stmt.setObject(2, testAttempt.getTestId());

            if (testAttempt.getDateOfAttempt() != null) {
                stmt.setDate(3, Date.valueOf(testAttempt.getDateOfAttempt()));
            } else {
                stmt.setNull(3, Types.DATE);
            }

            if (testAttempt.getPoint() != null) {
                stmt.setInt(4, testAttempt.getPoint());
            } else {
                stmt.setNull(4, Types.INTEGER);
            }

            if (testAttempt.getCertificateId() != null) {
                stmt.setObject(5, testAttempt.getCertificateId());
            } else {
                stmt.setNull(5, Types.OTHER);
            }

            if (testAttempt.getAttemptVersion() != null && !testAttempt.getAttemptVersion().isBlank()) {
                stmt.setObject(6, testAttempt.getAttemptVersion(), Types.OTHER);
            } else {
                stmt.setNull(6, Types.OTHER);
            }

            if (testAttempt.getAttemptSnapshot() != null) {
                stmt.setString(7, testAttempt.getAttemptSnapshot());
            } else {
                stmt.setNull(7, Types.VARCHAR);
            }

            if (testAttempt.getCompleted() != null) {
                stmt.setBoolean(8, testAttempt.getCompleted());
            } else {
                stmt.setNull(8, Types.BOOLEAN);
            }

            stmt.setObject(5, testAttempt.getCertificateId());
            stmt.setObject(6, testAttempt.getAttemptVersion());
            stmt.setString(7, testAttempt.getAttemptSnapshot());
            stmt.setBoolean(8, testAttempt.getCompleted());

            stmt.setObject(9, testAttempt.getId());

            int updated = stmt.executeUpdate();
            if (updated == 0) {
                throw new RuntimeException("Попытка с ID " + testAttempt.getId() + " не найдена");
            }

            return testAttempt;

        } catch (SQLException e) {
            logger.error("Ошибка БД при обновлении попытки", e);
            throw new RuntimeException("Ошибка базы данных при обновлении попытки", e);
        }
    }

    /**
     * Находит попытку теста по её идентификатору.
     *
     * @param id уникальный идентификатор попытки
     * @return Optional с найденной попыткой или пустой Optional, если попытка не
     *         найдена
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public Optional<TestAttemptModel> findById(UUID id) {
        String sql = """
                SELECT id, student_id, test_id, date_of_attempt, point,
                       certificate_id, attempt_version, attempt_snapshot, completed
                FROM testing.test_attempt_b
                WHERE id = ?
                """;

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {

            stmt.setObject(1, id);

            try (ResultSet rs = stmt.executeQuery()) {
                if (rs.next()) {
                    return Optional.of(mapRowToTestAttempt(rs));
                }
            }
            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка БД при поиске попытки по id={}", id, e);
            throw new RuntimeException("Ошибка базы данных при поиске попытки", e);
        }
    }

    @Override
    public List<TestAttemptModel> findAll() {
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_ALL);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка БД при получении всех попыток", e);
            throw new RuntimeException("Ошибка базы данных при получении попыток", e);
        }
    }

    @Override
    public boolean deleteById(UUID id) {
        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(DELETE_BY_ID)) {

            stmt.setObject(1, id);
            return stmt.executeUpdate() > 0;

        } catch (SQLException e) {
            logger.error("Ошибка БД при удалении попытки id={}", id, e);
            throw new RuntimeException("Ошибка базы данных при удалении попытки", e);
        }
    }

    @Override
    public boolean existsById(UUID id) {
        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(EXISTS_BY_ID)) {

            stmt.setObject(1, id);
            try (ResultSet rs = stmt.executeQuery()) {
                return rs.next();
            }

        } catch (SQLException e) {
            logger.error("Ошибка БД при проверке существования попытки id={}", id, e);
            throw new RuntimeException("Ошибка базы данных при проверке попытки", e);
        }
    }

    // ---------------------------------------------------------------------
    // Queries
    // ---------------------------------------------------------------------

    @Override
    public List<TestAttemptModel> findByStudentId(UUID studentId) {
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_STUDENT)) {

            stmt.setObject(1, studentId);

            try (ResultSet rs = stmt.executeQuery()) {
                while (rs.next()) {
                    attempts.add(mapRowToTestAttempt(rs));
                }
            }
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка БД при поиске попыток по student_id={}", studentId, e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток по студенту", e);
        }
    }

    @Override
    public List<TestAttemptModel> findByTestId(UUID testId) {
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_TEST)) {

            stmt.setObject(1, testId);

            try (ResultSet rs = stmt.executeQuery()) {
                while (rs.next()) {
                    attempts.add(mapRowToTestAttempt(rs));
                }
            }
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка БД при поиске попыток по test_id={}", testId, e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток по тесту", e);
        }
    }

    /**
     * Находит все попытки для указанного студента и теста.
     *
     * @param studentId идентификатор студента
     * @param testId    идентификатор теста
     * @return список попыток, отсортированных по дате (новые сначала)
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<TestAttemptModel> findByStudentIdAndTestId(UUID studentId, UUID testId) {
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_BY_STUDENT_AND_TEST)) {

            stmt.setObject(1, studentId);
            stmt.setObject(2, testId);

            try (ResultSet rs = stmt.executeQuery()) {
                while (rs.next()) {
                    attempts.add(mapRowToTestAttempt(rs));
                }
            }
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка БД при поиске попыток по student_id={} и test_id={}", studentId, testId, e);
            throw new RuntimeException("Ошибка базы данных при поиске попыток", e);
        }
    }

    @Override
    public int countByStudentId(UUID studentId) {
        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(COUNT_BY_STUDENT)) {

            stmt.setObject(1, studentId);

            try (ResultSet rs = stmt.executeQuery()) {
                return rs.next() ? rs.getInt(1) : 0;
            }

        } catch (SQLException e) {
            logger.error("Ошибка БД при подсчёте попыток по student_id={}", studentId, e);
            throw new RuntimeException("Ошибка базы данных при подсчёте попыток", e);
        }
    }

    @Override
    public int countByTestId(UUID testId) {
        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(COUNT_BY_TEST)) {

            stmt.setObject(1, testId);

            try (ResultSet rs = stmt.executeQuery()) {
                return rs.next() ? rs.getInt(1) : 0;
            }

        } catch (SQLException e) {
            logger.error("Ошибка БД при подсчёте попыток по test_id={}", testId, e);
            throw new RuntimeException("Ошибка базы данных при подсчёте попыток", e);
        }
    }

    @Override
    public List<TestAttemptModel> findCompletedAttempts() {
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_COMPLETED);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка БД при поиске завершённых попыток", e);
            throw new RuntimeException("Ошибка базы данных при поиске завершённых попыток", e);
        }
    }

    @Override
    public List<TestAttemptModel> findIncompleteAttempts() {
        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_INCOMPLETE);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка БД при поиске незавершённых попыток", e);
            throw new RuntimeException("Ошибка базы данных при поиске незавершённых попыток", e);
        }
    }

    // ---------------------------------------------------------------------
    // UI: attempt_version
    // ---------------------------------------------------------------------

    @Override
    public Optional<String> findAttemptVersion(UUID studentId, UUID testId, String date) {
        LocalDate parsed = parseDateOrNull(date);
        if (parsed == null) {
            return Optional.empty();
        }

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(SELECT_ATTEMPT_VERSION)) {

            stmt.setObject(1, studentId);
            stmt.setObject(2, testId);
            stmt.setDate(3, Date.valueOf(parsed));

            try (ResultSet rs = stmt.executeQuery()) {
                if (rs.next()) {
                    return Optional.ofNullable(rs.getString("attempt_version"));
                }
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка БД при чтении attempt_version (student={}, test={}, date={})",
                    studentId, testId, date, e);
            throw new RuntimeException("Ошибка базы данных при поиске attempt_version", e);
        }
    }

    @Override
    public void upsertAttemptVersion(UUID studentId, UUID testId, String date, String attemptVersionJson) {
        LocalDate parsed = parseDateOrNull(date);
        if (parsed == null) {
            throw new IllegalArgumentException("date_of_attempt is invalid: " + date);
        }

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(UPSERT_ATTEMPT_VERSION)) {

            stmt.setObject(1, UUID.randomUUID());
            stmt.setObject(2, studentId);
            stmt.setObject(3, testId);
            stmt.setDate(4, Date.valueOf(parsed));

            if (attemptVersionJson != null && !attemptVersionJson.isBlank()) {
                stmt.setObject(5, attemptVersionJson, Types.OTHER);
            } else {
                stmt.setNull(5, Types.OTHER);
            }

            stmt.executeUpdate();

        } catch (SQLException e) {
            logger.error("Ошибка БД при upsert attempt_version (student={}, test={}, date={})",
                    studentId, testId, date, e);
            throw new RuntimeException("Ошибка базы данных при сохранении attempt_version", e);
        }
    }

    private LocalDate parseDateOrNull(String date) {
        if (date == null || date.isBlank()) {
            return null;
        }
        try {
            return LocalDate.parse(date);
        } catch (Exception e) {
            return null;
        }
    }

    private TestAttemptModel mapRowToTestAttempt(ResultSet rs) throws SQLException {
        Date date = rs.getDate("date_of_attempt");

        LocalDate localDate = (date != null) ? date.toLocalDate() : null;

        Boolean completed = null;
        Object completedRaw = rs.getObject("completed");
        if (completedRaw != null) {
            completed = rs.getBoolean("completed");
        }

        return new TestAttemptModel(
                rs.getObject("id", UUID.class),
                rs.getObject("student_id", UUID.class),
                rs.getObject("test_id", UUID.class),
                localDate,
                rs.getObject("point", Integer.class),
                rs.getObject("certificate_id", UUID.class),
                rs.getString("attempt_version"),
                rs.getString("attempt_snapshot"),
                completed);
    }
}