package com.example.lms.content.domain.repository;

import com.example.lms.content.domain.model.ContentModel;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Интерфейс репозитория для работы с элементами контента.
 *
 * Соответствует таблице content_d:
 *  - id              UUID    (PK, not null)
 *  - order           INT     — порядковый номер в курсе
 *  - content         VARCHAR — текст контента/вопроса
 *  - type_of_content BOOLEAN — тип контента (true/false для различных типов)
 *  - question_id     UUID    — ссылка на вопрос (может быть null)
 *  - answer_id       UUID    — ссылка на ответ (может быть null)
 *
 * Отвечает за доступ к данным:
 *  - создание / обновление / удаление элементов контента;
 *  - выборка по идентификатору и различным критериям;
 *  - поиск по содержимому и типу контента.
 */
public interface ContentRepositoryInterface {

    /**
     * Сохранить новый элемент контента.
     *
     * @param content доменная модель элемента контента (обычно без id до сохранения)
     * @return сохранённая модель с присвоенным идентификатором
     */
    ContentModel save(ContentModel content);

    /**
     * Обновить существующий элемент контента.
     *
     * @param content модель элемента контента с актуальными данными (id должен быть заполнен)
     * @return обновлённая модель
     */
    ContentModel update(ContentModel content);

    /**
     * Найти элемент контента по его идентификатору.
     *
     * @param id идентификатор элемента контента (content_d.id)
     * @return Optional с элементом контента, если найден
     */
    Optional<ContentModel> findById(UUID id);

    /**
     * Получить список всех элементов контента.
     *
     * @return список всех элементов контента, отсортированных по порядку
     */
    List<ContentModel> findAll();

    /**
     * Найти элементы контента по содержимому (частичное совпадение, без учета регистра).
     *
     * @param content часть содержимого для поиска
     * @return список элементов контента, содержащих указанную строку в содержимом
     */
    List<ContentModel> findByContentContaining(String content);

    /**
     * Найти элементы контента по типу контента.
     *
     * @param typeOfContent тип контента для поиска
     * @return список элементов контента с указанным типом
     */
    List<ContentModel> findByTypeOfContent(Boolean typeOfContent);

    /**
     * Найти элементы контента по ID вопроса.
     *
     * @param questionId ID вопроса
     * @return список элементов контента, связанных с указанным вопросом
     */
    List<ContentModel> findByQuestionId(UUID questionId);

    /**
     * Найти элементы контента по ID ответа.
     *
     * @param answerId ID ответа
     * @return список элементов контента, связанных с указанным ответом
     */
    List<ContentModel> findByAnswerId(UUID answerId);

    /**
     * Удалить элемент контента по его идентификатору.
     *
     * @param id идентификатор элемента контента
     * @return true, если элемент контента был удалён; false — если запись не найдена
     */
    boolean deleteById(UUID id);

    /**
     * Проверить существование элемента контента по его идентификатору.
     *
     * @param id идентификатор элемента контента
     * @return true, если элемент контента существует
     */
    boolean existsById(UUID id);
}