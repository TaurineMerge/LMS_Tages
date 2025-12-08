package com.example.lms.test_attempt.api.domain.repository;

import com.example.lms.test_attempt.api.domain.model.Test_AttemptModel;
import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Интерфейс репозитория для работы с попытками прохождения тестов
 * Соответствует таблице test_attempt_b
 */
public interface Test_AttemptRepositoryInterface {
    
    /**
     * Сохранить новую попытку теста
     * @param attempt модель попытки
     * @return сохраненная модель с присвоенным ID
     */
    Test_AttemptModel save(Test_AttemptModel attempt);
    
    /**
     * Обновить существующую попытку
     * @param attempt модель попытки с обновленными данными
     * @return обновленная модель
     */
    Test_AttemptModel update(Test_AttemptModel attempt);
    
    /**
     * Найти попытку по ID
     * @param id идентификатор попытки
     * @return Optional с попыткой, если найдена
     */
    Optional<Test_AttemptModel> findById(UUID id);
    
    /**
     * Найти все попытки
     * @return список всех попыток
     */
    List<Test_AttemptModel> findAll();
    
    /**
     * Найти попытки по ID студента
     * @param studentId идентификатор студента
     * @return список попыток студента
     */
    List<Test_AttemptModel> findByStudentId(UUID studentId);
    
    /**
     * Найти попытки по ID теста
     * @param testId идентификатор теста
     * @return список попыток для теста
     */
    List<Test_AttemptModel> findByTestId(UUID testId);
    
    /**
     * Найти попытку студента для конкретного теста
     * @param studentId идентификатор студента
     * @param testId идентификатор теста
     * @return список попыток (обычно 0 или 1, но может быть несколько)
     */
    List<Test_AttemptModel> findByStudentAndTest(UUID studentId, UUID testId);
    
    /**
     * Удалить попытку по ID
     * @param id идентификатор попытки
     * @return true если попытка была удалена
     */
    boolean deleteById(UUID id);
    
    /**
     * Проверить существование попытки
     * @param id идентификатор попытки
     * @return true если попытка существует
     */
    boolean existsById(UUID id);
    
    /**
     * Получить количество попыток студента для теста
     * @param studentId идентификатор студента
     * @param testId идентификатор теста
     * @return количество попыток
     */
    int countAttemptsByStudentAndTest(UUID studentId, UUID testId);
    
    /**
     * Найти попытки по дате
     * @param date дата попытки
     * @return список попыток за указанную дату
     */
    List<Test_AttemptModel> findByDate(LocalDate date);
    
    /**
     * Найти завершенные попытки (с результатами)
     * @return список завершенных попыток
     */
    List<Test_AttemptModel> findCompletedAttempts();
    
    /**
     * Найти незавершенные попытки (без результатов)
     * @return список незавершенных попыток
     */
    List<Test_AttemptModel> findIncompleteAttempts();
    
    /**
     * Получить лучший результат студента по тесту
     * @param studentId идентификатор студента
     * @param testId идентификатор теста
     * @return Optional с лучшей попыткой (по баллам)
     */
    Optional<Test_AttemptModel> findBestAttemptByStudentAndTest(UUID studentId, UUID testId);
}