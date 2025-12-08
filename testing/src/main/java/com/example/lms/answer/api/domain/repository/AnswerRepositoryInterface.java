package com.example.lms.answer.api.domain.repository;

import com.example.lms.answer.api.domain.model.AnswerModel;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Интерфейс репозитория для работы с ответами
 */
public interface AnswerRepositoryInterface {
    
    /**
     * Сохранить новый ответ
     */
    AnswerModel save(AnswerModel answer);
    
    /**
     * Обновить существующий ответ
     */
    AnswerModel update(AnswerModel answer);
    
    /**
     * Найти ответ по ID
     */
    Optional<AnswerModel> findById(UUID id);
    
    /**
     * Найти все ответы
     */
    List<AnswerModel> findAll();
    
    /**
     * Найти ответы по ID вопроса
     */
    List<AnswerModel> findByQuestionId(UUID questionId);
    
    /**
     * Найти правильные ответы по ID вопроса
     */
    List<AnswerModel> findCorrectAnswersByQuestionId(UUID questionId);
    
    /**
     * Удалить ответ по ID
     */
    boolean deleteById(UUID id);
    
    /**
     * Удалить все ответы для вопроса
     */
    int deleteByQuestionId(UUID questionId);
    
    /**
     * Проверить существование ответа
     */
    boolean existsById(UUID id);
    
    /**
     * Проверить, существует ли ответ для такого вопроса
     */
    boolean existsByQuestionIdAndText(UUID questionId, String answerText);
    
    /**
     * Получить количество ответов для вопроса
     */
    int countByQuestionId(UUID questionId);
    
    /**
     * Получить количество правильных ответов для вопроса
     */
    int countCorrectAnswersByQuestionId(UUID questionId);
}
