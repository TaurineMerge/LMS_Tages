package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.List;
import java.util.UUID;

import com.example.lms.test_attempt.api.dto.TestAttempt;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;

/**
 * Сервис для работы с попытками тестов.
 */
public class TestAttemptService {

    private final TestAttemptRepositoryInterface repository;

    public TestAttemptService(TestAttemptRepositoryInterface repository) {
        this.repository = repository;
    }

    // DTO -> MODEL
    private TestAttemptModel toModel(TestAttempt dto) {
        // Преобразуем строку даты в LocalDate
        LocalDate date = null;
        if (dto.getDate_of_attempt() != null && !dto.getDate_of_attempt().isEmpty()) {
            try {
                date = LocalDate.parse(dto.getDate_of_attempt());
            } catch (Exception e) {
                // Логируем ошибку, но не падаем
                System.err.println("Ошибка парсинга даты: " + dto.getDate_of_attempt());
            }
        }
        
        return new TestAttemptModel(
                dto.getId() != null ? UUID.fromString(dto.getId()) : null,
                dto.getStudent_id() != null ? UUID.fromString(dto.getStudent_id()) : null,
                dto.getTest_id() != null ? UUID.fromString(dto.getTest_id()) : null,
                date, // LocalDate или null
                dto.getPoint(),
                dto.getCertificate_id() != null ? UUID.fromString(dto.getCertificate_id()) : null,
                dto.getAttempt_version(),
                dto.getAttempt_snapshot(),
                dto.getCompleted());
    }

    // MODEL → DTO
    private TestAttempt toDto(TestAttemptModel model) {
        // Преобразуем LocalDate в строку
        String dateStr = null;
        if (model.getDateOfAttempt() != null) {
            dateStr = model.getDateOfAttempt().toString(); // ISO формат: "2025-12-18"
        }
        
        return new TestAttempt(
                model.getId() != null ? model.getId().toString() : null,
                model.getStudentId() != null ? model.getStudentId().toString() : null,
                model.getTestId() != null ? model.getTestId().toString() : null,
                dateStr, // String или null
                model.getPoint(),
                model.getCertificateId() != null ? model.getCertificateId().toString() : null,
                model.getAttemptVersion(),
                model.getAttemptSnapshot(),
                model.getCompleted());
    }

    // PUBLIC API METHODS
    public List<TestAttempt> getAllTestAttempts() {
        return repository.findAll().stream()
                .map(this::toDto)
                .toList();
    }

    public TestAttempt createTestAttempt(TestAttempt dto) {
        TestAttemptModel model = toModel(dto);
        TestAttemptModel saved = repository.save(model);
        return toDto(saved);
    }

    public TestAttempt getTestAttemptById(String id) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();
        return toDto(model);
    }

    public TestAttempt updateTestAttempt(TestAttempt dto) {
        TestAttemptModel model = toModel(dto);
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    public boolean deleteTestAttempt(String id) {
        UUID uuid = UUID.fromString(id);
        return repository.deleteById(uuid);
    }

    public TestAttempt completeTestAttempt(String id, Integer finalPoint) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();
        model.completeAttempt(finalPoint);
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    public TestAttempt updateSnapshot(String id, String snapshot) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();
        model.updateSnapshot(snapshot);
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    public List<TestAttempt> getTestAttemptsByStudentId(String studentId) {
        UUID uuid = UUID.fromString(studentId);
        return repository.findByStudentId(uuid).stream()
                .map(this::toDto)
                .toList();
    }

    public List<TestAttempt> getTestAttemptsByTestId(String testId) {
        UUID uuid = UUID.fromString(testId);
        return repository.findByTestId(uuid).stream()
                .map(this::toDto)
                .toList();
    }

    public List<TestAttempt> getTestAttemptsByStudentIdAndTestId(String studentId, String testId) {
        UUID studentUuid = UUID.fromString(studentId);
        UUID testUuid = UUID.fromString(testId);
        return repository.findByStudentIdAndTestId(studentUuid, testUuid).stream()
                .map(this::toDto)
                .toList();
    }

    public int countByStudentId(String studentId) {
        UUID uuid = UUID.fromString(studentId);
        return repository.countByStudentId(uuid);
    }

    public int countByTestId(String testId) {
        UUID uuid = UUID.fromString(testId);
        return repository.countByTestId(uuid);
    }

    public List<TestAttempt> getCompletedTestAttempts() {
        return repository.findCompletedAttempts().stream()
                .map(this::toDto)
                .toList();
    }
}