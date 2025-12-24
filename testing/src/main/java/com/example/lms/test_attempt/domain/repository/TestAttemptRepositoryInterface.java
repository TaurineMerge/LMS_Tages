package com.example.lms.test_attempt.domain.repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import com.example.lms.test_attempt.domain.model.TestAttemptModel;

/**
 * Репозиторий для работы с попытками прохождения тестов.
 *
 * Таблица: testing.test_attempt_b
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

    // ---------------- UI: attempt_version by attemptId (PK) ----------------
    Optional<String> findAttemptVersionByAttemptId(UUID attemptId);
    void updateAttemptVersionByAttemptId(UUID attemptId, String attemptVersionJson);
}