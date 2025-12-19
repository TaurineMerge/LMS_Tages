package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ArrayNode;
import com.fasterxml.jackson.databind.node.ObjectNode;

/**
 * Сервис для работы с попытками прохождения тестов.
 */
public class TestAttemptService {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private final TestAttemptRepositoryInterface testAttemptRepository;

    public TestAttemptService(TestAttemptRepositoryInterface testAttemptRepository) {
        this.testAttemptRepository = testAttemptRepository;
    }

    /**
     * ✅ НОВОЕ: получить attempt_version для UI.
     */
    public Optional<String> getAttemptVersion(UUID studentId, UUID testId, LocalDate date) {
        return testAttemptRepository.findAttemptVersion(studentId, testId, date);
    }

    /**
     * Сохраняет ответ студента на конкретный вопрос в attempt_version (JSON) таблицы test_attempt_b.
     *
     * Формат attempt_version:
     * {
     *   "answers": [
     *     { "question": "<questionUuid>", "answer": "<answerUuid>" },
     *     ...
     *   ]
     * }
     *
     * Правило "без возможности исправления":
     * если в массиве answers уже есть объект с question == questionId, второй раз НЕ добавляем/не меняем.
     */
    public void saveAnswer(UUID studentId, UUID testId, UUID questionId, UUID answerId) {
        LocalDate today = LocalDate.now();

        try {
            String existingJson = testAttemptRepository
                    .findAttemptVersion(studentId, testId, today)
                    .orElse(null);

            ObjectNode root = toObjectNodeOrEmpty(existingJson);
            ArrayNode answersArray = ensureAnswersArray(root);

            String q = questionId.toString();
            String a = answerId.toString();

            // уже отвечал — не даём менять
            for (JsonNode item : answersArray) {
                if (item != null && item.isObject()) {
                    JsonNode qNode = item.get("question");
                    if (qNode != null && q.equals(qNode.asText())) {
                        logger.info("Ответ уже сохранён (student={}, test={}, date={}, question={}). Повтор игнорируем.",
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
            testAttemptRepository.upsertAttemptVersion(studentId, testId, today, updatedJson);

            logger.info("Ответ сохранён (student={}, test={}, date={}, question={}, answer={})",
                    studentId, testId, today, questionId, answerId);

        } catch (Exception e) {
            logger.error("Ошибка при сохранении ответа в attempt_version", e);
            throw new RuntimeException("Ошибка при сохранении ответа", e);
        }
    }

    private ArrayNode ensureAnswersArray(ObjectNode root) {
        JsonNode answers = root.get("answers");

        if (answers != null && answers.isArray()) {
            return (ArrayNode) answers;
        }

        // если вдруг лежал старый формат объектом — конвертируем в массив
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
        } catch (Exception ignored) {}
        return OBJECT_MAPPER.createObjectNode();
    }

    // ---- твои существующие методы ниже (без изменений) ----

    public TestAttemptModel createTestAttempt(UUID studentId, UUID testId, String attemptVersion) {
        TestAttemptModel testAttempt = new TestAttemptModel(studentId, testId);
        testAttempt.validate();
        return testAttemptRepository.save(testAttempt);
    }

    public TestAttemptModel getTestAttemptById(UUID id) {
        return testAttemptRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("No test attempt found with id: " + id));
    }

    public TestAttemptModel updateTestAttempt(TestAttemptModel testAttempt) {
        UUID id = testAttempt.getId();
        if (id == null) {
            throw new IllegalArgumentException("Test Attempt ID cannot be null for update");
        }

        TestAttemptModel existingAttempt = testAttemptRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("No test attempt found with id: " + id));

        existingAttempt.setPoint(testAttempt.getPoint());
        existingAttempt.validate();

        return testAttemptRepository.update(existingAttempt);
    }

    public void completeTestAttempt(UUID id, int points) {
        TestAttemptModel testAttempt = getTestAttemptById(id);
        testAttempt.complete(points);
        testAttempt.validate();
        testAttemptRepository.update(testAttempt);
    }

    public void attachCertificate(UUID id, UUID certificateId) {
        TestAttemptModel testAttempt = getTestAttemptById(id);
        testAttempt.validate();
        testAttemptRepository.update(testAttempt);
    }

    public void deleteTestAttempt(UUID id) {
        testAttemptRepository.deleteById(id);
    }

    public List<TestAttemptModel> getAllTestAttempts() {
        return testAttemptRepository.findAll();
    }

    public Optional<TestAttemptModel> findById(UUID id) {
        return testAttemptRepository.findById(id);
    }

    public List<TestAttemptModel> getTestAttemptsByStudentId(UUID studentId) {
        return testAttemptRepository.findByStudentId(studentId);
    }

    public List<TestAttemptModel> getTestAttemptsByTestId(UUID testId) {
        return testAttemptRepository.findByTestId(testId);
    }

    public List<TestAttemptModel> getTestAttemptsByDate(LocalDate date) {
        return testAttemptRepository.findByDate(date);
    }

    public List<TestAttemptModel> getCompletedTestAttempts() {
        return testAttemptRepository.findCompletedAttempts();
    }

    public List<TestAttemptModel> getIncompleteTestAttempts() {
        return testAttemptRepository.findIncompleteAttempts();
    }

    public List<TestAttemptModel> getAttemptsByStudentAndTest(UUID studentId, UUID testId) {
        return testAttemptRepository.findByStudentAndTest(studentId, testId);
    }

    public int countAttemptsByStudentAndTest(UUID studentId, UUID testId) {
        return testAttemptRepository.countAttemptsByStudentAndTest(studentId, testId);
    }

    public boolean existsById(UUID id) {
        return testAttemptRepository.existsById(id);
    }

    public void completeAttemptForToday(UUID studentId, UUID testId, int points) {
        var today = java.time.LocalDate.now();

        // ищем попытку за сегодня
        var attempts = testAttemptRepository.findByStudentAndTest(studentId, testId);
        com.example.lms.test_attempt.domain.model.TestAttemptModel todayAttempt = null;

        for (var a : attempts) {
            if (today.equals(a.getDateOfAttempt())) {
                todayAttempt = a;
                break;
            }
        }

        // если нет — создаём пустую попытку
        if (todayAttempt == null) {
            todayAttempt = new com.example.lms.test_attempt.domain.model.TestAttemptModel(studentId, testId);
            todayAttempt.validate();
            todayAttempt = testAttemptRepository.save(todayAttempt);
        }

        todayAttempt.setPoint(points);
        todayAttempt.validate();
        testAttemptRepository.update(todayAttempt);
    }

}