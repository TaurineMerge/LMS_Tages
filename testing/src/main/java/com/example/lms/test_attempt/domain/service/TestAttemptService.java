package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test_attempt.api.dto.TestAttempt;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.example.lms.shared.storage.StorageServiceInterface;
import com.example.lms.shared.storage.dto.UploadResult;

/**
 * Сервис для работы с попытками тестов.
 *
 * ✅ Сохранён весь функционал HEAD (DTO ↔ Model + CRUD/поиск),
 * ✅ и добавлена UI-часть (attempt_version JSON, saveAnswer,
 * completeAttemptForToday).
 */
public class TestAttemptService {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private final TestAttemptRepositoryInterface repository;
    private final StorageServiceInterface storageService;

    public TestAttemptService(TestAttemptRepositoryInterface repository, StorageServiceInterface storageService) {
        this.repository = repository;
        this.storageService = storageService;
    }

    // ---------------------------------------------------------------------
    // DTO -> MODEL (как в HEAD)
    // ---------------------------------------------------------------------
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

    // ---------------------------------------------------------------------
    // MODEL -> DTO (как в HEAD)
    // ---------------------------------------------------------------------
    private TestAttempt toDto(TestAttemptModel model) {
        String dateStr = (model.getDateOfAttempt() == null) ? null : model.getDateOfAttempt().toString();

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

    // ---------------------------------------------------------------------
    // PUBLIC API METHODS (как в HEAD)
    // ---------------------------------------------------------------------
    public List<TestAttempt> getAllTestAttempts() {
        return repository.findAll().stream().map(this::toDto).toList();
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
                        snapshot, model.getAttemptVersion(), model.getDateOfAttempt());

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
        return repository.findByStudentId(uuid).stream().map(this::toDto).toList();
    }

    public List<TestAttempt> getTestAttemptsByTestId(String testId) {
        UUID uuid = UUID.fromString(testId);
        return repository.findByTestId(uuid).stream().map(this::toDto).toList();
    }

    public List<TestAttempt> getTestAttemptsByStudentIdAndTestId(String studentId, String testId) {
        UUID studentUuid = UUID.fromString(studentId);
        UUID testUuid = UUID.fromString(testId);
        return repository.findByStudentIdAndTestId(studentUuid, testUuid).stream().map(this::toDto).toList();
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
        return repository.findCompletedAttempts().stream().map(this::toDto).toList();
    }

    // ---------------------------------------------------------------------
    // UI METHODS (чтобы UI часть работала)
    // ---------------------------------------------------------------------

    /**
     * ✅ UI: получить attempt_version (JSON).
     */
    public Optional<String> getAttemptVersion(UUID studentId, UUID testId, LocalDate date) {
        if (date == null) {
            date = LocalDate.now();
        }
        return repository.findAttemptVersion(studentId, testId, date.toString());
    }

    /**
     * UI: Сохраняет ответ студента на конкретный вопрос в attempt_version (JSON).
     *
     * Формат attempt_version:
     * {
     * "answers": [
     * { "question": "<questionUuid>", "answer": "<answerUuid>" }
     * ]
     * }
     *
     * Правило "без возможности исправления":
     * если в answers уже есть запись с question == questionId, повтор игнорируем.
     */
    public void saveAnswer(UUID studentId, UUID testId, UUID questionId, UUID answerId) {
        String today = LocalDate.now().toString();

        try {
            String existingJson = repository.findAttemptVersion(studentId, testId, today).orElse(null);

            ObjectNode root = toObjectNodeOrEmpty(existingJson);
            ArrayNode answersArray = ensureAnswersArray(root);

            String q = questionId.toString();
            String a = answerId.toString();

            for (JsonNode item : answersArray) {
                if (item != null && item.isObject()) {
                    JsonNode qNode = item.get("question");
                    if (qNode != null && q.equals(qNode.asText())) {
                        logger.info(
                                "Ответ уже сохранён (student={}, test={}, date={}, question={}). Повтор игнорируем.",
                                studentId, testId, today, questionId);
                        return;
                    }
                }
            }

            ObjectNode entry = OBJECT_MAPPER.createObjectNode();
            entry.put("question", q);
            entry.put("answer", a);
            answersArray.add(entry);

            String updatedJson = OBJECT_MAPPER.writeValueAsString(root);
            repository.upsertAttemptVersion(studentId, testId, today, updatedJson);

            logger.info("Ответ сохранён (student={}, test={}, date={}, question={}, answer={})",
                    studentId, testId, today, questionId, answerId);

        } catch (Exception e) {
            logger.error("Ошибка при сохранении ответа в attempt_version", e);
            throw new RuntimeException("Ошибка при сохранении ответа", e);
        }
    }

    /**
     * UI: завершить попытку за сегодня (если записи за сегодня нет — создаём).
     */
    public void completeAttemptForToday(UUID studentId, UUID testId, int points) {
        LocalDate today = LocalDate.now();

        List<TestAttemptModel> attempts = repository.findByStudentIdAndTestId(studentId, testId);

        TestAttemptModel todayAttempt = null;
        for (TestAttemptModel a : attempts) {
            if (a.getDateOfAttempt() != null && today.equals(a.getDateOfAttempt())) {
                todayAttempt = a;
                break;
            }
        }

        if (todayAttempt == null) {
            todayAttempt = new TestAttemptModel(studentId, testId);
            todayAttempt.validate();
            todayAttempt = repository.save(todayAttempt);
        }

        todayAttempt.completeAttempt(points);
        todayAttempt.validate();
        repository.update(todayAttempt);
    }

    // ---------------------------------------------------------------------
    // JSON helpers
    // ---------------------------------------------------------------------

    private ArrayNode ensureAnswersArray(ObjectNode root) {
        JsonNode answers = root.get("answers");

        if (answers != null && answers.isArray()) {
            return (ArrayNode) answers;
        }

        if (answers != null && answers.isObject()) {
            ArrayNode arr = OBJECT_MAPPER.createArrayNode();
            answers.fields().forEachRemaining(entry -> {
                ObjectNode item = OBJECT_MAPPER.createObjectNode();
                item.put("question", entry.getKey());
                JsonNode v = entry.getValue();
                item.put("answer", v == null ? null : v.asText());
                arr.add(item);
            });
            root.set("answers", arr);
            return arr;
        }

        ArrayNode arr = OBJECT_MAPPER.createArrayNode();
        root.set("answers", arr);
        return arr;
    }

    private ObjectNode toObjectNodeOrEmpty(String json) {
        if (json == null || json.isBlank()) {
            return OBJECT_MAPPER.createObjectNode();
        }
        try {
            JsonNode node = OBJECT_MAPPER.readTree(json);
            if (node != null && node.isObject()) {
                return (ObjectNode) node;
            }
        } catch (Exception ignored) {
        }
        return OBJECT_MAPPER.createObjectNode();
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
                    model.getAttemptVersion(),
                    snapshot, model.getDateOfAttempt());

            logger.info("Снепшот сохранен в MinIO: {}", result.getObjectPath());

        } catch (Exception e) {
            logger.error("Ошибка при сохранении снепшота в MinIO", e);
            // Не падаем, продолжаем работу (снепшот остается в БД)
        }
    }
}