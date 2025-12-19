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
 * Репозиторий для работы с тестами в базе данных.
 * Реализует интерфейс TestRepositoryInterface для операций CRUD и поиска
 * тестов.
 */
public class TestRepository implements TestRepositoryInterface {

	private final DatabaseConfig dbConfig;

	/**
	 * Конструктор репозитория.
	 *
	 * @param dbConfig конфигурация базы данных, содержащая параметры подключения
	 */
	public TestRepository(DatabaseConfig dbConfig) {
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
	 * Сохраняет новый тест в базе данных.
	 *
	 * @param test объект TestModel для сохранения
	 * @return сохраненный объект TestModel с присвоенным идентификатором
	 * @throws RuntimeException если не удалось сохранить тест или произошла
	 *                          SQL-ошибка
	 */
	@Override
	public TestModel save(TestModel test) {
		test.validate();

		String sql = """
				INSERT INTO testing.test_d (course_id, title, min_point, description)
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

	/**
	 * Обновляет существующий тест в базе данных.
	 *
	 * @param test объект TestModel с обновленными данными
	 * @return обновленный объект TestModel
	 * @throws IllegalArgumentException если тест не имеет идентификатора
	 * @throws RuntimeException         если тест не найден или произошла SQL-ошибка
	 */
	@Override
	public TestModel update(TestModel test) {
		if (test.getId() == null) {
			throw new IllegalArgumentException("Тест должен иметь ID");
		}

		test.validate();

		String sql = """
				UPDATE testing.test_d
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

	/**
	 * Находит тест по его идентификатору.
	 *
	 * @param id уникальный идентификатор теста
	 * @return Optional с найденным тестом или пустой Optional, если тест не найден
	 * @throws RuntimeException если произошла SQL-ошибка
	 */
	@Override
	public Optional<TestModel> findById(UUID id) {
		String sql = """
				SELECT id, course_id, title, min_point, description
				FROM testing.test_d
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

	/**
	 * Получает все тесты из базы данных.
	 *
	 * @return список всех тестов, отсортированных по названию
	 * @throws RuntimeException если произошла SQL-ошибка
	 */
	@Override
	public List<TestModel> findAll() {
		String sql = """
				SELECT id, course_id, title, min_point, description
				FROM testing.test_d
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

	/**
	 * Удаляет тест по его идентификатору.
	 *
	 * @param id уникальный идентификатор теста для удаления
	 * @return true если тест был удален, false если тест не найден
	 * @throws RuntimeException если произошла SQL-ошибка
	 */
	@Override
	public boolean deleteById(UUID id) {
		String sql = "DELETE FROM testing.test_d WHERE id = ?";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setObject(1, id);
			return stmt.executeUpdate() > 0;

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при удалении теста", e);
		}
	}

	/**
	 * Проверяет существование теста по идентификатору.
	 *
	 * @param id уникальный идентификатор теста
	 * @return true если тест существует, false если нет
	 * @throws RuntimeException если произошла SQL-ошибка
	 */
	@Override
	public boolean existsById(UUID id) {
		String sql = "SELECT 1 FROM testing.test_d WHERE id = ?";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setObject(1, id);
			return stmt.executeQuery().next();

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при проверке существования теста", e);
		}
	}

	/**
	 * Ищет тесты по частичному совпадению названия (без учета регистра).
	 *
	 * @param title часть названия для поиска
	 * @return список тестов, содержащих указанную строку в названии
	 * @throws RuntimeException если произошла SQL-ошибка
	 */
	@Override
	public List<TestModel> findByTitleContaining(String title) {
		String sql = """
				SELECT id, course_id, title, min_point, description
				FROM testing.test_d
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

	/**
	 * Подсчитывает количество тестов для указанного курса.
	 *
	 * @param courseId уникальный идентификатор курса
	 * @return количество тестов в указанном курсе
	 * @throws RuntimeException если произошла SQL-ошибка
	 */
	@Override
	public int countByCourseId(UUID courseId) {
		String sql = "SELECT COUNT(*) FROM testing.test_d WHERE course_id = ?";

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
				FROM testing.test_d
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

	/**
	 * Преобразует строку ResultSet в объект TestModel.
	 *
	 * @param rs ResultSet, содержащий данные теста
	 * @return объект TestModel, созданный из данных ResultSet
	 * @throws SQLException если произошла ошибка при чтении данных из ResultSet
	 */
	private TestModel mapRowToTest(ResultSet rs) throws SQLException {
		return new TestModel(
				rs.getObject("id", UUID.class),
				rs.getObject("course_id", UUID.class),
				rs.getString("title"),
				rs.getObject("min_point", Integer.class),
				rs.getString("description"));
	}
}
