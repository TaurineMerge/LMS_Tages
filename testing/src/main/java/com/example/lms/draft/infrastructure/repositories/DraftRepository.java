package com.example.lms.draft.infrastructure.repositories;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Types;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.draft.domain.model.DraftModel;
import com.example.lms.draft.domain.repository.DraftRepositoryInterface;

/**
 * Репозиторий для работы с черновиками тестов в базе данных.
 * Реализует интерфейс {@link DraftRepositoryInterface} для операций CRUD и поиска черновиков.
 *
 * Соответствует таблице draft_b:
 *  - id          UUID    (PK, not null, unique)
 *  - title       VARCHAR
 *  - min_point   INT
 *  - description TEXT
 *  - test_id     UUID    (может быть null)
 *  - course_id   UUID    (может быть null)
 */
public class DraftRepository implements DraftRepositoryInterface {

	private final DatabaseConfig dbConfig;

	/**
	 * Создаёт репозиторий с настройками подключения к БД.
	 *
	 * @param dbConfig конфигурация базы данных
	 */
	public DraftRepository(DatabaseConfig dbConfig) {
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
	 * Создаёт новый черновик в базе данных.
	 *
	 * @param draftModel объект {@link DraftModel} для сохранения
	 * @return сохранённый объект {@link DraftModel} с присвоенным идентификатором
	 * @throws RuntimeException если произошла SQL-ошибка или вставка не вернула id
	 */
	@Override
	public DraftModel create(DraftModel draftModel) {
		draftModel.validate();

		String sql = """
				INSERT INTO testing.draft_b (title, min_point, description, test_id, course_id)
				VALUES (?, ?, ?, ?, ?)
				RETURNING id
				""";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setString(1, draftModel.getTitle());

			if (draftModel.getMinPoint() != null)
				stmt.setInt(2, draftModel.getMinPoint());
			else
				stmt.setNull(2, Types.INTEGER);

			stmt.setString(3, draftModel.getDescription());
			stmt.setObject(4, draftModel.getTestId());
			stmt.setObject(5, draftModel.getCourseId());

			try (ResultSet rs = stmt.executeQuery()) {
				if (rs.next()) {
					draftModel.setId(rs.getObject("id", UUID.class));
					return draftModel;
				}
			}

			throw new RuntimeException("Не удалось сохранить черновик (draft)");

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при сохранении черновика (draft)", e);
		}
	}

	/**
	 * Находит черновик по идентификатору.
	 *
	 * @param id UUID черновика
	 * @return Optional с черновиком или пустой Optional, если не найден
	 */
	@Override
	public Optional<DraftModel> findById(UUID id) {
		String sql = """
				SELECT id, title, min_point, description, test_id, course_id
				FROM testing.draft_b
				WHERE id = ?
				""";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setObject(1, id);

			try (ResultSet rs = stmt.executeQuery()) {
				if (rs.next()) {
					return Optional.of(mapRowToDraft(rs));
				}
			}

			return Optional.empty();

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при поиске черновика по id: " + id, e);
		}
	}

	/**
	 * Находит черновик по идентификатору теста.
	 *
	 * @param testId UUID теста
	 * @return Optional с черновиком или пустой Optional, если не найден
	 */
	@Override
	public Optional<DraftModel> findByTestId(UUID testId) {
		String sql = """
				SELECT id, title, min_point, description, test_id, course_id
				FROM testing.draft_b
				WHERE test_id = ?
				""";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setObject(1, testId);

			try (ResultSet rs = stmt.executeQuery()) {
				if (rs.next()) {
					return Optional.of(mapRowToDraft(rs));
				}
			}

			return Optional.empty();

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при поиске черновика по test_id: " + testId, e);
		}
	}

	/**
	 * Находит черновики по идентификатору курса.
	 *
	 * @param courseId UUID курса
	 * @return список черновиков для указанного курса
	 */
	@Override
	public List<DraftModel> findByCourseId(UUID courseId) {
		String sql = """
				SELECT id, title, min_point, description, test_id, course_id
				FROM testing.draft_b
				WHERE course_id = ?
				ORDER BY title
				""";

		List<DraftModel> result = new ArrayList<>();

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setObject(1, courseId);

			try (ResultSet rs = stmt.executeQuery()) {
				while (rs.next()) {
					result.add(mapRowToDraft(rs));
				}
			}

			return result;

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при поиске черновиков по course_id: " + courseId, e);
		}
	}

	/**
	 * Возвращает список всех черновиков.
	 *
	 * @return список черновиков
	 */
	@Override
	public List<DraftModel> findAll() {
		String sql = """
				SELECT id, title, min_point, description, test_id, course_id
				FROM testing.draft_b
				ORDER BY title
				""";

		List<DraftModel> result = new ArrayList<>();

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql);
				ResultSet rs = stmt.executeQuery()) {

			while (rs.next()) {
				result.add(mapRowToDraft(rs));
			}

			return result;

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при получении списка черновиков", e);
		}
	}

	/**
	 * Обновляет существующий черновик.
	 *
	 * @param draftModel объект {@link DraftModel} с заполненным id
	 * @return обновлённый объект {@link DraftModel}
	 */
	@Override
	public DraftModel update(DraftModel draftModel) {
		draftModel.validate();

		if (draftModel.getId() == null) {
			throw new IllegalArgumentException("Нельзя обновить черновик без id");
		}

		String sql = """
				UPDATE testing.draft_b
				SET title = ?, min_point = ?, description = ?, test_id = ?, course_id = ?
				WHERE id = ?
				""";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setString(1, draftModel.getTitle());

			if (draftModel.getMinPoint() != null)
				stmt.setInt(2, draftModel.getMinPoint());
			else
				stmt.setNull(2, Types.INTEGER);

			stmt.setString(3, draftModel.getDescription());
			stmt.setObject(4, draftModel.getTestId());
			stmt.setObject(5, draftModel.getCourseId());
			stmt.setObject(6, draftModel.getId());

			int updated = stmt.executeUpdate();
			if (updated == 0) {
				throw new RuntimeException("Черновик с ID " + draftModel.getId() + " не найден");
			}

			return draftModel;

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при обновлении черновика (draft)", e);
		}
	}

	/**
	 * Удаляет черновик по идентификатору.
	 *
	 * @param id UUID черновика
	 * @return true если удалено, false если записи не было
	 */
	@Override
	public boolean deleteById(UUID id) {
		Connection conn = null;
		
		try {
			conn = getConnection();
			conn.setAutoCommit(false);
			
			// 1. Находим все вопросы черновика
			String findQuestionsSQL = "SELECT id FROM testing.question_d WHERE draft_id = ?";
			List<UUID> questionIds = new ArrayList<>();
			
			try (PreparedStatement findQuestionsStmt = conn.prepareStatement(findQuestionsSQL)) {
				findQuestionsStmt.setObject(1, id);
				ResultSet rs = findQuestionsStmt.executeQuery();
				
				while (rs.next()) {
					questionIds.add(rs.getObject("id", UUID.class));
				}
			}
			
			System.out.println("Found " + questionIds.size() + " draft questions to delete");
			
			// 2. Удаляем ответы для этих вопросов
			if (!questionIds.isEmpty()) {
				String deleteAnswersSQL = "DELETE FROM testing.answer_d WHERE question_id = ?";
				try (PreparedStatement deleteAnswersStmt = conn.prepareStatement(deleteAnswersSQL)) {
					for (UUID questionId : questionIds) {
						deleteAnswersStmt.setObject(1, questionId);
						deleteAnswersStmt.addBatch();
					}
					int[] answerDeletions = deleteAnswersStmt.executeBatch();
					int totalAnswersDeleted = 0;
					for (int deleted : answerDeletions) {
						totalAnswersDeleted += deleted;
					}
					System.out.println("Deleted " + totalAnswersDeleted + " draft answers");
				}
			}
			
			// 3. Удаляем вопросы черновика
			String deleteQuestionsSQL = "DELETE FROM testing.question_d WHERE draft_id = ?";
			try (PreparedStatement stmt = conn.prepareStatement(deleteQuestionsSQL)) {
				stmt.setObject(1, id);
				int questionsDeleted = stmt.executeUpdate();
				System.out.println("Deleted " + questionsDeleted + " draft questions");
			}
			
			// 4. Удаляем сам черновик - ИСПРАВЛЕНО!
			String deleteDraftSQL = "DELETE FROM testing.draft_b WHERE id = ?";
			
			try (PreparedStatement stmt = conn.prepareStatement(deleteDraftSQL)) {
				stmt.setObject(1, id);
				int deleted = stmt.executeUpdate();
				
				conn.commit();
				System.out.println("Deleted " + deleted + " drafts");
				return deleted > 0;
			}
			
		} catch (SQLException e) {
			if (conn != null) {
				try { conn.rollback(); } catch (SQLException ignored) {}
			}
			throw new RuntimeException("Ошибка при удалении черновика: " + e.getMessage(), e);
		} finally {
			if (conn != null) {
				try {
					conn.setAutoCommit(true);
					conn.close();
				} catch (SQLException ignored) {}
			}
		}
	}
	/**
	 * Удаляет черновики по идентификатору курса.
	 *
	 * @param courseId UUID курса
	 * @return true если были удалены записи, false если записей не было
	 */
	@Override
	public boolean deleteByCourseId(UUID courseId) {
		String sql = """
				DELETE FROM testing.draft_b
				WHERE course_id = ?
				""";

		try (Connection conn = getConnection();
				PreparedStatement stmt = conn.prepareStatement(sql)) {

			stmt.setObject(1, courseId);
			return stmt.executeUpdate() > 0;

		} catch (SQLException e) {
			throw new RuntimeException("Ошибка при удалении черновиков по course_id: " + courseId, e);
		}
	}

	/**
	 * Маппит строку ResultSet в {@link DraftModel}.
	 *
	 * @param rs ResultSet с данными черновика
	 * @return объект {@link DraftModel}
	 * @throws SQLException если произошла ошибка при чтении данных
	 */
	private DraftModel mapRowToDraft(ResultSet rs) throws SQLException {
		return new DraftModel(
				rs.getObject("id", UUID.class),
				rs.getObject("test_id", UUID.class),
				rs.getObject("course_id", UUID.class),
				rs.getString("title"),
				rs.getObject("min_point", Integer.class),
				rs.getString("description"));
	}
}