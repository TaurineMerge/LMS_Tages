package com.example.lms.test_attempt.domain.repository;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import com.example.lms.test_attempt.domain.model.TestAttemptModel;

/**
 * Репозиторий для работы с попытками прохождения тестов.
 * Соответствует таблице TEST_ATTEMPT_B.
 */
public interface TestAttemptRepositoryInterface {

	/**
	 * Сохранить новую попытку.
	 */
	TestAttemptModel save(TestAttemptModel attempt);

	/**
	 * Обновить существующую попытку.
	 */
	TestAttemptModel update(TestAttemptModel attempt);

	/**
	 * Найти попытку по её ID.
	 */
	Optional<TestAttemptModel> findById(UUID id);

	/**
	 * Получить все попытки.
	 */
	List<TestAttemptModel> findAll();

	/**
	 * Получить все попытки студента.
	 */
	List<TestAttemptModel> findByStudentId(UUID studentId);

	/**
	 * Получить все попытки по тесту.
	 */
	List<TestAttemptModel> findByTestId(UUID testId);

	/**
	 * Получить все попытки конкретного студента по конкретному тесту.
	 */
	List<TestAttemptModel> findByStudentAndTest(UUID studentId, UUID testId);

	/**
	 * Удалить попытку по ID.
	 */
	boolean deleteById(UUID id);

	/**
	 * Проверить существование попытки по ID.
	 */
	boolean existsById(UUID id);

	/**
	 * Подсчитать количество попыток студента по тесту.
	 */
	int countAttemptsByStudentAndTest(UUID studentId, UUID testId);

	/**
	 * Найти попытки по дате прохождения.
	 */
	List<TestAttemptModel> findByDate(LocalDate date);

	/**
	 * Найти все завершённые попытки (где point IS NOT NULL).
	 */
	List<TestAttemptModel> findCompletedAttempts();

	/**
	 * Найти все незавершённые попытки (где point IS NULL).
	 */
	List<TestAttemptModel> findIncompleteAttempts();
}
