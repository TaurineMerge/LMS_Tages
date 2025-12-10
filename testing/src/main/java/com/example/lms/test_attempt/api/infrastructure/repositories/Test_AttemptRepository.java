package com.example.lms.test_attempt.api.infrastructure.repositories;

import com.example.lms.test_attempt.api.domain.model.Test_AttemptModel;
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
 * Репозиторий для работы с попытками прохождения тестов.
 * <p>
 * Реализация {@link Test_AttemptRepositoryInterface} на основе JDBC.
 * Работает с таблицей {@code test_attempt_b} в PostgreSQL и выполняет:
 * <ul>
 *     <li>CRUD-операции над попытками</li>
 *     <li>поиск по студенту, тесту, дате</li>
 *     <li>поиск завершённых/незавершённых попыток</li>
 *     <li>подсчёт количества попыток</li>
 *     <li>поиск наилучшей попытки студента по тесту</li>
 * </ul>
 */
public class Test_AttemptRepository implements Test_AttemptRepositoryInterface {

    private static final Logger logger = LoggerFactory.getLogger(Test_AttemptRepository.class);

    /** Источник соединений с БД. */
    private final DataSource dataSource;

    // SQL запросы для таблицы test_attempt_b
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

    private static final String SELECT_BEST_ATTEMPT = """
        SELECT id, student_id, test_id, date_of_attempt, 
               point, certificate_id, attempt_version
        FROM test_attempt_b 
        WHERE student_id = ? AND test_id = ? AND point IS NOT NULL
        ORDER BY point DESC, date_of_attempt DESC
        LIMIT 1
        """;

    /**
     * Создаёт репозиторий попыток тестов.
     *
     * @param dataSource источник соединений с PostgreSQL
     */
    public Test_AttemptRepository(DataSource dataSource) {
        this.dataSource = dataSource;
    }

    /**
     * Сохраняет новую попытку теста.
     * <p>
     * После успешной вставки в БД устанавливает сгенерированный {@code id}
     * в переданную модель.
     *
     * @param attempt доменная модель попытки
     * @return сохранённая модель с установленным ID
     * @throws RuntimeException при ошибке работы с БД или отсутствии результата
     */
    @Override
    public Test_AttemptModel save(Test_AttemptModel attempt) {
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
                stmt.setNull(5, Types.OTHER); // тип UUID в PostgreSQL
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

    /**
     * Обновляет существующую попытку теста.
     *
     * @param attempt модель с изменёнными данными
     * @return обновлённая модель
     * @throws IllegalArgumentException если у попытки нет ID
     * @throws RuntimeException         если запись не найдена или при ошибках SQL
     */
    @Override
    public Test_AttemptModel update(Test_AttemptModel attempt) {
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

    /**
     * Находит попытку по её идентификатору.
     *
     * @param id ID попытки
     * @return Optional с моделью или пустой Optional, если не найдено
     */
    @Override
    public Optional<Test_AttemptModel> findById(UUID id) {
        logger.debug("Поиск попытки по ID: {}", id);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BY_ID)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                Test_AttemptModel attempt = mapRowToTestAttempt(rs);
                return Optional.of(attempt);
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка при поиске попытки по ID: {}", id, e);
            throw new RuntimeException("Ошибка базы данных при поиске попытки", e);
        }
    }

    /**
     * Возвращает список всех попыток.
     *
     * @return список всех попыток, отсортированных по дате (по убыванию)
     */
    @Override
    public List<Test_AttemptModel> findAll() {
        logger.debug("Получение всех попыток");
        List<Test_AttemptModel> attempts = new ArrayList<>();

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

    /**
     * Возвращает все попытки указанного студента.
     *
     * @param studentId ID студента
     * @return список попыток, отсортированных по дате (по убыванию)
     */
    @Override
    public List<Test_AttemptModel> findByStudentId(UUID studentId) {
        logger.debug("Поиск попыток по студенту: {}", studentId);
        List<Test_AttemptModel> attempts = new ArrayList<>();

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

    /**
     * Возвращает все попытки по конкретному тесту.
     *
     * @param testId ID теста
     * @return список попыток, отсортированных по дате (по убыванию)
     */
    @Override
    public List<Test_AttemptModel> findByTestId(UUID testId) {
        logger.debug("Поиск попыток по тесту: {}", testId);
        List<Test_AttemptModel> attempts = new ArrayList<>();

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

    /**
     * Возвращает попытки конкретного студента по конкретному тесту.
     *
     * @param studentId ID студента
     * @param testId    ID теста
     * @return список попыток (могут быть как завершённые, так и нет)
     */
    @Override
    public List<Test_AttemptModel> findByStudentAndTest(UUID studentId, UUID testId) {
        logger.debug("Поиск попыток студента {} по тесту {}", studentId, testId);
        List<Test_AttemptModel> attempts = new ArrayList<>();

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

    /**
     * Удаляет попытку по ID.
     *
     * @param id ID попытки
     * @return true — если попытка была удалена; false — если не найдена
     */
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

    /**
     * Проверяет существование попытки по ID.
     *
     * @param id ID попытки
     * @return true — если запись существует, иначе false
     */
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

    /**
     * Считает количество попыток конкретного студента по определённому тесту.
     *
     * @param studentId ID студента
     * @param testId    ID теста
     * @return количество попыток
     */
    @Override
    public int countAttemptsByStudentAndTest(UUID studentId, UUID testId) {
        logger.debug("Подсчет попыток студента {} по тесту {}", studentId, testId);

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
            logger.error("Ошибка при подсчете попыток", e);
            throw new RuntimeException("Ошибка базы данных при подсчете попыток", e);
        }
    }

    /**
     * Возвращает список попыток за указанную дату.
     *
     * @param date дата попытки
     * @return список попыток за этот день
     */
    @Override
    public List<Test_AttemptModel> findByDate(LocalDate date) {
        logger.debug("Поиск попыток за дату: {}", date);
        List<Test_AttemptModel> attempts = new ArrayList<>();

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

    /**
     * Возвращает список завершённых попыток (где {@code point IS NOT NULL}).
     *
     * @return список завершённых попыток
     */
    @Override
    public List<Test_AttemptModel> findCompletedAttempts() {
        logger.debug("Поиск завершенных попыток");
        List<Test_AttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_COMPLETED);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} завершенных попыток", attempts.size());
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске завершенных попыток", e);
            throw new RuntimeException("Ошибка базы данных при поиске завершенных попыток", e);
        }
    }

    /**
     * Возвращает список незавершённых попыток (где {@code point IS NULL}).
     *
     * @return список незавершённых попыток
     */
    @Override
    public List<Test_AttemptModel> findIncompleteAttempts() {
        logger.debug("Поиск незавершенных попыток");
        List<Test_AttemptModel> attempts = new ArrayList<>();

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_INCOMPLETE);
             ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                attempts.add(mapRowToTestAttempt(rs));
            }

            logger.debug("Найдено {} незавершенных попыток", attempts.size());
            return attempts;

        } catch (SQLException e) {
            logger.error("Ошибка при поиске незавершенных попыток", e);
            throw new RuntimeException("Ошибка базы данных при поиске незавершенных попыток", e);
        }
    }

    /**
     * Ищет лучшую попытку студента по тесту:
     * <ul>
     *     <li>с максимальным количеством баллов</li>
     *     <li>при равенстве баллов — с более поздней датой</li>
     * </ul>
     *
     * @param studentId ID студента
     * @param testId    ID теста
     * @return Optional с лучшей попыткой или пустой Optional
     */
    @Override
    public Optional<Test_AttemptModel> findBestAttemptByStudentAndTest(UUID studentId, UUID testId) {
        logger.debug("Поиск лучшей попытки студента {} по тесту {}", studentId, testId);

        try (Connection conn = dataSource.getConnection();
             PreparedStatement stmt = conn.prepareStatement(SELECT_BEST_ATTEMPT)) {

            stmt.setObject(1, studentId);
            stmt.setObject(2, testId);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                Test_AttemptModel attempt = mapRowToTestAttempt(rs);
                return Optional.of(attempt);
            }

            return Optional.empty();

        } catch (SQLException e) {
            logger.error("Ошибка при поиске лучшей попытки", e);
            throw new RuntimeException("Ошибка базы данных при поиске лучшей попытки", e);
        }
    }

    /**
     * Маппит текущую строку {@link ResultSet} в доменную модель {@link Test_AttemptModel}.
     *
     * @param rs ResultSet, позиционированный на нужной строке
     * @return объект доменной модели
     * @throws SQLException при ошибках доступа к ResultSet
     */
    private Test_AttemptModel mapRowToTestAttempt(ResultSet rs) throws SQLException {
        return new Test_AttemptModel(
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