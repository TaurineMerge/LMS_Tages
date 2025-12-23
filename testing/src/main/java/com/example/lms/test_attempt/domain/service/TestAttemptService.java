package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.ArrayList;
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

/**
 * Сервис для попыток тестов.
 *
 * UI-логика:
 * - attempt_version хранится по attemptId (PK)
 * - attempt_version инициализируется сразу всеми вопросами
 * - attemptNo пишем в JSON
 *
 * Поддержка multi-answer:
 * - в JSON используем "answerIds": ["...","..."]
 * - для обратной совместимости оставляем "answerId" если выбран ровно 1 ответ
 */
public class TestAttemptService {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private final TestAttemptRepositoryInterface repository;

    public TestAttemptService(TestAttemptRepositoryInterface repository) {
        this.repository = repository;
    }

    // =========================================================
    // DTO <-> MODEL (HEAD)
    // =========================================================

    private TestAttemptModel toModel(TestAttempt dto) {
        LocalDate date = null;
        if (dto.getDate_of_attempt() != null && !dto.getDate_of_attempt().isEmpty()) {
            try {
                date = LocalDate.parse(dto.getDate_of_attempt());
            } catch (Exception e) {
                logger.warn("Ошибка парсинга даты: {}", dto.getDate_of_attempt());
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
                dto.getCompleted()
        );
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
                model.getCompleted()
        );
    }

    // =========================================================
    // CRUD / Queries (HEAD)
    // =========================================================

    public List<TestAttempt> getAllTestAttempts() {
        return repository.findAll().stream().map(this::toDto).toList();
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
        // m.setDateOfAttempt(LocalDate.now()); // ❌ НЕ ДЕЛАЕМ
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
     *     {"order":1,"questionId":"...","answerId":null,"answerIds":[]},
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
                item.putNull("answerId");               // legacy
                item.set("answerIds", OBJECT_MAPPER.createArrayNode()); // multi
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
     * Старый single API: сохраняем 1 ответ.
     * Для совместимости оставляем метод, но внутри пишем через saveAnswers().
     */
    public void saveAnswer(UUID attemptId, UUID questionId, UUID answerId) {
        saveAnswers(attemptId, questionId, List.of(answerId));
    }

    /**
     * Multi API: сохраняем выбранные ответы.
     * Если уже есть ответ(ы) — повтор игнорируем.
     */
    public void saveAnswers(UUID attemptId, UUID questionId, List<UUID> answerIds) {
        if (answerIds == null || answerIds.isEmpty()) {
            throw new IllegalArgumentException("answerIds пустой");
        }

        try {
            String json = repository.findAttemptVersionByAttemptId(attemptId).orElse(null);
            ObjectNode root = toObjectNodeOrEmpty(json);
            ArrayNode answers = ensureAnswersArray(root);

            String q = questionId.toString();

            // нормализуем: без null + без дублей, но сохраняем порядок
            List<String> ids = new ArrayList<>();
            for (UUID id : answerIds) {
                if (id == null) continue;
                String s = id.toString();
                if (!ids.contains(s)) ids.add(s);
            }
            if (ids.isEmpty()) {
                throw new IllegalArgumentException("answerIds после нормализации пустой");
            }

            for (JsonNode item : answers) {
                if (item != null && item.isObject()) {
                    String qid = item.path("questionId").asText(null);
                    if (qid != null && qid.equals(q)) {

                        // если уже есть выбранные — игнорируем повтор (как у тебя было)
                        if (hasAnySelected((ObjectNode) item)) {
                            logger.info("Ответ уже был выбран (attemptId={}, questionId={}) — игнор", attemptId, questionId);
                            return;
                        }

                        applySelectedAnswers((ObjectNode) item, ids);

                        String updated = OBJECT_MAPPER.writeValueAsString(root);
                        repository.updateAttemptVersionByAttemptId(attemptId, updated);
                        return;
                    }
                }
            }

            // если вопроса вдруг нет в массиве — добавим (на всякий)
            ObjectNode newItem = OBJECT_MAPPER.createObjectNode();
            newItem.put("questionId", q);
            applySelectedAnswers(newItem, ids);
            answers.add(newItem);

            String updated = OBJECT_MAPPER.writeValueAsString(root);
            repository.updateAttemptVersionByAttemptId(attemptId, updated);

        } catch (Exception e) {
            logger.error("Ошибка saveAnswers", e);
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

    private boolean hasAnySelected(ObjectNode item) {
        // multi
        JsonNode arr = item.get("answerIds");
        if (arr != null && arr.isArray() && arr.size() > 0) return true;

        // legacy single
        JsonNode single = item.get("answerId");
        if (single != null && !single.isNull() && !single.asText("").isBlank()) return true;

        return false;
    }

    private void applySelectedAnswers(ObjectNode item, List<String> ids) {
        // legacy: если ровно 1 — кладём в answerId, иначе null
        if (ids.size() == 1) item.put("answerId", ids.get(0));
        else item.putNull("answerId");

        ArrayNode arr = OBJECT_MAPPER.createArrayNode();
        for (String id : ids) arr.add(id);
        item.set("answerIds", arr);
    }

    // ---------------------------------------------------------------------
    // BACKWARD COMPATIBILITY (for TestAttemptController)
    // ---------------------------------------------------------------------

    /**
     * Старый API — завершение попытки по ID (используется API-контроллером).
     * ДЛЯ ОБРАТНОЙ СОВМЕСТИМОСТИ.
     */
    public TestAttempt completeTestAttempt(String id, Integer finalPoint) {
        UUID attemptId = UUID.fromString(id);

        TestAttemptModel model = repository.findById(attemptId)
                .orElseThrow(() -> new IllegalArgumentException("Attempt not found: " + id));

        model.completeAttempt(finalPoint);
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }

    /**
     * Старый API — обновление snapshot (используется API-контроллером).
     * ДЛЯ ОБРАТНОЙ СОВМЕСТИМОСТИ.
     */
    public TestAttempt updateSnapshot(String id, String snapshot) {
        UUID attemptId = UUID.fromString(id);

        TestAttemptModel model = repository.findById(attemptId)
                .orElseThrow(() -> new IllegalArgumentException("Attempt not found: " + id));

        model.updateSnapshot(snapshot);
        TestAttemptModel updated = repository.update(model);
        return toDto(updated);
    }
}
