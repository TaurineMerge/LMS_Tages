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
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;

/**
 * Сервис для работы с попытками тестов.
 *
 * ВАЖНО для UI:
 * - attempt_version хранит тексты вопросов/ответов, чтобы результаты жили даже после удаления/изменения теста
 * - добавлен upsertAnswers(...) для кнопки "Сохранить" (перезапись ответов -> можно исправлять)
 * - старый saveAnswers(...) оставлен (legacy-логика: если ответ уже есть — игнор)
 *
 * Формат attempt_version (расширенный, обратно-совместимый):
 * {
 *   "attemptNo": 1,
 *   "testTitle": "....",        // optional
 *   "minPoint": 10,             // optional
 *   "answers": [
 *     {
 *       "order": 1,
 *       "questionId": "uuid",
 *       "questionText": "....", // optional
 *       "maxPoints": 3,         // optional
 *       "answerIds": ["..."],
 *       "answerTexts": ["..."], // optional (NEW)
 *       "answerPoints": [1,0],
 *       "earnedPoints": 1
 *     }
 *   ]
 * }
 */
public class TestAttemptService {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private final TestAttemptRepositoryInterface repository;
    private final StorageServiceInterface storageService; // может быть null

    public TestAttemptService(TestAttemptRepositoryInterface repository, StorageServiceInterface storageService) {
        this.repository = repository;
        this.storageService = storageService;
    }

    // =========================================================
    // DTO <-> MODEL
    // =========================================================

    private TestAttemptModel toModel(TestAttempt dto) {
        LocalDate date = null;
        if (dto.getDate_of_attempt() != null && !dto.getDate_of_attempt().isBlank()) {
            try {
                date = LocalDate.parse(dto.getDate_of_attempt().trim());
            } catch (Exception e) {
                logger.warn("Не удалось распарсить date_of_attempt: {}", dto.getDate_of_attempt(), e);
            }
        }

        UUID certId = null;
        if (dto.getCertificate_id() != null && !dto.getCertificate_id().isBlank()) {
            try {
                certId = UUID.fromString(dto.getCertificate_id().trim());
            } catch (Exception e) {
                logger.warn("Некорректный certificate_id: {}", dto.getCertificate_id(), e);
            }
        }

        return new TestAttemptModel(
                dto.getId() != null ? UUID.fromString(dto.getId()) : null,
                dto.getStudent_id() != null ? UUID.fromString(dto.getStudent_id()) : null,
                dto.getTest_id() != null ? UUID.fromString(dto.getTest_id()) : null,
                date,
                dto.getPoint(),
                certId,
                dto.getAttempt_version(),
                dto.getAttempt_snapshot(),
                dto.getCompleted()
        );
    }

    private TestAttempt toDto(TestAttemptModel model) {
        return new TestAttempt(
                model.getId() != null ? model.getId().toString() : null,
                model.getStudentId() != null ? model.getStudentId().toString() : null,
                model.getTestId() != null ? model.getTestId().toString() : null,
                model.getDateOfAttempt() != null ? model.getDateOfAttempt().toString() : null,
                model.getPoint(),
                model.getCertificateId() != null ? model.getCertificateId().toString() : null,
                model.getAttemptVersion(),
                model.getAttemptSnapshot(),
                model.getCompleted()
        );
    }

    // =========================================================
    // CRUD (API Controller использует это)
    // =========================================================

    public List<TestAttempt> getAllTestAttempts() {
        return repository.findAll().stream().map(this::toDto).toList();
    }

    public TestAttempt createTestAttempt(TestAttempt dto) {
        TestAttemptModel model = toModel(dto);
        model.validate();
        TestAttemptModel saved = repository.save(model);
        return toDto(saved);
    }

    public TestAttempt getTestAttemptById(UUID uuid) {
        TestAttemptModel model = repository.findById(uuid).orElseThrow();
        return toDto(model);
    }

    public TestAttempt updateTestAttempt(TestAttempt dto) {
        TestAttemptModel model = toModel(dto);
        model.validate();
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    public boolean deleteTestAttempt(String id) {
        UUID uuid = UUID.fromString(id);

        // если MinIO есть — удаляем снепшот
        if (storageService != null) {
            try {
                TestAttemptModel model = repository.findById(uuid).orElse(null);
                if (model != null) {
                    storageService.deleteSnapshot(
                            model.getStudentId().toString(),
                            model.getTestId().toString(),
                            model.getId().toString()
                    );
                }
            } catch (Exception e) {
                logger.warn("Не удалось удалить snapshot в MinIO для attemptId={}", id, e);
            }
        }

        return repository.deleteById(uuid);
    }

    public TestAttempt completeTestAttempt(String id, Integer finalPoint) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();

        model.completeAttempt(finalPoint != null ? finalPoint : 0);
        model.validate();

        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    public TestAttempt updateSnapshot(String id, String snapshot) {
        UUID uuid = UUID.fromString(id);
        TestAttemptModel model = repository.findById(uuid).orElseThrow();

        model.setAttemptSnapshot(snapshot);

        // Если есть MinIO — грузим туда, но УЖЕ по твоей сигнатуре (6 аргументов)
        if (storageService != null) {
            try {
                String attemptVersionJson = model.getAttemptVersion();
                LocalDate attemptDate = model.getDateOfAttempt();
                if (attemptDate == null) attemptDate = LocalDate.now();

                UploadResult result = storageService.uploadSnapshot(
                        model.getStudentId().toString(),
                        model.getTestId().toString(),
                        model.getId().toString(),
                        snapshot,
                        attemptVersionJson,
                        attemptDate
                );
                logger.info("Snapshot uploaded to MinIO: {}", result.getObjectPath());
            } catch (Exception e) {
                logger.warn("Не удалось загрузить snapshot в MinIO (останется в БД): attemptId={}", id, e);
            }
        }

        model.validate();
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    public List<TestAttempt> getTestAttemptsByStudentId(String studentId) {
        UUID uuid = UUID.fromString(studentId);
        return repository.findByStudentId(uuid).stream().map(this::toDto).toList();
    }

    public List<TestAttempt> getTestAttemptsByTestId(String testId) {
        UUID uuid = UUID.fromString(testId);
        return repository.findByTestId(uuid).stream().map(this::toDto).toList();
    }

    // =========================================================
    // UI helpers (для UiTestController)
    // =========================================================

    public int countCompletedAttempts(UUID studentId, UUID testId) {
        return (int) repository.findByStudentIdAndTestId(studentId, testId).stream()
                .filter(a -> Boolean.TRUE.equals(a.getCompleted()) || a.getPoint() != null)
                .count();
    }

    public Optional<TestAttempt> getLatestIncompleteAttempt(UUID studentId, UUID testId) {
        return repository.findLatestIncompleteByStudentAndTestId(studentId, testId).map(this::toDto);
    }

    public Optional<TestAttempt> getLatestCompletedAttempt(UUID studentId, UUID testId) {
        return repository.findLatestCompletedByStudentAndTestId(studentId, testId).map(this::toDto);
    }

    public Optional<String> getAttemptVersionByAttemptId(UUID attemptId) {
        return repository.findAttemptVersionByAttemptId(attemptId);
    }

    public static class QuestionInit {
        public final UUID questionId;
        public final Integer order;

        // NEW: чтобы результаты жили без test/question/answer
        public final String questionText;
        public final Integer maxPoints;

        public QuestionInit(UUID questionId, Integer order) {
            this(questionId, order, null, null);
        }

        public QuestionInit(UUID questionId, Integer order, String questionText, Integer maxPoints) {
            this.questionId = questionId;
            this.order = order;
            this.questionText = questionText;
            this.maxPoints = maxPoints;
        }
    }

    /**
     * Создать новую попытку (attemptId = PK).
     * В твоём TestAttemptModel НЕТ create(...), поэтому используем конструктор.
     */
    public UUID createNewAttempt(UUID studentId, UUID testId) {
        TestAttemptModel model = new TestAttemptModel(studentId, testId);
        model.validate();
        TestAttemptModel saved = repository.save(model);
        return saved.getId();
    }

    /**
     * Инициализировать attempt_version (если пустой) — сразу всеми вопросами.
     * Пишем questionText/maxPoints + поля для ответов.
     */
    public void initAttemptVersionIfEmpty(UUID attemptId,
                                         int attemptNo,
                                         List<QuestionInit> questions,
                                         String testTitle,
                                         Integer minPoint) {
        String existing = repository.findAttemptVersionByAttemptId(attemptId).orElse(null);
        if (existing != null && !existing.isBlank()) return;

        try {
            ObjectNode root = OBJECT_MAPPER.createObjectNode();
            root.put("attemptNo", attemptNo);

            if (testTitle != null) root.put("testTitle", testTitle);
            if (minPoint != null) root.put("minPoint", minPoint);

            ArrayNode answers = OBJECT_MAPPER.createArrayNode();
            for (QuestionInit q : questions) {
                ObjectNode item = OBJECT_MAPPER.createObjectNode();
                if (q.order != null) item.put("order", q.order);
                item.put("questionId", q.questionId.toString());

                if (q.questionText != null) item.put("questionText", q.questionText);
                if (q.maxPoints != null) item.put("maxPoints", Math.max(0, q.maxPoints));

                item.set("answerIds", OBJECT_MAPPER.createArrayNode());
                item.set("answerTexts", OBJECT_MAPPER.createArrayNode());
                item.set("answerPoints", OBJECT_MAPPER.createArrayNode());
                item.put("earnedPoints", 0);

                answers.add(item);
            }
            root.set("answers", answers);

            repository.updateAttemptVersionByAttemptId(attemptId, OBJECT_MAPPER.writeValueAsString(root));
        } catch (Exception e) {
            logger.error("Ошибка initAttemptVersionIfEmpty", e);
            throw new RuntimeException("Ошибка initAttemptVersionIfEmpty", e);
        }
    }

    /** Совместимость */
    public void initAttemptVersionIfEmpty(UUID attemptId, int attemptNo, List<QuestionInit> questions) {
        initAttemptVersionIfEmpty(attemptId, attemptNo, questions, null, null);
    }

    /** legacy */
    public void saveAnswers(UUID attemptId, UUID questionId, List<UUID> answerIds) {
        saveAnswers(attemptId, questionId, answerIds, List.of(), 0);
    }

    /**
     * Legacy-сохранение: если уже есть ответ — игнор (чтобы не ломать старый флоу).
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
                for (int i = 0; i < answerIds.size(); i++) ptsNode.add(0);
                earnedPoints = 0;
            }

            for (JsonNode item : answers) {
                if (item != null && item.isObject()) {
                    String qid = item.path("questionId").asText(null);
                    if (q.equals(qid)) {
                        ObjectNode obj = (ObjectNode) item;

                        JsonNode existingIds = obj.get("answerIds");
                        if (existingIds != null && existingIds.isArray() && existingIds.size() > 0) {
                            return; // legacy: повтор игнор
                        }

                        obj.set("answerIds", idsNode);
                        obj.set("answerPoints", ptsNode);
                        obj.put("earnedPoints", Math.max(0, earnedPoints));
                        obj.remove("answerId");

                        repository.updateAttemptVersionByAttemptId(attemptId, OBJECT_MAPPER.writeValueAsString(root));
                        return;
                    }
                }
            }

            ObjectNode newItem = OBJECT_MAPPER.createObjectNode();
            newItem.put("questionId", q);
            newItem.set("answerIds", idsNode);
            newItem.set("answerPoints", ptsNode);
            newItem.put("earnedPoints", Math.max(0, earnedPoints));
            answers.add(newItem);

            repository.updateAttemptVersionByAttemptId(attemptId, OBJECT_MAPPER.writeValueAsString(root));
        } catch (Exception e) {
            logger.error("Ошибка saveAnswers", e);
            throw new RuntimeException("Ошибка saveAnswers", e);
        }
    }

    /**
     * NEW: upsert для кнопки "Сохранить" (перезаписывает, можно исправлять).
     * Плюс пишет тексты (questionText/answerTexts) и maxPoints для устойчивых результатов.
     */
    public void upsertAnswers(UUID attemptId,
                              UUID questionId,
                              String questionText,
                              int maxPoints,
                              List<UUID> answerIds,
                              List<String> answerTexts,
                              List<Integer> answerPoints,
                              int earnedPoints) {
        try {
            String json = repository.findAttemptVersionByAttemptId(attemptId).orElse(null);
            ObjectNode root = toObjectNodeOrEmpty(json);
            ArrayNode answers = ensureAnswersArray(root);

            String q = questionId.toString();

            ArrayNode idsNode = OBJECT_MAPPER.createArrayNode();
            if (answerIds != null) {
                for (UUID id : answerIds) if (id != null) idsNode.add(id.toString());
            }

            ArrayNode textsNode = OBJECT_MAPPER.createArrayNode();
            if (answerTexts != null) {
                for (String t : answerTexts) {
                    if (t != null && !t.isBlank()) textsNode.add(t);
                }
            }

            ArrayNode ptsNode = OBJECT_MAPPER.createArrayNode();
            if (answerPoints != null && answerIds != null && answerPoints.size() == answerIds.size()) {
                for (Integer p : answerPoints) ptsNode.add(Math.max(0, p == null ? 0 : p));
            } else {
                int n = (answerIds == null) ? 0 : answerIds.size();
                for (int i = 0; i < n; i++) ptsNode.add(0);
                earnedPoints = 0;
            }

            for (JsonNode item : answers) {
                if (item != null && item.isObject()) {
                    String qid = item.path("questionId").asText(null);
                    if (q.equals(qid)) {
                        ObjectNode obj = (ObjectNode) item;

                        if (questionText != null) obj.put("questionText", questionText);
                        obj.put("maxPoints", Math.max(0, maxPoints));

                        obj.set("answerIds", idsNode);
                        obj.set("answerTexts", textsNode);
                        obj.set("answerPoints", ptsNode);
                        obj.put("earnedPoints", Math.max(0, earnedPoints));

                        obj.remove("answerId");

                        repository.updateAttemptVersionByAttemptId(attemptId, OBJECT_MAPPER.writeValueAsString(root));
                        return;
                    }
                }
            }

            ObjectNode newItem = OBJECT_MAPPER.createObjectNode();
            newItem.put("questionId", q);
            if (questionText != null) newItem.put("questionText", questionText);
            newItem.put("maxPoints", Math.max(0, maxPoints));
            newItem.set("answerIds", idsNode);
            newItem.set("answerTexts", textsNode);
            newItem.set("answerPoints", ptsNode);
            newItem.put("earnedPoints", Math.max(0, earnedPoints));
            answers.add(newItem);

            repository.updateAttemptVersionByAttemptId(attemptId, OBJECT_MAPPER.writeValueAsString(root));
        } catch (Exception e) {
            logger.error("Ошибка upsertAnswers", e);
            throw new RuntimeException("Ошибка upsertAnswers", e);
        }
    }

    /**
     * Завершить попытку по attemptId (для UI).
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

    private ObjectNode toObjectNodeOrEmpty(String json) {
        if (json == null || json.isBlank()) {
            ObjectNode root = OBJECT_MAPPER.createObjectNode();
            root.set("answers", OBJECT_MAPPER.createArrayNode());
            return root;
        }
        try {
            JsonNode node = OBJECT_MAPPER.readTree(json);
            if (node != null && node.isObject()) return (ObjectNode) node;
        } catch (Exception ignored) { }
        ObjectNode root = OBJECT_MAPPER.createObjectNode();
        root.set("answers", OBJECT_MAPPER.createArrayNode());
        return root;
    }

    private ArrayNode ensureAnswersArray(ObjectNode root) {
        JsonNode answers = root.get("answers");
        if (answers != null && answers.isArray()) return (ArrayNode) answers;
        ArrayNode arr = OBJECT_MAPPER.createArrayNode();
        root.set("answers", arr);
        return arr;
    }
}
