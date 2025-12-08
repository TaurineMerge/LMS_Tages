package com.example.lms.test.domain.repository;

import com.example.lms.test.domain.model.TestModel;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Интерфейс репозитория для работы с тестами
 * ОПРЕДЕЛЯЕТ ТОЛЬКО МЕТОДЫ - без моделей!
 */
public interface TestRepositoryInterface {
    
    /**
     * Сохранить новый тест
     */
    TestModel save(TestModel test);
    
    /**
     * Обновить существующий тест
     */
    TestModel update(TestModel test);
    
    /**
     * Найти тест по ID
     */
    Optional<TestModel> findById(UUID id);
    
    /**
     * Найти все тесты
     */
    List<TestModel> findAll();
    
    /**
     * Найти тесты по ID курса
     */
    List<TestModel> findByCourseId(UUID courseId);
    
    /**
     * Удалить тест по ID
     */
    boolean deleteById(UUID id);
    
    /**
     * Проверить существование теста
     */
    boolean existsById(UUID id);
    
    /**
     * Поиск тестов по названию
     */
    List<TestModel> findByTitleContaining(String title);
    
    /**
     * Получить количество тестов в курсе
     */
    int countByCourseId(UUID courseId);
}