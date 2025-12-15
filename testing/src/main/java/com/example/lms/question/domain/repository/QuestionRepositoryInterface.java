package com.example.lms.question.domain.repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import com.example.lms.question.domain.model.QuestionModel;

/**
 * Интерфейс репозитория для работы с вопросами тестов
 * Соответствует таблице QUESTION_D
 */
public interface QuestionRepositoryInterface {
    
    /**
     * Сохранить новый вопрос
     * @param question модель вопроса
     * @return сохраненная модель с присвоенным ID
     */
    QuestionModel save(QuestionModel question);
    
    /**
     * Обновить существующий вопрос
     * @param question модель вопроса с обновленными данными
     * @return обновленная модель
     */
    QuestionModel update(QuestionModel question);
    
    /**
     * Найти вопрос по ID
     * @param id идентификатор вопроса
     * @return Optional с вопросом, если найден
     */
    Optional<QuestionModel> findById(UUID id);
    
    /**
     * Найти все вопросы
     * @return список всех вопросов
     */
    List<QuestionModel> findAll();
    
    /**
     * Найти вопросы по ID теста
     * @param testId идентификатор теста
     * @return список вопросов теста, отсортированный по порядку
     */
    List<QuestionModel> findByTestId(UUID testId);
    
    /**
     * Удалить вопрос по ID
     * @param id идентификатор вопроса
     * @return true если вопрос был удален
     */
    boolean deleteById(UUID id);
    
    
    /**
     * Получить количество вопросов в тесте
     * @param testId идентификатор теста
     * @return количество вопросов
     */
    int countByTestId(UUID testId);
    
    /**
     * Найти вопросы по тексту (регистронезависимый поиск)
     * @param text часть текста вопроса для поиска
     * @return список найденных вопросов
     */
    List<QuestionModel> findByTextContaining(String text);
    
    /**
     * Получить следующий порядковый номер для нового вопроса в тесте
     * @param testId идентификатор теста
     * @return следующий доступный порядковый номер
     */
    int getNextOrderForTest(UUID testId);
    
    /**
     * Обновить порядок вопросов в тесте
     * Сдвигает все вопросы начиная с указанного порядка
     * @param testId идентификатор теста
     * @param fromOrder порядок, с которого начинать сдвиг
     * @param shiftBy на сколько сдвигать (может быть отрицательным)
     * @return количество обновленных вопросов
     */
    int shiftQuestionsOrder(UUID testId, int fromOrder, int shiftBy);
}