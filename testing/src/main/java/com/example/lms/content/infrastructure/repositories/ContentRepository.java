package com.example.lms.content.infrastructure.repositories;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.content.domain.model.ContentModel;
import com.example.lms.content.domain.repository.ContentRepositoryInterface;

import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с контентом в базе данных.
 * Реализует интерфейс ContentRepositoryInterface для операций CRUD.
 */
public class ContentRepository implements ContentRepositoryInterface {

    private final DatabaseConfig dbConfig;

    /**
     * Конструктор репозитория.
     *
     * @param dbConfig конфигурация базы данных, содержащая параметры подключения
     */
    public ContentRepository(DatabaseConfig dbConfig) {
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
     * Сохраняет новый элемент контента в базе данных.
     *
     * @param content объект ContentModel для сохранения
     * @return сохраненный объект ContentModel с присвоенным идентификатором
     * @throws RuntimeException если не удалось сохранить или произошла SQL-ошибка
     */
    @Override
    public ContentModel save(ContentModel content) {
        content.validate();

        String sql = """
                INSERT INTO testing.content_d ("order", content, type_of_content, question_id, answer_id)
                VALUES (?, ?, ?, ?, ?)
                RETURNING id
                """.trim();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setInt(1, content.getOrder());
            stmt.setString(2, content.getContent());

            if (content.getTypeOfContent() != null) {
                stmt.setBoolean(3, content.getTypeOfContent());
            } else {
                stmt.setNull(3, Types.BOOLEAN);
            }

            if (content.getQuestionId() != null) {
                stmt.setObject(4, content.getQuestionId());
            } else {
                stmt.setNull(4, Types.OTHER);
            }

            if (content.getAnswerId() != null) {
                stmt.setObject(5, content.getAnswerId());
            } else {
                stmt.setNull(5, Types.OTHER);
            }

            ResultSet rs = stmt.executeQuery();
            if (rs.next()) {
                content.setId(rs.getObject("id", UUID.class));
                return content;
            }

            throw new RuntimeException("Не удалось сохранить контент");

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при сохранении контента", e);
        }
    }

    /**
     * Обновляет существующий элемент контента в базе данных.
     *
     * @param content объект ContentModel с обновленными данными
     * @return обновленный объект ContentModel
     * @throws IllegalArgumentException если контент не имеет идентификатора
     * @throws RuntimeException         если контент не найден или произошла SQL-ошибка
     */
    @Override
    public ContentModel update(ContentModel content) {
        if (content.getId() == null) {
            throw new IllegalArgumentException("Контент должен иметь ID");
        }

        content.validate();

        String sql = """
                UPDATE testing.content_d
                SET "order" = ?, content = ?, type_of_content = ?, question_id = ?, answer_id = ?
                WHERE id = ?
                """;

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setInt(1, content.getOrder());
            stmt.setString(2, content.getContent());

            if (content.getTypeOfContent() != null) {
                stmt.setBoolean(3, content.getTypeOfContent());
            } else {
                stmt.setNull(3, Types.BOOLEAN);
            }

            if (content.getQuestionId() != null) {
                stmt.setObject(4, content.getQuestionId());
            } else {
                stmt.setNull(4, Types.OTHER);
            }

            if (content.getAnswerId() != null) {
                stmt.setObject(5, content.getAnswerId());
            } else {
                stmt.setNull(5, Types.OTHER);
            }

            stmt.setObject(6, content.getId());

            int updated = stmt.executeUpdate();
            if (updated == 0) {
                throw new RuntimeException("Контент с ID " + content.getId() + " не найден");
            }

            return content;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при обновлении контента", e);
        }
    }

    /**
     * Находит элемент контента по его идентификатору.
     *
     * @param id уникальный идентификатор контента
     * @return Optional с найденным контентом или пустой Optional, если не найден
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public Optional<ContentModel> findById(UUID id) {
        String sql = """
                SELECT id, "order", content, type_of_content, question_id, answer_id
                FROM testing.content_d
                WHERE id = ?
                """;

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            ResultSet rs = stmt.executeQuery();

            if (rs.next()) {
                return Optional.of(mapRowToContent(rs));
            }

            return Optional.empty();

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске контента по ID", e);
        }
    }

    /**
     * Получает все элементы контента из базы данных.
     *
     * @return список всех элементов контента, отсортированных по порядку
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<ContentModel> findAll() {
        String sql = """
                SELECT id, "order", content, type_of_content, question_id, answer_id
                FROM testing.content_d
                ORDER BY "order"
                """;

        List<ContentModel> contents = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql);
                ResultSet rs = stmt.executeQuery()) {

            while (rs.next()) {
                contents.add(mapRowToContent(rs));
            }

            return contents;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при получении всего контента", e);
        }
    }

    /**
     * Удаляет элемент контента по его идентификатору.
     *
     * @param id уникальный идентификатор контента для удаления
     * @return true если контент был удален, false если контент не найден
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public boolean deleteById(UUID id) {
        String sql = "DELETE FROM testing.content_d WHERE id = ?";

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            return stmt.executeUpdate() > 0;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при удалении контента", e);
        }
    }

    /**
     * Проверяет существование элемента контента по идентификатору.
     *
     * @param id уникальный идентификатор контента
     * @return true если контент существует, false если нет
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public boolean existsById(UUID id) {
        String sql = "SELECT 1 FROM testing.content_d WHERE id = ?";

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, id);
            return stmt.executeQuery().next();

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при проверке существования контента", e);
        }
    }

    /**
     * Находит элементы контента по содержимому (частичное совпадение, без учета регистра).
     *
     * @param content часть содержимого для поиска
     * @return список элементов контента, содержащих указанную строку в содержимом
     * @throws RuntimeException если произошла SQL-ошибка
     */
    @Override
    public List<ContentModel> findByContentContaining(String content) {
        String sql = """
                SELECT id, "order", content, type_of_content, question_id, answer_id
                FROM testing.content_d
                WHERE LOWER(content) LIKE LOWER(?)
                ORDER BY "order"
                """;

        List<ContentModel> contents = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setString(1, "%" + content + "%");
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                contents.add(mapRowToContent(rs));
            }

            return contents;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске контента по содержимому", e);
        }
    }

    /**
     * Находит элементы контента по типу контента.
     *
     * @param typeOfContent тип контента для поиска
     * @return список элементов контента с указанным типом
     * @throws RuntimeException если произошла SQL-ошибка
     */
    public List<ContentModel> findByTypeOfContent(Boolean typeOfContent) {
        String sql = """
                SELECT id, "order", content, type_of_content, question_id, answer_id
                FROM testing.content_d
                WHERE type_of_content = ?
                ORDER BY "order"
                """;

        List<ContentModel> contents = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            if (typeOfContent != null) {
                stmt.setBoolean(1, typeOfContent);
            } else {
                stmt.setNull(1, Types.BOOLEAN);
            }

            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                contents.add(mapRowToContent(rs));
            }

            return contents;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске контента по типу", e);
        }
    }

    /**
     * Находит элементы контента по ID вопроса.
     *
     * @param questionId ID вопроса
     * @return список элементов контента, связанных с указанным вопросом
     * @throws RuntimeException если произошла SQL-ошибка
     */
    public List<ContentModel> findByQuestionId(UUID questionId) {
        String sql = """
                SELECT id, "order", content, type_of_content, question_id, answer_id
                FROM testing.content_d
                WHERE question_id = ?
                ORDER BY "order"
                """;

        List<ContentModel> contents = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, questionId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                contents.add(mapRowToContent(rs));
            }

            return contents;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске контента по question_id", e);
        }
    }

    /**
     * Находит элементы контента по ID ответа.
     *
     * @param answerId ID ответа
     * @return список элементов контента, связанных с указанным ответом
     * @throws RuntimeException если произошла SQL-ошибка
     */
    public List<ContentModel> findByAnswerId(UUID answerId) {
        String sql = """
                SELECT id, "order", content, type_of_content, question_id, answer_id
                FROM testing.content_d
                WHERE answer_id = ?
                ORDER BY "order"
                """;

        List<ContentModel> contents = new ArrayList<>();

        try (Connection conn = getConnection();
                PreparedStatement stmt = conn.prepareStatement(sql)) {

            stmt.setObject(1, answerId);
            ResultSet rs = stmt.executeQuery();

            while (rs.next()) {
                contents.add(mapRowToContent(rs));
            }

            return contents;

        } catch (SQLException e) {
            throw new RuntimeException("Ошибка при поиске контента по answer_id", e);
        }
    }

    /**
     * Преобразует строку ResultSet в объект ContentModel.
     *
     * @param rs ResultSet, содержащий данные контента
     * @return объект ContentModel, созданный из данных ResultSet
     * @throws SQLException если произошла ошибка при чтении данных из ResultSet
     */
    private ContentModel mapRowToContent(ResultSet rs) throws SQLException {
        return new ContentModel(
                rs.getObject("id", UUID.class),
                rs.getInt("order"),
                rs.getString("content"),
                rs.getObject("type_of_content", Boolean.class),
                rs.getObject("question_id", UUID.class),
                rs.getObject("answer_id", UUID.class));
    }
}