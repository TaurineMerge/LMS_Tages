package com.example.lms.ui;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;

import com.example.lms.question.api.dto.Question;
import com.example.lms.question.domain.service.QuestionService;

import com.example.lms.answer.api.dto.Answer;
import com.example.lms.answer.domain.service.AnswerService;

import com.example.lms.test_attempt.domain.service.TestAttemptService;

import io.javalin.http.Context;

import java.util.*;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * UI-контроллер прохождения теста (без JS).
 *
 * Функционал:
 * - single/multi определяется по количеству ответов со score>0
 * - если "Ответить" без выбора -> редирект на take с errorQ/errorType, показываем сообщение под вопросом
 * - "Завершить" без модалки: если есть пропуски -> ведём на /finish (страница подтверждения)
 *   если пропусков нет -> сразу POST /finish
 */
public class UiTestController {

    private static final String STUDENT_ID_COOKIE = "testing_student_id";
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private final TestService testService;
    private final QuestionService questionService;
    private final AnswerService answerService;
    private final TestAttemptService testAttemptService;

    public UiTestController(
            TestService testService,
            QuestionService questionService,
            AnswerService answerService,
            TestAttemptService testAttemptService
    ) {
        this.testService = testService;
        this.questionService = questionService;
        this.answerService = answerService;
        this.testAttemptService = testAttemptService;
    }

    // =========================================================
    // RETRY
    // =========================================================

    /**
     * POST /ui/tests/{testId}/retry
     */
    public void startNewAttempt(Context ctx) {
        UUID testId = parseUuidOr400(ctx, ctx.pathParam("testId"), "testId");
        if (testId == null) return;

        UUID studentId = resolveOrCreateStudentId(ctx);

        UUID newAttemptId = testAttemptService.createNewAttempt(studentId, testId);
        ctx.redirect("/testing/ui/tests/" + testId + "/take?attemptId=" + newAttemptId);
    }

    // =========================================================
    // TAKE
    // =========================================================

    /**
     * GET /ui/tests/{testId}/take?attemptId=...
     */
    public void showTakePage(Context ctx) {
        String testIdStr = ctx.pathParam("testId");
        UUID testId = parseUuidOr400(ctx, testIdStr, "testId");
        if (testId == null) return;

        Test testDto;
        try {
            testDto = testService.getTestById(testIdStr);
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Тест не найден: " + testIdStr);
            return;
        }

        UUID studentId = resolveOrCreateStudentId(ctx);

        UUID attemptId = tryParseUuid(ctx.queryParam("attemptId"));
        if (attemptId == null) {
            UUID newAttemptId = testAttemptService.createNewAttempt(studentId, testId);
            ctx.redirect("/testing/ui/tests/" + testId + "/take?attemptId=" + newAttemptId);
            return;
        }

        // параметры для серверной валидации без JS
        UUID errorQ = tryParseUuid(ctx.queryParam("errorQ"));
        String errorType = ctx.queryParam("errorType"); // single|multi

        List<Question> questions = questionService.getQuestionsByTestId(testId);

        int completedCount = testAttemptService.countCompletedAttempts(studentId, testId);
        int attemptNo = completedCount + 1;

        // init attempt_version, если пустой
        List<TestAttemptService.QuestionInit> init = new ArrayList<>();
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;
            init.add(new TestAttemptService.QuestionInit(qid, q.getOrder()));
        }
        testAttemptService.initAttemptVersionIfEmpty(attemptId, attemptNo, init);

        Map<String, List<String>> selectedAnswerByQuestion = loadSelectedAnswers(attemptId);

        List<Map<String, Object>> questionViews = new ArrayList<>();
        List<Map<String, Object>> missingQuestions = new ArrayList<>();
        int answeredCount = 0;

        for (Question q : questions) {
            UUID questionId = safeUuid(q.getId(), "вопроса", ctx);
            if (questionId == null) return;

            String qIdStr = questionId.toString();
            List<String> selectedAnswerIds = selectedAnswerByQuestion.getOrDefault(qIdStr, Collections.emptyList());
            boolean answered = selectedAnswerIds != null && !selectedAnswerIds.isEmpty();

            if (answered) answeredCount++;
            else {
                Map<String, Object> miss = new HashMap<>();
                miss.put("id", qIdStr);
                miss.put("order", q.getOrder());
                missingQuestions.add(miss);
            }

            List<Answer> answers = answerService.getAnswersByQuestionId(questionId);

            int maxPoints = 0;
            int positiveCount = 0;
            for (Answer a : answers) {
                int sc = extractAnswerScore(a);
                if (sc > 0) positiveCount++;
                maxPoints += Math.max(0, sc);
            }

            boolean multi = positiveCount > 1;
            String inputType = multi ? "checkbox" : "radio";
            String inputName = multi ? "answer_ids" : "answer_id";
            String hint = multi ? "Можно выбрать несколько вариантов" : "Выберите один вариант";

            List<Map<String, Object>> answerViews = new ArrayList<>();
            List<String> selectedTexts = new ArrayList<>();

            for (Answer a : answers) {
                String aIdStr = String.valueOf(a.getId());
                boolean selected = answered && selectedAnswerIds.contains(aIdStr);

                Map<String, Object> av = new HashMap<>();
                av.put("id", aIdStr);
                av.put("text", a.getText());
                av.put("selected", selected);
                av.put("disabled", answered);

                if (selected) selectedTexts.add(a.getText());
                answerViews.add(av);
            }

            Map<String, Object> qView = new HashMap<>();
            qView.put("id", qIdStr);
            qView.put("order", q.getOrder());
            qView.put("text_of_question", q.getTextOfQuestion());
            qView.put("answers", answerViews);

            qView.put("answered", answered);
            qView.put("selectedAnswerText", selectedTexts.isEmpty() ? null : String.join(", ", selectedTexts));
            qView.put("maxPoints", maxPoints);

            qView.put("multi", multi);
            qView.put("hint", hint);
            qView.put("inputType", inputType);
            qView.put("inputName", inputName);

            // Сообщение ошибки (без JS)
            if (!answered && errorQ != null && errorQ.equals(questionId)) {
                String msg = "multi".equalsIgnoreCase(errorType)
                        ? "Выберите хотя бы один вариант ответа."
                        : "Выберите один вариант ответа.";
                qView.put("validationError", msg);
            }

            qView.put("buttonText", answered ? "Ответ принят" : "Ответить");
            qView.put("buttonDisabled", answered);
            qView.put("buttonClass", answered ? "btn pressed" : "btn");

            questionViews.add(qView);
        }

        int totalQuestions = questionViews.size();
        int missingCount = missingQuestions.size();

        Map<String, Object> model = new HashMap<>();
        model.put("test", testDto);
        model.put("questions", questionViews);

        model.put("student_id", studentId.toString());
        model.put("attempt_id", attemptId.toString());
        model.put("attempt_no", attemptNo);

        model.put("totalQuestions", totalQuestions);
        model.put("answeredCount", answeredCount);

        model.put("missingCount", missingCount);
        model.put("missingQuestions", missingQuestions);
        model.put("allAnswered", missingCount == 0);

        String html = HbsRenderer.render("test-take", model);
        ctx.contentType("text/html; charset=utf-8");
        ctx.result(html);
    }

    /**
     * POST /ui/tests/{testId}/questions/{questionId}/answer
     * Без JS:
     * - если ничего не выбрали -> редирект назад на take с errorQ/errorType и якорем
     */
    public void submitAnswer(Context ctx) {
        UUID testId = parseUuidOr400(ctx, ctx.pathParam("testId"), "testId");
        if (testId == null) return;

        UUID questionId = parseUuidOr400(ctx, ctx.pathParam("questionId"), "questionId");
        if (questionId == null) return;

        UUID attemptId = tryParseUuid(ctx.formParam("attempt_id"));
        if (attemptId == null) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Нет attempt_id (hidden)");
            return;
        }

        // multi
        List<String> answerIdsRaw = ctx.formParams("answer_ids");
        // single
        String answerIdSingle = ctx.formParam("answer_id");

        boolean hasMulti = (answerIdsRaw != null && !answerIdsRaw.isEmpty());
        boolean hasSingle = (answerIdSingle != null && !answerIdSingle.isBlank());

        if (!hasMulti && !hasSingle) {
            boolean isMulti = isMultiQuestion(questionId);

            ctx.redirect("/testing/ui/tests/" + testId
                    + "/take?attemptId=" + attemptId
                    + "&errorQ=" + questionId
                    + "&errorType=" + (isMulti ? "multi" : "single")
                    + "#q_" + questionId);
            return;
        }

        List<String> chosenIds = hasMulti ? answerIdsRaw : List.of(answerIdSingle);

        List<UUID> answerIds = new ArrayList<>();
        for (String s : chosenIds) {
            UUID id = tryParseUuid(s);
            if (id == null) {
                ctx.status(400).contentType("text/plain; charset=utf-8")
                        .result("Некорректный UUID answerId: " + s);
                return;
            }
            answerIds.add(id);
        }

        resolveStudentIdForPost(ctx);

        try {
            testAttemptService.saveAnswers(attemptId, questionId, answerIds);
        } catch (Exception e) {
            ctx.status(500).contentType("text/plain; charset=utf-8")
                    .result("Ошибка при сохранении ответа: " + e.getMessage());
            return;
        }

        ctx.redirect("/testing/ui/tests/" + testId + "/take?attemptId=" + attemptId + "#q_" + questionId);
    }

    // =========================================================
    // FINISH
    // =========================================================

    /**
     * GET /ui/tests/{testId}/finish?attemptId=...
     * Страница подтверждения (без JS).
     */
    public void showFinishPage(Context ctx) {
        String testIdStr = ctx.pathParam("testId");
        UUID testId = parseUuidOr400(ctx, testIdStr, "testId");
        if (testId == null) return;

        UUID attemptId = tryParseUuid(ctx.queryParam("attemptId"));
        if (attemptId == null) {
            ctx.status(400).contentType("text/plain; charset=utf-8").result("Нет attemptId в query");
            return;
        }

        Test testDto;
        try { testDto = testService.getTestById(testIdStr); }
        catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8").result("Тест не найден");
            return;
        }

        UUID studentId = resolveOrCreateStudentId(ctx);

        List<Question> questions = questionService.getQuestionsByTestId(testId);
        Map<String, List<String>> selected = loadSelectedAnswers(attemptId);

        List<Map<String, Object>> missing = new ArrayList<>();
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            List<String> aids = selected.getOrDefault(qid.toString(), Collections.emptyList());
            if (aids == null || aids.isEmpty()) {
                Map<String, Object> m = new HashMap<>();
                m.put("id", qid.toString());
                m.put("order", q.getOrder());
                missing.add(m);
            }
        }

        int total = questions.size();
        int missingCount = missing.size();
        int answeredCount = total - missingCount;

        Map<String, Object> model = new HashMap<>();
        model.put("test", testDto);
        model.put("student_id", studentId.toString());
        model.put("attempt_id", attemptId.toString());

        model.put("totalQuestions", total);
        model.put("answeredCount", answeredCount);
        model.put("missingCount", missingCount);
        model.put("missingQuestions", missing);
        model.put("allAnswered", missingCount == 0);

        String html = HbsRenderer.render("test-finish", model);
        ctx.contentType("text/html; charset=utf-8");
        ctx.result(html);
    }

    /**
     * POST /ui/tests/{testId}/finish
     */
    public void finishAttempt(Context ctx) {
        UUID testId = parseUuidOr400(ctx, ctx.pathParam("testId"), "testId");
        if (testId == null) return;

        UUID attemptId = tryParseUuid(ctx.formParam("attempt_id"));
        if (attemptId == null) {
            ctx.status(400).contentType("text/plain; charset=utf-8").result("Нет attempt_id");
            return;
        }

        boolean force = "true".equalsIgnoreCase(ctx.formParam("force"));
        resolveStudentIdForPost(ctx);

        List<Question> questions = questionService.getQuestionsByTestId(testId);
        Map<String, List<String>> selected = loadSelectedAnswers(attemptId);

        int missingCount = 0;
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            List<String> aids = selected.getOrDefault(qid.toString(), Collections.emptyList());
            if (aids == null || aids.isEmpty()) missingCount++;
        }

        if (missingCount > 0 && !force) {
            ctx.redirect("/testing/ui/tests/" + testId + "/finish?attemptId=" + attemptId);
            return;
        }

        int points = 0;

        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            List<String> selectedAids = selected.getOrDefault(qid.toString(), Collections.emptyList());
            if (selectedAids == null || selectedAids.isEmpty()) continue;

            List<Answer> answers = answerService.getAnswersByQuestionId(qid);
            for (Answer a : answers) {
                if (selectedAids.contains(String.valueOf(a.getId()))) {
                    points += Math.max(0, extractAnswerScore(a));
                }
            }
        }

        try {
            testAttemptService.completeAttemptById(attemptId, points);
        } catch (Exception ignored) {}

        ctx.redirect("/testing/ui/tests/" + testId + "/results?attemptId=" + attemptId);
    }

    // =========================================================
    // Helpers
    // =========================================================

    private boolean isMultiQuestion(UUID questionId) {
        try {
            List<Answer> answers = answerService.getAnswersByQuestionId(questionId);
            int positiveCount = 0;
            for (Answer a : answers) {
                if (extractAnswerScore(a) > 0) positiveCount++;
                if (positiveCount > 1) return true;
            }
        } catch (Exception ignored) {}
        return false;
    }

    private int readAttemptNoFromJson(UUID attemptId) {
        try {
            Optional<String> jsonOpt = testAttemptService.getAttemptVersionByAttemptId(attemptId);
            if (jsonOpt.isEmpty()) return 0;

            String json = jsonOpt.get();
            if (json == null || json.isBlank()) return 0;

            JsonNode root = OBJECT_MAPPER.readTree(json);
            JsonNode n = root.get("attemptNo");
            if (n != null && n.isInt()) return n.asInt();
            if (n != null && n.isTextual()) {
                try { return Integer.parseInt(n.asText()); } catch (Exception ignored) {}
            }
        } catch (Exception ignored) {}
        return 0;
    }

    private Integer extractMinPoint(Test testDto) {
        try {
            Object v = testDto.getClass().getMethod("getMin_point").invoke(testDto);
            return (v == null) ? null : ((Number) v).intValue();
        } catch (Exception ignored1) {
            try {
                Object v = testDto.getClass().getMethod("getMinPoint").invoke(testDto);
                return (v == null) ? null : ((Number) v).intValue();
            } catch (Exception ignored2) {
                return null;
            }
        }
    }

    private Map<String, List<String>> loadSelectedAnswers(UUID attemptId) {
        Map<String, List<String>> map = new HashMap<>();
        try {
            Optional<String> jsonOpt = testAttemptService.getAttemptVersionByAttemptId(attemptId);
            if (jsonOpt.isEmpty()) return map;

            String json = jsonOpt.get();
            if (json == null || json.isBlank()) return map;

            JsonNode root = OBJECT_MAPPER.readTree(json);
            JsonNode answers = root.path("answers");
            if (!answers.isArray()) return map;

            for (JsonNode item : answers) {
                if (item == null || !item.isObject()) continue;

                String qid = item.path("questionId").asText(null);
                if (qid == null || qid.isBlank()) continue;

                JsonNode arr = item.get("answerIds");
                if (arr != null && arr.isArray() && arr.size() > 0) {
                    List<String> ids = new ArrayList<>();
                    for (JsonNode n : arr) {
                        String v = n.asText(null);
                        if (v != null && !v.isBlank()) ids.add(v);
                    }
                    if (!ids.isEmpty()) {
                        map.put(qid, ids);
                        continue;
                    }
                }

                JsonNode single = item.get("answerId");
                if (single != null && !single.isNull()) {
                    String aid = single.asText(null);
                    if (aid != null && !aid.isBlank()) {
                        map.put(qid, List.of(aid));
                    }
                }
            }
        } catch (Exception ignored) {}
        return map;
    }

    private UUID resolveOrCreateStudentId(Context ctx) {
        String fromJwt = ctx.attribute("userId");
        UUID jwtUuid = tryParseUuid(fromJwt);
        if (jwtUuid != null) {
            ctx.cookie(STUDENT_ID_COOKIE, jwtUuid.toString());
            return jwtUuid;
        }

        String fromQuery = ctx.queryParam("studentId");
        UUID queryUuid = tryParseUuid(fromQuery);
        if (queryUuid != null) {
            ctx.cookie(STUDENT_ID_COOKIE, queryUuid.toString());
            return queryUuid;
        }

        String fromCookie = ctx.cookie(STUDENT_ID_COOKIE);
        UUID cookieUuid = tryParseUuid(fromCookie);
        if (cookieUuid != null) return cookieUuid;

        UUID generated = UUID.randomUUID();
        ctx.cookie(STUDENT_ID_COOKIE, generated.toString());
        return generated;
    }

    private UUID resolveStudentIdForPost(Context ctx) {
        String fromForm = ctx.formParam("student_id");
        UUID formUuid = tryParseUuid(fromForm);
        if (formUuid != null) {
            ctx.cookie(STUDENT_ID_COOKIE, formUuid.toString());
            return formUuid;
        }
        return resolveOrCreateStudentId(ctx);
    }

    private UUID tryParseUuid(String value) {
        if (value == null || value.isBlank()) return null;
        try { return UUID.fromString(value.trim()); }
        catch (IllegalArgumentException e) { return null; }
    }

    private UUID parseUuidOr400(Context ctx, String raw, String field) {
        UUID id = tryParseUuid(raw);
        if (id == null) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID " + field + ": " + raw);
        }
        return id;
    }

    private UUID safeUuid(Object idObj, String entity, Context ctx) {
        try {
            if (idObj instanceof UUID) return (UUID) idObj;
            return UUID.fromString(String.valueOf(idObj));
        } catch (Exception e) {
            ctx.status(500).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID " + entity + ": " + idObj);
            return null;
        }
    }

    private int extractAnswerScore(Answer a) {
        try {
            Object v = a.getClass().getMethod("getScore").invoke(a);
            if (v == null) return 0;
            return ((Number) v).intValue();
        } catch (Exception ignored) {
            return 0;
        }
    }

    /**
     * GET /ui/tests/{testId}/results?attemptId=...
     */
    public void showResultsPage(Context ctx) {
        String testIdStr = ctx.pathParam("testId");
        UUID testId = parseUuidOr400(ctx, testIdStr, "testId");
        if (testId == null) return;

        UUID attemptId = tryParseUuid(ctx.queryParam("attemptId"));
        if (attemptId == null) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Нет attemptId в query");
            return;
        }

        Test testDto;
        try {
            testDto = testService.getTestById(testIdStr);
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Тест не найден");
            return;
        }

        UUID studentId = resolveOrCreateStudentId(ctx);

        List<Question> questions = questionService.getQuestionsByTestId(testId);
        Map<String, List<String>> selected = loadSelectedAnswers(attemptId);

        int points = 0;
        int maxPossible = 0;
        int answeredCount = 0;

        List<Map<String, Object>> rows = new ArrayList<>();

        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            List<Answer> answers = answerService.getAnswersByQuestionId(qid);

            int qMax = 0;
            for (Answer a : answers) qMax += Math.max(0, extractAnswerScore(a));
            maxPossible += qMax;

            List<String> selectedAids = selected.getOrDefault(qid.toString(), Collections.emptyList());

            int earned = 0;
            List<String> selectedTexts = new ArrayList<>();

            if (selectedAids != null && !selectedAids.isEmpty()) {
                answeredCount++;
                for (Answer a : answers) {
                    if (selectedAids.contains(String.valueOf(a.getId()))) {
                        earned += Math.max(0, extractAnswerScore(a));
                        selectedTexts.add(a.getText());
                    }
                }
                points += earned;
            }

            Map<String, Object> row = new HashMap<>();
            row.put("order", q.getOrder());
            row.put("text", q.getTextOfQuestion());
            row.put("answered", selectedAids != null && !selectedAids.isEmpty());
            row.put("selectedText", selectedTexts.isEmpty() ? null : String.join(", ", selectedTexts));
            row.put("earned", earned);
            row.put("maxPoints", qMax);
            rows.add(row);
        }

        Integer minPoint = extractMinPoint(testDto);
        boolean passed = (minPoint != null) && points >= minPoint;

        int attemptNo = readAttemptNoFromJson(attemptId);

        Map<String, Object> model = new HashMap<>();
        model.put("test", testDto);
        model.put("student_id", studentId.toString());
        model.put("attempt_id", attemptId.toString());
        model.put("attempt_no", attemptNo);

        model.put("points", points);
        model.put("maxPossible", maxPossible);

        model.put("totalQuestions", questions.size());
        model.put("answeredCount", answeredCount);

        model.put("minPoint", minPoint);
        model.put("passed", passed);

        model.put("rows", rows);

        String html = HbsRenderer.render("test-results", model);
        ctx.contentType("text/html; charset=utf-8");
        ctx.result(html);
    }

}