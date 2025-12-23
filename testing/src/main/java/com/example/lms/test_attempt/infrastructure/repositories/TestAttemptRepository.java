package com.example.lms.test_attempt.infrastructure.repositories;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;

import java.sql.*;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с попытками тестов в базе данных.
 * Реализует интерфейс TestAttemptRepositoryInterface для операций CRUD и поиска
 * попыток тестов.
 */
public class TestAttemptRepository implements TestAttemptRepositoryInterface {

    private final DatabaseConfig dbConfig;

    /**
     * Конструктор репозитория.
     *
     * @param dbConfig конфигурация базы данных, содержащая параметры подключения
     */
    public TestAttemptRepository(DatabaseConfig dbConfig) {
        this.dbConfig = dbConfig;
    }

    /**
     * Устанавливает соединение с базой данных.
     *
     * @return активное соединение с базой данных
     * @throws SQLException если произошла ошибка при установке соединения
     */
    private Connection getConnection() throws SQLException {
        return DriverManager.getConnection(
                dbConfig.getUrl(),
                dbConfig.getUser(),
                dbConfig.getPassword());
    }

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

        String sql = """
                INSERT INTO testing.test_attempt_b
                (student_id, test_id, date_of_attempt, point, certificate_id,
                 attempt_version, attempt_snapshot, completed)
                VALUES (?, ?, ?, ?, ?, ?::json, ?, ?)
                RETURNING id
                """;

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, testAttempt.getStudentId());
            stmt.setObject(2, testAttempt.getTestId());

            if (testAttempt.getDateOfAttempt() != null)
                stmt.setDate(3, Date.valueOf(testAttempt.getDateOfAttempt()));
            else
                stmt.setNull(3, Types.DATE);

            if (testAttempt.getPoint() != null)
                stmt.setInt(4, testAttempt.getPoint());
            else
                stmt.setNull(4, Types.INTEGER);

            stmt.setObject(5, testAttempt.getCertificateId());
            stmt.setString(6, testAttempt.getAttemptVersion());
            stmt.setString(7, testAttempt.getAttemptSnapshot());
            stmt.setBoolean(8, testAttempt.getCompleted());

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                testAttempt.setId(rs.getObject("id", UUID.class));
                return testAttempt;
            }

            throw new RuntimeException("Не удалось сохранить попытку теста");

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при сохранении попытки теста", e);
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
            throw new IllegalArgumentException("Попытка теста должна иметь ID");
        }

        testAttempt.validate();

        String sql = """
                UPDATE testing.test_attempt_b
                SET student_id = ?, test_id = ?, date_of_attempt = ?, point = ?,
                    certificate_id = ?, attempt_version = ?::json, attempt_snapshot = ?,
                    completed = ?
                WHERE id = ?
                """;

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, testAttempt.getStudentId());
            stmt.setObject(2, testAttempt.getTestId());

            if (testAttempt.getDateOfAttempt() != null)
                stmt.setDate(3, Date.valueOf(testAttempt.getDateOfAttempt()));
            else
                stmt.setNull(3, Types.DATE);

            if (testAttempt.getPoint() != null)
                stmt.setInt(4, testAttempt.getPoint());
            else
                stmt.setNull(4, Types.INTEGER);

            stmt.setObject(5, testAttempt.getCertificateId());
            stmt.setObject(6, testAttempt.getAttemptVersion());
            stmt.setString(7, testAttempt.getAttemptSnapshot());
            stmt.setBoolean(8, testAttempt.getCompleted());
            stmt.setObject(9, testAttempt.getId());

            int updated = stmt.executeUpdate();
            if (updated == 0)
                throw new RuntimeException("Попытка теста с ID " + testAttempt.getId() + " не найдена");

            return testAttempt;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при обновлении попытки теста", e);
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
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next())
                return Optional.of(mapRowToTestAttempt(rs));

            return Optional.empty();

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске попытки теста по ID", e);
        }
    }

    /**
     * Получает все попытки тестов из базы данных.
     *
     * @return список всех попыток, отсортированных по дате (новые сначала)
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<TestAttemptModel> findAll() {
        String sql = """
                SELECT id, student_id, test_id, date_of_attempt, point,
                       certificate_id, attempt_version, attempt_snapshot, completed
                FROM testing.test_attempt_b
                ORDER BY date_of_attempt DESC NULLS LAST
                """;

        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next())
                attempts.add(mapRowToTestAttempt(rs));

            return attempts;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при получении всех попыток тестов", e);
        }
    }

    /**
     * Удаляет попытку теста по её идентификатору.
     *
     * @param id уникальный идентификатор попытки для удаления
     * @return true если попытка была удалена, false если попытка не найдена
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public boolean deleteById(UUID id) {
        String sql = "DELETE FROM testing.test_attempt_b WHERE id = ?";

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            return stmt.executeUpdate() > 0;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при удалении попытки теста", e);
        }
    }

    /**
     * Проверяет существование попытки теста по идентификатору.
     *
     * @param id уникальный идентификатор попытки
     * @return true если попытка существует, false если нет
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public boolean existsById(UUID id) {
        String sql = "SELECT 1 FROM testing.test_attempt_b WHERE id = ?";

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            return stmt.executeQuery().next();

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при проверке существования попытки теста", e);
        }
    }

    /**
     * Находит все попытки теста для указанного студента.
     *
     * @param studentId идентификатор студента
     * @return список попыток студента, отсортированных по дате (новые сначала)
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<TestAttemptModel> findByStudentId(UUID studentId) {
        String sql = """
                SELECT id, student_id, test_id, date_of_attempt, point,
                       certificate_id, attempt_version, attempt_snapshot, completed
                FROM testing.test_attempt_b
                WHERE student_id = ?
                ORDER BY date_of_attempt DESC NULLS LAST
                """;

        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, studentId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next())
                attempts.add(mapRowToTestAttempt(rs));

            return attempts;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске попыток теста по student_id", e);
        }
    }

    /**
     * Находит все попытки для указанного теста.
     *
     * @param testId идентификатор теста
     * @return список попыток теста, отсортированных по дате (новые сначала)
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<TestAttemptModel> findByTestId(UUID testId) {
        String sql = """
                SELECT id, student_id, test_id, date_of_attempt, point,
                       certificate_id, attempt_version, attempt_snapshot, completed
                FROM testing.test_attempt_b
                WHERE test_id = ?
                ORDER BY date_of_attempt DESC NULLS LAST
                """;

        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next())
                attempts.add(mapRowToTestAttempt(rs));

            return attempts;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске попыток теста по test_id", e);
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
        String sql = """
                SELECT id, student_id, test_id, date_of_attempt, point,
                       certificate_id, attempt_version, attempt_snapshot, completed
                FROM testing.test_attempt_b
                WHERE student_id = ? AND test_id = ?
                ORDER BY date_of_attempt DESC NULLS LAST
                """;

        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, studentId);
            stmt.setObject(2, testId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next())
                attempts.add(mapRowToTestAttempt(rs));

            return attempts;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске попыток теста по student_id и test_id", e);
        }
    }

    /**
     * Подсчитывает количество попыток для указанного студента.
     *
     * @param studentId идентификатор студента
     * @return количество попыток студента
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public int countByStudentId(UUID studentId) {
        String sql = "SELECT COUNT(*) FROM testing.test_attempt_b WHERE student_id = ?";

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, studentId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next())
                return rs.getInt(1);

            return 0;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при подсчёте попыток теста для студента", e);
        }
    }

    /**
     * Подсчитывает количество попыток для указанного теста.
     *
     * @param testId идентификатор теста
     * @return количество попыток теста
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public int countByTestId(UUID testId) {
        String sql = "SELECT COUNT(*) FROM testing.test_attempt_b WHERE test_id = ?";

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, testId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next())
                return rs.getInt(1);

            return 0;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при подсчёте попыток теста", e);
        }
    }

    /**
     * Находит все завершенные попытки.
     *
     * @return список завершенных попыток
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<TestAttemptModel> findCompletedAttempts() {
        String sql = """
                SELECT id, student_id, test_id, date_of_attempt, point,
                       certificate_id, attempt_version, attempt_snapshot, completed
                FROM testing.test_attempt_b
                WHERE completed = true
                ORDER BY date_of_attempt DESC
                """;

        List<TestAttemptModel> attempts = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next())
                attempts.add(mapRowToTestAttempt(rs));

            return attempts;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске завершенных попыток теста", e);
        }
    }

    /**
     * Преобразует строку ResultSet в объект TestAttemptModel.
     *
     * @param rs ResultSet, содержащий данные попытки
     * @return объект TestAttemptModel, созданный из данных ResultSet
     * @throws SQLException если произошла ошибка при чтении данных из ResultSet
     */
    private TestAttemptModel mapRowToTestAttempt(ResultSet rs) throws SQLException {
        Date date = rs.getDate("date_of_attempt");
        LocalDate localDate = date != null ? date.toLocalDate() : null;

        return new TestAttemptModel(
                rs.getObject("id", UUID.class),
                rs.getObject("student_id", UUID.class),
                rs.getObject("test_id", UUID.class),
                localDate,
                rs.getObject("point", Integer.class),
                rs.getObject("certificate_id", UUID.class),
                rs.getString("attempt_version"),
                rs.getString("attempt_snapshot"),
                rs.getBoolean("completed"));
    }
}