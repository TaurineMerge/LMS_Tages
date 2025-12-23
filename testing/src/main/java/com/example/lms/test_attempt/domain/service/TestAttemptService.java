package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.shared.storage.StorageServiceInterface;
import com.example.lms.shared.storage.dto.UploadResult;
import com.example.lms.test_attempt.api.dto.TestAttempt;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Сервис для работы с попытками тестов.
 * <p>
 * Теперь интегрирован с MinIO для хранения снепшотов попыток.
 */
public class TestAttemptService {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);

    private final TestAttemptRepositoryInterface repository;
    private final StorageServiceInterface storageService;
    private final ObjectMapper objectMapper;

    public TestAttemptService(TestAttemptRepositoryInterface repository, StorageServiceInterface storageService) {
        this.repository = repository;
        this.storageService = storageService;
        this.objectMapper = new ObjectMapper();
    }

    // Старый конструктор для обратной совместимости
    public TestAttemptService(TestAttemptRepositoryInterface repository) {
        this(repository, null);
    }

    // DTO -> MODEL
    private TestAttemptModel toModel(TestAttempt dto) {
        LocalDate date = null;
        if (dto.getDate_of_attempt() != null && !dto.getDate_of_attempt().isEmpty()) {
            try {
                date = LocalDate.parse(dto.getDate_of_attempt());
            } catch (Exception e) {
                logger.error("Ошибка парсинга даты: {}", dto.getDate_of_attempt(), e);
            }
        }

        return new TestAttemptModel(
                dto.getId() != null ? UUID.fromString(dto.getId()) : null,
                dto.getStudent_id() != null ? UUID.fromString(dto.getStudent_id()) : null,
                dto.getTest_id() != null ? UUID.fromString(dto.getTest_id()) : null,
                date,
                dto.getPoint(),
                dto.getCertificate_id() != null ? UUID.fromString(dto.getCertificate_id()) : null,
                dto.getAttempt_version(),
                dto.getAttempt_snapshot(),
                dto.getCompleted());
    }

    // MODEL → DTO
    private TestAttempt toDto(TestAttemptModel model) {
        String dateStr = null;
        if (model.getDateOfAttempt() != null) {
            dateStr = model.getDateOfAttempt().toString();
        }

        return new TestAttempt(
                model.getId() != null ? model.getId().toString() : null,
                model.getStudentId() != null ? model.getStudentId().toString() : null,
                model.getTestId() != null ? model.getTestId().toString() : null,
                dateStr,
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

    public TestAttempt getTestAttemptById(UUID uuid) {
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

        // Если используем MinIO, удаляем снепшот из хранилища
        if (storageService != null) {
            try {
                TestAttemptModel model = repository.findById(uuid).orElse(null);
                if (model != null) {
                    storageService.deleteSnapshot(
                            model.getStudentId().toString(),
                            model.getTestId().toString(),
                            model.getId().toString());
                }
            } catch (Exception e) {
                logger.warn("Не удалось удалить снепшот из MinIO для попытки: {}", id, e);
            }
        }

        return repository.deleteById(uuid);
    }

    public TestAttempt completeTestAttempt(String id, Integer finalPoint) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();
        model.completeAttempt(finalPoint);

        TestAttemptModel updated = repository.update(model);

        // Сохраняем финальный снепшот в MinIO
        if (storageService != null && updated.getAttemptSnapshot() != null) {
            saveSnapshotToMinio(updated);
        }

        return toDto(updated);
    }

    /**
     * Обновляет снепшот попытки и сохраняет его в MinIO.
     * 
     * @param id       ID попытки
     * @param snapshot JSON-снепшот
     * @return обновленная попытка
     */
    public TestAttempt updateSnapshot(String id, String snapshot) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();

        // Сохраняем снепшот в MinIO
        if (storageService != null) {
            try {
                UploadResult result = storageService.uploadSnapshot(
                        model.getStudentId().toString(),
                        model.getTestId().toString(),
                        model.getId().toString(),
                        snapshot);

                logger.info("Снепшот сохранен в MinIO: {}", result.getObjectPath());

                // В БД сохраняем только ссылку на файл в MinIO
                model.updateSnapshot(result.getObjectPath());

            } catch (Exception e) {
                logger.error("Ошибка при сохранении снепшота в MinIO", e);
                // Fallback: сохраняем в БД напрямую
                model.updateSnapshot(snapshot);
            }
        } else {
            // Если MinIO не настроен, сохраняем в БД
            model.updateSnapshot(snapshot);
        }

        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    /**
     * Получает полный снепшот попытки из MinIO.
     * 
     * @param id ID попытки
     * @return JSON-снепшот, или null если не найден
     */
    public String getFullSnapshot(String id) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();

        // Если снепшот хранится в MinIO (путь начинается с "snapshots/")
        if (storageService != null && model.getAttemptSnapshot() != null
                && model.getAttemptSnapshot().startsWith("snapshots/")) {

            Optional<String> snapshot = storageService.downloadSnapshot(
                    model.getStudentId().toString(),
                    model.getTestId().toString(),
                    model.getId().toString());

            return snapshot.orElse(null);
        }

        // Иначе возвращаем из БД
        return model.getAttemptSnapshot();
    }

    /**
     * Получает список всех снепшотов студента для конкретного теста из MinIO.
     * 
     * @param studentId ID студента
     * @param testId    ID теста
     * @return список ID попыток
     */
    public List<String> getSnapshotsList(String studentId, String testId) {
        if (storageService == null) {
            logger.warn("StorageService не настроен, возвращаем пустой список");
            return List.of();
        }

        return storageService.listSnapshots(studentId, testId);
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

    // ======================================================================
    // PRIVATE HELPER METHODS
    // ======================================================================

    /**
     * Сохраняет снепшот попытки в MinIO.
     */
    private void saveSnapshotToMinio(TestAttemptModel model) {
        try {
            String snapshot = model.getAttemptSnapshot();
            if (snapshot == null || snapshot.isEmpty()) {
                logger.warn("Снепшот пустой, пропускаем сохранение в MinIO");
                return;
            }

            UploadResult result = storageService.uploadSnapshot(
                    model.getStudentId().toString(),
                    model.getTestId().toString(),
                    model.getId().toString(),
                    snapshot);

            logger.info("Снепшот сохранен в MinIO: {}", result.getObjectPath());

        } catch (Exception e) {
            logger.error("Ошибка при сохранении снепшота в MinIO", e);
            // Не падаем, продолжаем работу (снепшот остается в БД)
        }
    }
}