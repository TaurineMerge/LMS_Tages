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
 * Сервис для попыток тестов.
 *
 * UI-логика:
 * - attempt_version хранится по attemptId (PK)
 * - attempt_version инициализируется сразу всеми вопросами
 * - attemptNo пишем в JSON
 *
 * Формат attempt_version (АКТУАЛЬНЫЙ):
 * {
 *   "attemptNo": 1,
 *   "answers": [
 *     {"order":1,"questionId":"...","answerIds":[],"answerPoints":[],"earnedPoints":0},
 *     ...
 *   ]
 * }
 *
 * Важно:
 * - answerId (legacy) НЕ пишем вообще
 * - но при чтении старых данных можем встретить answerId и не падаем
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

    // =========================================================
    // DTO <-> MODEL
    // =========================================================

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

    // =========================================================
    // CRUD / Queries
    // =========================================================

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

    // =========================================================
    // UI helpers types
    // =========================================================

    public static class QuestionInit {
        public final UUID questionId;
        public final Integer order;

        public QuestionInit(UUID questionId, Integer order) {
            this.questionId = questionId;
            this.order = order;
        }
    }

    // =========================================================
    // UI API (attemptId)
    // =========================================================

    /**
     * Создать новую попытку (attemptId = PK).
     *
     * Важно: НЕ ставим date_of_attempt автоматически, иначе при UNIQUE(student_id,test_id,date_of_attempt)
     * повторная попытка в тот же день будет падать.
     */
    public UUID createNewAttempt(UUID studentId, UUID testId) {
        TestAttemptModel m = new TestAttemptModel(studentId, testId);
        m.setDateOfAttempt(LocalDate.now()); 
        m.setCompleted(false);
        m.validate();
        TestAttemptModel saved = repository.save(m);
        return saved.getId();

    }

    /**
     * Считает только завершённые попытки (completed=true ИЛИ point != null).
     */
    public int countCompletedAttempts(UUID studentId, UUID testId) {
        return (int) repository.findByStudentIdAndTestId(studentId, testId).stream()
                .filter(a -> Boolean.TRUE.equals(a.getCompleted()) || a.getPoint() != null)
                .count();
    }

    /**
     * Последняя незавершённая попытка студента по тесту (если есть).
     * Незавершённая = completed=false AND point IS NULL.
     */
    public Optional<TestAttempt> getLatestIncompleteAttempt(UUID studentId, UUID testId) {
        return repository.findLatestIncompleteByStudentAndTestId(studentId, testId).map(this::toDto);
    }

    /**
     * Последняя завершённая попытка студента по тесту (если есть).
     * Завершённая = completed=true OR point != null.
     */
    public Optional<TestAttempt> getLatestCompletedAttempt(UUID studentId, UUID testId) {
        return repository.findLatestCompletedByStudentAndTestId(studentId, testId).map(this::toDto);
    }

    public Optional<String> getAttemptVersionByAttemptId(UUID attemptId) {
        return repository.findAttemptVersionByAttemptId(attemptId);
    }

    /**
     * Инициализируем JSON, если attempt_version пустой.
     *
     * Формат:
     * {
     *   "attemptNo": 1,
     *   "answers": [
     *     {"order":1,"questionId":"...","answerIds":[],"answerPoints":[],"earnedPoints":0},
     *     ...
     *   ]
     * }
     */
    public void initAttemptVersionIfEmpty(UUID attemptId, int attemptNo, List<QuestionInit> questions) {
        String existing = repository.findAttemptVersionByAttemptId(attemptId).orElse(null);
        if (existing != null && !existing.isBlank()) return;

        try {
            ObjectNode root = OBJECT_MAPPER.createObjectNode();
            root.put("attemptNo", attemptNo);

            ArrayNode answers = OBJECT_MAPPER.createArrayNode();
            for (QuestionInit q : questions) {
                ObjectNode item = OBJECT_MAPPER.createObjectNode();
                if (q.order != null) item.put("order", q.order);
                item.put("questionId", q.questionId.toString());

                // ✅ только новый формат
                item.set("answerIds", OBJECT_MAPPER.createArrayNode());
                item.set("answerPoints", OBJECT_MAPPER.createArrayNode());
                item.put("earnedPoints", 0);

                answers.add(item);
            }
            root.set("answers", answers);

            String json = OBJECT_MAPPER.writeValueAsString(root);
            repository.updateAttemptVersionByAttemptId(attemptId, json);
        } catch (Exception e) {
            logger.error("Ошибка при init attempt_version", e);
            throw new RuntimeException("Ошибка init attempt_version", e);
        }
    }

    /**
     * Оставляем старую сигнатуру для совместимости компиляции (если где-то ещё вызывается).
     * Баллы по ответам не запишет (0), но хотя бы сохранит answerIds.
     */
    public void saveAnswers(UUID attemptId, UUID questionId, List<UUID> answerIds) {
        saveAnswers(attemptId, questionId, answerIds, List.of(), 0);
    }

    /**
     * Multi API: сохраняем выбранные ответы + баллы.
     * Если уже есть ответ(ы) — повтор игнорируем.
     *
     * Пишем:
     * - answerIds: [...]
     * - answerPoints: [...]
     * - earnedPoints: N
     */
    public void saveAnswers(UUID attemptId,
                            UUID questionId,
                            List<UUID> answerIds,
                            List<Integer> answerPoints,
                            int earnedPoints) {
        try {
            String json = repository.findAttemptVersionByAttemptId(attemptId).orElse(null);
            ObjectNode root = toObjectNodeOrEmpty(json);
            ArrayNode answers = ensureAnswersArray(root);

            String q = questionId.toString();

            ArrayNode idsNode = OBJECT_MAPPER.createArrayNode();
            for (UUID id : answerIds) idsNode.add(id.toString());

            ArrayNode ptsNode = OBJECT_MAPPER.createArrayNode();
            if (answerPoints != null && answerPoints.size() == answerIds.size()) {
                for (Integer p : answerPoints) ptsNode.add(Math.max(0, p == null ? 0 : p));
            } else {
                // если points не передали или размер не совпал — пишем нули по количеству ответов
                for (int i = 0; i < answerIds.size(); i++) ptsNode.add(0);
                earnedPoints = 0;
            }

            for (JsonNode item : answers) {
                if (item != null && item.isObject()) {
                    String qid = item.path("questionId").asText(null);
                    if (q.equals(qid)) {
                        ObjectNode obj = (ObjectNode) item;

                        // если уже отвечено — игнор
                        JsonNode existingIds = obj.get("answerIds");
                        if (existingIds != null && existingIds.isArray() && existingIds.size() > 0) {
                            logger.info("Ответ уже был выбран (attemptId={}, questionId={}) — игнор", attemptId, questionId);
                            return;
                        }

                        obj.set("answerIds", idsNode);
                        obj.set("answerPoints", ptsNode);
                        obj.put("earnedPoints", Math.max(0, earnedPoints));

                        // на всякий: если в старых данных был answerId — удалим, чтобы формат стал единым
                        obj.remove("answerId");

                        String updated = OBJECT_MAPPER.writeValueAsString(root);
                        repository.updateAttemptVersionByAttemptId(attemptId, updated);
                        return;
                    }
                }
            }

            // если вопроса нет — добавим
            ObjectNode newItem = OBJECT_MAPPER.createObjectNode();
            newItem.put("questionId", q);
            newItem.set("answerIds", idsNode);
            newItem.set("answerPoints", ptsNode);
            newItem.put("earnedPoints", Math.max(0, earnedPoints));
            answers.add(newItem);

            String updated = OBJECT_MAPPER.writeValueAsString(root);
            repository.updateAttemptVersionByAttemptId(attemptId, updated);

        } catch (Exception e) {
            logger.error("Ошибка saveAnswers(with points)", e);
            throw new RuntimeException("Ошибка сохранения ответа", e);
        }
    }

    /**
     * Завершить попытку по attemptId.
     *
     * Важно: date_of_attempt здесь НЕ проставляем автоматически.
     */
    public void completeAttemptById(UUID attemptId, int points) {
        TestAttemptModel model = repository.findById(attemptId).orElseThrow();
        model.completeAttempt(points);
        model.validate();
        repository.update(model);
    }

    // =========================================================
    // JSON helpers
    // =========================================================

    private ArrayNode ensureAnswersArray(ObjectNode root) {
        JsonNode answers = root.get("answers");
        if (answers != null && answers.isArray()) {
            return (ArrayNode) answers;
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
            if (node != null && node.isObject()) return (ObjectNode) node;
        } catch (Exception ignored) {}
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
    // ---------------------------------------------------------------------
    // BACKWARD COMPATIBILITY (for TestAttemptController)
    // ---------------------------------------------------------------------
}