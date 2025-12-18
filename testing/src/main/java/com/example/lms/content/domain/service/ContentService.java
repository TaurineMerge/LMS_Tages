package com.example.lms.content.domain.service;

import java.util.List;
import java.util.UUID;
import java.util.NoSuchElementException;

import com.example.lms.content.api.dto.Content;
import com.example.lms.content.domain.model.ContentModel;
import com.example.lms.content.domain.repository.ContentRepositoryInterface;

/**
 * Сервис для работы с элементами контента.
 * <p>
 * Отвечает за:
 * <ul>
 * <li>конвертацию DTO ↔ Domain Model</li>
 * <li>вызовы репозитория ({@link ContentRepositoryInterface})</li>
 * <li>бизнес-логику CRUD-операций над элементами контента</li>
 * <li>предоставление дополнительных методов поиска и фильтрации</li>
 * </ul>
 *
 * Сервисный слой отделяет контроллеры от репозитория
 * и обеспечивает единое место для обработки бизнес-процессов.
 */
public class ContentService {

    /** Репозиторий элементов контента (слой работы с базой данных). */
    private final ContentRepositoryInterface repository;

    /**
     * Создаёт сервис элементов контента.
     *
     * @param repository репозиторий, выполняющий операции с БД
     */
    public ContentService(ContentRepositoryInterface repository) {
        this.repository = repository;
    }

    // ---------------------------------------------------------------------
    // DTO -> MODEL
    // ---------------------------------------------------------------------

    /**
     * Преобразует DTO в доменную модель {@link ContentModel}.
     * <p>
     * Доменная модель используется на сервисном и репозиторном уровнях.
     * DTO применяется только для API и не содержит бизнес-логики.
     *
     * @param dto объект API DTO
     * @return доменная модель элемента контента
     */
    private ContentModel toModel(Content dto) {
        return new ContentModel(
                dto.getId() != null ? UUID.fromString(dto.getId()) : null,
                dto.getOrder(),
                dto.getContent(),
                dto.getTypeOfContent(),
                dto.getQuestionId() != null ? UUID.fromString(dto.getQuestionId()) : null,
                dto.getAnswerId() != null ? UUID.fromString(dto.getAnswerId()) : null);
    }

    // ---------------------------------------------------------------------
    // MODEL → DTO
    // ---------------------------------------------------------------------

    /**
     * Преобразует доменную модель элемента контента в DTO для API.
     *
     * @param model доменная модель элемента контента
     * @return DTO, отправляемый наружу
     */
    private Content toDto(ContentModel model) {
        return new Content(
                model.getId() != null ? model.getId().toString() : null,
                model.getOrder(),
                model.getContent(),
                model.getTypeOfContent(),
                model.getQuestionId() != null ? model.getQuestionId().toString() : null,
                model.getAnswerId() != null ? model.getAnswerId().toString() : null);
    }

    // ---------------------------------------------------------------------
    // PUBLIC API METHODS (CRUD)
    // ---------------------------------------------------------------------

    /**
     * Возвращает список всех элементов контента, конвертируя их в DTO.
     *
     * @return список элементов контента в формате DTO
     */
    public List<Content> getAllContents() {
        return repository.findAll().stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Создаёт новый элемент контента.
     *
     * @param dto входные данные элемента контента
     * @return созданный элемент контента в виде DTO
     * @throws IllegalArgumentException если данные невалидны
     */
    public Content createContent(Content dto) {
        ContentModel model = toModel(dto);
        ContentModel saved = repository.save(model);
        return toDto(saved);
    }

    /**
     * Получает элемент контента по ID.
     *
     * @param id строковый UUID элемента контента
     * @return DTO найденного элемента контента
     * @throws NoSuchElementException если элемент контента не найден
     */
    public Content getContentById(String id) {
        UUID uuid = UUID.fromString(id);
        ContentModel model = repository.findById(uuid)
                .orElseThrow(() -> new NoSuchElementException("Элемент контента с ID " + id + " не найден"));
        return toDto(model);
    }

    /**
     * Обновляет существующий элемент контента.
     *
     * @param dto данные элемента контента с актуальными полями
     * @return DTO обновлённого элемента контента
     * @throws NoSuchElementException если элемент контента не найден
     * @throws IllegalArgumentException если данные невалидны
     */
    public Content updateContent(Content dto) {
        if (dto.getId() == null) {
            throw new IllegalArgumentException("ID элемента контента не может быть null для обновления");
        }
        
        // Проверяем существование элемента
        UUID uuid = UUID.fromString(dto.getId());
        if (!repository.existsById(uuid)) {
            throw new NoSuchElementException("Элемент контента с ID " + dto.getId() + " не найден");
        }
        
        ContentModel model = toModel(dto);
        ContentModel updated = repository.update(model);
        return toDto(updated);
    }

    /**
     * Удаляет элемент контента по ID.
     *
     * @param id строковый UUID элемента контента
     * @return true — если элемент контента был удалён; false — если не найден
     */
    public boolean deleteContent(String id) {
        UUID uuid = UUID.fromString(id);
        return repository.deleteById(uuid);
    }

    /**
     * Проверяет существование элемента контента по ID.
     *
     * @param id строковый UUID элемента контента
     * @return true — если элемент контента существует; false — если нет
     */
    public boolean existsById(String id) {
        UUID uuid = UUID.fromString(id);
        return repository.existsById(uuid);
    }

    // ---------------------------------------------------------------------
    // PUBLIC API METHODS (SEARCH AND FILTER)
    // ---------------------------------------------------------------------

    /**
     * Находит элементы контента по содержимому (частичное совпадение, без учета регистра).
     *
     * @param content часть содержимого для поиска
     * @return список элементов контента, содержащих указанную строку в содержимом
     */
    public List<Content> findByContentContaining(String content) {
        return repository.findByContentContaining(content).stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Находит элементы контента по типу контента.
     *
     * @param typeOfContent тип контента для поиска
     * @return список элементов контента с указанным типом
     */
    public List<Content> findByTypeOfContent(Boolean typeOfContent) {
        return repository.findByTypeOfContent(typeOfContent).stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Находит элементы контента по ID вопроса.
     *
     * @param questionId строковый UUID вопроса
     * @return список элементов контента, связанных с указанным вопросом
     */
    public List<Content> findByQuestionId(String questionId) {
        UUID uuid = UUID.fromString(questionId);
        return repository.findByQuestionId(uuid).stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Находит элементы контента по ID ответа.
     *
     * @param answerId строковый UUID ответа
     * @return список элементов контента, связанных с указанным ответом
     */
    public List<Content> findByAnswerId(String answerId) {
        UUID uuid = UUID.fromString(answerId);
        return repository.findByAnswerId(uuid).stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Находит элементы контента по содержимому и типу контента.
     *
     * @param content часть содержимого для поиска
     * @param typeOfContent тип контента для поиска
     * @return список элементов контента, удовлетворяющих обоим критериям
     */
    public List<Content> findByContentAndType(String content, Boolean typeOfContent) {
        // Фильтруем результаты на стороне сервиса для простоты
        List<Content> byContent = findByContentContaining(content);
        if (typeOfContent == null) {
            return byContent;
        }
        
        return byContent.stream()
                .filter(c -> typeOfContent.equals(c.getTypeOfContent()))
                .toList();
    }

    /**
     * Получает элементы контента с максимальным порядковым номером.
     *
     * @return список элементов контента с максимальным порядковым номером
     */
    public List<Content> findWithMaxOrder() {
        List<Content> allContents = getAllContents();
        if (allContents.isEmpty()) {
            return List.of();
        }
        
        int maxOrder = allContents.stream()
                .mapToInt(Content::getOrder)
                .max()
                .orElse(0);
        
        return allContents.stream()
                .filter(c -> c.getOrder() == maxOrder)
                .toList();
    }

    /**
     * Проверяет валидность данных элемента контента перед сохранением.
     *
     * @param dto данные для проверки
     * @throws IllegalArgumentException если данные невалидны
     */
    public void validateContentData(Content dto) {
        if (dto.getOrder() == null || dto.getOrder() < 0) {
            throw new IllegalArgumentException("Order cannot be null and must be non-negative");
        }
        if (dto.getContent() == null || dto.getContent().trim().isEmpty()) {
            throw new IllegalArgumentException("Content cannot be empty");
        }
    }
}