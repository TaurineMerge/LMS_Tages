package com.example.lms.answer.api.domain.repository;

import com.example.lms.answer.api.domain.model.AnswerModel;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с ответами (ANSWER_D).
 *
 * Отвечает за доступ к данным:
 *  - чтение / сохранение / обновление / удаление ответов;
 *  - выборка ответов по вопросу;
 *  - подсчёт количества ответов.
 *
 * Принятые договорённости по доменной логике:
 *  - "правильный" ответ — это ответ с score > 0;
 *  - "неправильный" — ответ с score = 0.
 */
public interface AnswerRepositoryInterface {

    /**
     * Сохранить новый ответ.
     *
     * @param answer доменная модель нового ответа (без id или с временным id)
     * @return сохранённая модель (обычно уже с присвоенным id)
     */
    AnswerModel save(AnswerModel answer);

    /**
     * Обновить существующий ответ.
     *
     * @param answer модель с актуальными данными (id должен быть заполнен)
     * @return обновлённая модель
     */
    AnswerModel update(AnswerModel answer);

    /**
     * Найти ответ по его идентификатору.
     *
     * @param id идентификатор ответа (ANSWER_D.id)
     * @return Optional с найденным ответом или пустой, если не найден
     */
    Optional<AnswerModel> findById(UUID id);

    /**
     * Получить все ответы из таблицы.
     *
     * @return список всех ответов
     */
    List<AnswerModel> findAll();

    /**
     * Найти все ответы для указанного вопроса.
     *
     * @param questionId идентификатор вопроса (ANSWER_D.question_id)
     * @return список ответов, привязанных к этому вопросу
     */
    List<AnswerModel> findByQuestionId(UUID questionId);

    /**
     * Найти правильные ответы для вопроса.
     *
     * Под "правильными" понимаются ответы с score > 0.
     *
     * @param questionId идентификатор вопроса
     * @return список правильных ответов
     */
    List<AnswerModel> findCorrectAnswersByQuestionId(UUID questionId);

    /**
     * Удалить ответ по его идентификатору.
     *
     * @param id идентификатор ответа
     * @return true, если запись была удалена; false, если такой записи не было
     */
    boolean deleteById(UUID id);

    /**
     * Удалить все ответы, привязанные к указанному вопросу.
     *
     * @param questionId идентификатор вопроса
     * @return количество удалённых записей
     */
    int deleteByQuestionId(UUID questionId);

    /**
     * Проверить существование ответа по его идентификатору.
     *
     * @param id идентификатор ответа
     * @return true, если ответ существует
     */
    boolean existsById(UUID id);

    /**
     * Проверить, существует ли уже ответ с таким текстом для данного вопроса.
     *
     * @param questionId идентификатор вопроса
     * @param text       текст ответа
     * @return true, если такой ответ уже есть
     */
    boolean existsByQuestionIdAndText(UUID questionId, String text);

    /**
     * Получить количество ответов для вопроса.
     *
     * @param questionId идентификатор вопроса
     * @return число ответов
     */
    int countByQuestionId(UUID questionId);

    /**
     * Получить количество правильных ответов для вопроса.
     *
     * Под "правильными" понимаются ответы с score > 0.
     *
     * @param questionId идентификатор вопроса
     * @return число правильных ответов
     */
    int countCorrectAnswersByQuestionId(UUID questionId);
}