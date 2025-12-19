package com.example.lms.test_attempt.domain.repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import com.example.lms.test_attempt.domain.model.TestAttemptModel;

/**
 * Репозиторий для работы с попытками прохождения тестов.
 *
 * Таблица: testing.test_attempt_b
 * Поля:
 *  - id UUID (PK)
 *  - student_id UUID (not null)
 *  - test_id UUID (not null)
 *  - date_of_attempt DATE
 *  - point INTEGER
 *  - certificate_id UUID
 *  - attempt_version JSON
 *  - attempt_snapshot VARCHAR(256)
 *  - completed BOOLEAN
 */
public interface TestAttemptRepositoryInterface {

    // ---------------- CRUD ----------------

    TestAttemptModel save(TestAttemptModel testAttempt);

    TestAttemptModel update(TestAttemptModel testAttempt);

    Optional<TestAttemptModel> findById(UUID id);

    List<TestAttemptModel> findAll();

    boolean deleteById(UUID id);

    boolean existsById(UUID id);

    // ---------------- Queries ----------------

    List<TestAttemptModel> findByStudentId(UUID studentId);

    List<TestAttemptModel> findByTestId(UUID testId);

    List<TestAttemptModel> findByStudentIdAndTestId(UUID studentId, UUID testId);

    int countByStudentId(UUID studentId);

    int countByTestId(UUID testId);

    List<TestAttemptModel> findCompletedAttempts();

    List<TestAttemptModel> findIncompleteAttempts();

    // ---------------- UI: attempt_version ----------------

    /**
     * Вернуть attempt_version (JSON) для попытки (student_id + test_id + date_of_attempt).
     *
     * @param date ISO-строка формата YYYY-MM-DD (пример: "2025-12-19")
     */
    Optional<String> findAttemptVersion(UUID studentId, UUID testId, String date);

    /**
     * Insert или Update attempt_version (JSON) для попытки (student_id + test_id + date_of_attempt).
     *
     * @param date ISO-строка формата YYYY-MM-DD (пример: "2025-12-19")
     */
    void upsertAttemptVersion(UUID studentId, UUID testId, String date, String attemptVersionJson);
}