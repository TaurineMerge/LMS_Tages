package com.example.lms.question.infrastructure.repositories;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.question.api.domain.model.QuestionModel;
import com.example.lms.question.domain.repository.QuestionRepositoryInterface;

import javax.sql.DataSource;
import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для управления вопросами тестов (QUESTION_D) с использованием
 * JDBC.
 * <p>
 * Предоставляет операции сохранения, обновления, удаления и выборки
 * вопросов, а также методы для работы с порядком отображения вопросов
 * внутри одного теста.
 * <p>
 * Особенности:
 * <ul>
 * <li>все операции работают напрямую через JDBC</li>
 * <li>выполняется логирование ключевых операций</li>
 * <li>таблица QUESTION_D содержит поле "order", поэтому используется
 * экранирование кавычками</li>
 * </ul>
 */
public class QuestionRepository implements QuestionRepositoryInterface {

	private static final Logger logger = LoggerFactory.getLogger(QuestionRepository.class);

	/** Источник соединений с БД. */
	private final DatabaseConfig dbConfig;

	// language=SQL
	private static final String INSERT_SQL = """
			INSERT INTO tests.question_d (test_id, text_of_question, "order")
			VALUES (?, ?, ?)
			RETURNING id
			""".trim();

	// language=SQL
	private static final String UPDATE_SQL = """
			UPDATE tests.question_d
			SET test_id = ?, text_of_question = ?, "order" = ?
			WHERE id = ?
			""".trim();

	// language=SQL
	private static final String SELECT_BY_ID = """
			SELECT id, test_id, text_of_question, "order"
			FROM tests.question_d
			WHERE id = ?
			""".trim();

	// language=SQL
	private static final String SELECT_ALL = """
			SELECT id, test_id, text_of_question, "order"
			FROM tests.question_d
			ORDER BY test_id, "order"
			""".trim();

	// language=SQL
	private static final String SELECT_BY_TEST = """
			SELECT id, test_id, text_of_question, "order"
			FROM tests.question_d
			WHERE test_id = ?
			ORDER BY "order"
			""".trim();

	private static final String DELETE_BY_ID = "DELETE FROM tests.question_d WHERE id = ?";

	private static final String DELETE_BY_TEST = "DELETE FROM tests.question_d WHERE test_id = ?";

	// language=SQL
	private static final String COUNT_BY_TEST = """
			SELECT COUNT(*)
			FROM tests.question_d
			WHERE test_id = ?
			""".trim();

	// language=SQL
	private static final String SEARCH_BY_TEXT = """
			SELECT id, test_id, text_of_question, "order"
			FROM tests.question_d
			WHERE LOWER(text_of_question) LIKE LOWER(?)
			ORDER BY test_id, "order"
			""".trim();

	// language=SQL
	private static final String MAX_ORDER_BY_TEST = """
			SELECT COALESCE(MAX("order"), -1)
			FROM tests.question_d
			WHERE test_id = ?
			""".trim();

	// language=SQL
	private static final String SHIFT_ORDER_SQL = """
			UPDATE tests.question_d
			SET "order" = "order" + ?
			WHERE test_id = ? AND "order" >= ?
			""".trim();

	/**
	 * Создаёт репозиторий вопросов.
	 *
	 * @param dbConfig источник соединений PostgreSQL
	 */
	public QuestionRepository(DatabaseConfig dbConfig) {
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

	static {
		try {
			Class.forName("org.postgresql.Driver");
			logger.info("PostgreSQL драйвер зарегистрирован");
		} catch (ClassNotFoundException e) {
			logger.error("Не удалось зарегистрировать драйвер PostgreSQL", e);
			throw new RuntimeException("Драйвер БД не найден", e);
		}
	}

	// ------------------------------------------------------------
	// CRUD + Query методы
	// ------------------------------------------------------------

	//
	// public QuestionRepository(DatabaseConfig dbConfig) {
	// //TODO Auto-generated constructor stub
	// }
	//
	/**
	 * Сохраняет новый вопрос в БД.
	 * <p>
	 * После вставки устанавливает сгенерированный {@code id} в модель.
	 *
	 * @param question модель вопроса
	 * @return сохранённую модель с заполненным ID
	 * @throws RuntimeException при ошибках SQL или отсутствии результата
	 */
	@Override
	public QuestionModel save(QuestionModel question) {
		logger.info("Сохранение нового вопроса для теста: {}", question.getTestId());

		question.validate();

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(INSERT_SQL)) {

			stmt.setObject(1, question.getTestId());
			stmt.setString(2, question.getTextOfQuestion());
			stmt.setInt(3, question.getOrder());

			ResultSet rs = stmt.executeQuery();
			if (rs.next()) {
				UUID id = rs.getObject("id", UUID.class);
				question.setId(id);
				logger.info("Вопрос сохранён: id={}, order={}", id, question.getOrder());
				return question;
			}

			throw new RuntimeException("Не удалось сохранить вопрос");

		} catch (SQLException e) {
			logger.error("Ошибка при сохранении вопроса", e);
			throw new RuntimeException("Ошибка базы данных", e);
		}
	}

	/**
	 * Обновляет существующий вопрос.
	 *
	 * @param question модель с обновлёнными данными
	 * @return обновлённую модель
	 * @throws IllegalArgumentException если вопрос не имеет ID
	 */
	@Override
	public QuestionModel update(QuestionModel question) {
		logger.info("Обновление вопроса id={}", question.getId());

		if (question.getId() == null) {
			throw new IllegalArgumentException("Вопрос должен иметь ID для обновления");
		}

		question.validate();

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(UPDATE_SQL)) {

			stmt.setObject(1, question.getTestId());
			stmt.setString(2, question.getTextOfQuestion());
			stmt.setInt(3, question.getOrder());
			stmt.setObject(4, question.getId());

			int updated = stmt.executeUpdate();
			if (updated == 0) {
				throw new RuntimeException("Вопрос с ID " + question.getId() + " не найден");
			}

			logger.info("Вопрос обновлён id={}", question.getId());
			return question;

		} catch (SQLException e) {
			logger.error("Ошибка при обновлении вопроса", e);
			throw new RuntimeException("Ошибка базы данных", e);
		}
	}

	/**
	 * Ищет вопрос по ID.
	 *
	 * @param id идентификатор вопроса
	 * @return Optional с найденным вопросом или пустой Optional
	 */
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

	/**
	 * Возвращает список всех вопросов в системе.
	 *
	 * @return список всех вопросов отсортированных по test_id и order
	 */
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

	/**
	 * Возвращает все вопросы указанного теста.
	 *
	 * @param testId ID теста
	 * @return список вопросов теста
	 */
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

	/**
	 * Удаляет вопрос по ID.
	 *
	 * @param id идентификатор вопроса
	 * @return true — если вопрос был удалён, false — если не найден
	 */
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

	/**
	 * Подсчитывает число вопросов в тесте.
	 *
	 * @param testId ID теста
	 * @return количество вопросов
	 */
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

	/**
	 * Ищет вопросы, содержащие указанный текст.
	 *
	 * @param text фрагмент текста для поиска
	 * @return список совпадающих вопросов
	 */
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

	/**
	 * Возвращает следующий порядковый номер для вопроса теста.
	 * <p>
	 * Если вопросов нет — возвращает 0.
	 *
	 * @param testId ID теста
	 * @return номер следующей позиции (max(order) + 1)
	 */
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

	/**
	 * Сдвигает порядок вопросов, начиная с указанной позиции.
	 * Используется при вставке нового вопроса внутрь списка.
	 *
	 * @param testId    ID теста
	 * @param fromOrder начиная с какого порядка сдвигать
	 * @param shiftBy   на сколько сдвигать (может быть отрицательным)
	 * @return количество обновлённых записей
	 */
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

	// ------------------------------------------------------------
	// Mapping
	// ------------------------------------------------------------

	/**
	 * Преобразует строку результата SQL в доменную модель QuestionModel.
	 *
	 * @param rs ResultSet, указывающий на текущую строку
	 * @return QuestionModel построенная из данных строки
	 * @throws SQLException если поля недоступны
	 */
	private QuestionModel mapRowToQuestion(ResultSet rs) throws SQLException {
		return new QuestionModel(
				rs.getObject("id", UUID.class),
				rs.getObject("test_id", UUID.class),
				rs.getString("text_of_question"),
				rs.getInt("order"));
	}
}