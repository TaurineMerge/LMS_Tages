package com.example.lms.ui;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;

import com.example.lms.question.api.dto.Question;
import com.example.lms.question.domain.service.QuestionService;

import com.example.lms.answer.api.dto.Answer;
import com.example.lms.answer.domain.service.AnswerService;

import com.example.lms.test_attempt.domain.service.TestAttemptService;

import io.javalin.http.Context;

import java.time.LocalDate;
import java.util.*;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * UI-контроллер страницы прохождения теста + завершение + результаты.
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
    // TAKE PAGE
    // =========================================================

    public void showTakePage(Context ctx) {
        String testIdStr = ctx.pathParam("testId");

        UUID testIdUuid;
        try {
            testIdUuid = UUID.fromString(testIdStr);
        } catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID testId: " + testIdStr);
            return;
        }

        Test testDto;
        try {
            testDto = testService.getTestById(testIdStr); // у тебя getTestById(String)
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Тест не найден: " + testIdStr);
            return;
        }

        UUID studentId = resolveOrCreateStudentId(ctx);

        Map<String, String> selectedAnswerByQuestion =
                loadSelectedAnswers(studentId, testIdUuid, LocalDate.now());

        List<Question> questions = questionService.getQuestionsByTestId(testIdUuid);

        List<Map<String, Object>> questionViews = new ArrayList<>();
        List<Map<String, Object>> missingQuestions = new ArrayList<>();

        int answeredCount = 0;

        for (Question q : questions) {

            UUID questionId;
            try {
                Object qIdObj = q.getId();
                if (qIdObj instanceof UUID) {
                    questionId = (UUID) qIdObj;
                } else {
                    questionId = UUID.fromString(String.valueOf(qIdObj));
                }
            } catch (Exception ex) {
                ctx.status(500).contentType("text/plain; charset=utf-8")
                        .result("Некорректный UUID вопроса: " + q.getId());
                return;
            }

            String qIdStr = questionId.toString();
            String selectedAnswerId = selectedAnswerByQuestion.get(qIdStr);
            boolean answered = (selectedAnswerId != null);

            if (answered) {
                answeredCount++;
            } else {
                Map<String, Object> miss = new HashMap<>();
                miss.put("id", qIdStr);
                miss.put("order", q.getOrder());
                missingQuestions.add(miss);
            }

            List<Answer> answers = answerService.getAnswersByQuestionId(questionId);

            int maxPoints = 0;
            for (Answer a : answers) {
                try {
                    Integer score = a.getScore(); // если нет getScore() — скажи, поправлю
                    maxPoints += (score == null ? 0 : score);
                } catch (Exception ignored) {}
            }

            List<Map<String, Object>> answerViews = new ArrayList<>();
            String selectedAnswerText = null;

            for (Answer a : answers) {
                String aIdStr = String.valueOf(a.getId());
                boolean selected = answered && aIdStr.equals(selectedAnswerId);

                Map<String, Object> av = new HashMap<>();
                av.put("id", aIdStr);
                av.put("text", a.getText());
                av.put("selected", selected);
                av.put("disabled", answered);

                if (selected) selectedAnswerText = a.getText();

                answerViews.add(av);
            }

            Map<String, Object> qView = new HashMap<>();
            qView.put("id", qIdStr);
            qView.put("order", q.getOrder());
            qView.put("text_of_question", q.getTextOfQuestion());
            qView.put("answers", answerViews);

            qView.put("answered", answered);
            qView.put("selectedAnswerText", selectedAnswerText);

            qView.put("maxPoints", maxPoints);

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

        model.put("totalQuestions", totalQuestions);
        model.put("answeredCount", answeredCount);

        // ✅ для всплывающего окна
        model.put("missingCount", missingCount);
        model.put("missingQuestions", missingQuestions);
        model.put("allAnswered", missingCount == 0);

        String html = HbsRenderer.render("test-take", model);
        ctx.contentType("text/html; charset=utf-8");
        ctx.result(html);
    }


    public void submitAnswer(Context ctx) {
        String testIdStr = ctx.pathParam("testId");
        String questionIdStr = ctx.pathParam("questionId");
        String answerIdStr = ctx.formParam("answer_id");

        if (answerIdStr == null || answerIdStr.isBlank()) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Не выбран ответ (answer_id пустой)");
            return;
        }

        UUID testId;
        UUID questionId;
        UUID answerId;

        try { testId = UUID.fromString(testIdStr); }
        catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID testId: " + testIdStr);
            return;
        }

        try { questionId = UUID.fromString(questionIdStr); }
        catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID questionId: " + questionIdStr);
            return;
        }

        try { answerId = UUID.fromString(answerIdStr); }
        catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID answerId: " + answerIdStr);
            return;
        }

        UUID studentId = resolveStudentIdForPost(ctx);

        try {
            testAttemptService.saveAnswer(studentId, testId, questionId, answerId);
        } catch (Exception e) {
            ctx.status(500).contentType("text/plain; charset=utf-8")
                    .result("Ошибка при сохранении ответа: " + e.getMessage());
            return;
        }

        // якорь — чтобы не скроллило в начало
        ctx.redirect("/testing/ui/tests/" + testId + "/take#q_" + questionId);
    }

    // =========================================================
    // FINISH FLOW
    // =========================================================

    /**
     * GET /ui/tests/{testId}/finish
     * Показывает подтверждение, если есть неотвеченные вопросы.
     */
    public void showFinishPage(Context ctx) {
        String testIdStr = ctx.pathParam("testId");

        UUID testIdUuid;
        try {
            testIdUuid = UUID.fromString(testIdStr);
        } catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID testId: " + testIdStr);
            return;
        }

        Test testDto;
        try {
            testDto = testService.getTestById(testIdStr);
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Тест не найден: " + testIdStr);
            return;
        }

        UUID studentId = resolveOrCreateStudentId(ctx);

        List<Question> questions = questionService.getQuestionsByTestId(testIdUuid);
        Map<String, String> selected = loadSelectedAnswers(studentId, testIdUuid, LocalDate.now());

        List<Map<String, Object>> missing = new ArrayList<>();
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            if (!selected.containsKey(qid.toString())) {
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
     * Если force=true — завершает даже при незаполненных вопросах.
     * После завершения — редирект на результаты.
     */
    public void finishAttempt(Context ctx) {
        String testIdStr = ctx.pathParam("testId");

        UUID testIdUuid;
        try {
            testIdUuid = UUID.fromString(testIdStr);
        } catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID testId: " + testIdStr);
            return;
        }

        UUID studentId = resolveStudentIdForPost(ctx);
        boolean force = "true".equalsIgnoreCase(ctx.formParam("force"));

        List<Question> questions = questionService.getQuestionsByTestId(testIdUuid);
        Map<String, String> selected = loadSelectedAnswers(studentId, testIdUuid, LocalDate.now());

        // проверяем незаполненные
        int missingCount = 0;
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;
            if (!selected.containsKey(qid.toString())) missingCount++;
        }

        if (missingCount > 0 && !force) {
            // если пришли без force — отправляем на confirm-страницу
            ctx.redirect("testing/ui/tests/" + testIdUuid + "/finish");
            return;
        }

        // считаем баллы по выбранным ответам
        int points = 0;
        int maxPossible = 0;

        List<Map<String, Object>> resultRows = new ArrayList<>();

        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            String qidStr = qid.toString();
            String selectedAid = selected.get(qidStr);

            List<Answer> answers = answerService.getAnswersByQuestionId(qid);

            // maxPossible по твоей логике "сумма score всех ответов"
            int qMax = 0;
            for (Answer a : answers) qMax += Math.max(0, extractAnswerScore(a));
            maxPossible += qMax;

            int earned = 0;
            String selectedText = null;

            if (selectedAid != null) {
                for (Answer a : answers) {
                    String aid = String.valueOf(a.getId());
                    if (aid.equals(selectedAid)) {
                        earned = Math.max(0, extractAnswerScore(a));
                        selectedText = a.getText();
                        break;
                    }
                }
                points += earned;
            }

            Map<String, Object> row = new HashMap<>();
            row.put("order", q.getOrder());
            row.put("text", q.getTextOfQuestion());
            row.put("answered", selectedAid != null);
            row.put("selectedText", selectedText);
            row.put("earned", earned);
            row.put("maxPoints", qMax);
            resultRows.add(row);
        }

        // сохраняем point в попытку (на сегодня)
        try {
            testAttemptService.completeAttemptForToday(studentId, testIdUuid, points);
        } catch (Exception ignored) {
            // если вдруг не сохранится — результаты всё равно показываем
        }

        // кладем результаты в cookie? не нужно — просто редирект на results (там пересчитаем)
        ctx.redirect("/testing/ui/tests/" + testIdUuid + "/results");
    }

    /**
     * GET /ui/tests/{testId}/results
     */
    public void showResultsPage(Context ctx) {
        String testIdStr = ctx.pathParam("testId");

        UUID testIdUuid;
        try {
            testIdUuid = UUID.fromString(testIdStr);
        } catch (IllegalArgumentException e) {
            ctx.status(400).contentType("text/plain; charset=utf-8")
                    .result("Некорректный UUID testId: " + testIdStr);
            return;
        }

        Test testDto;
        try {
            testDto = testService.getTestById(testIdStr);
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Тест не найден: " + testIdStr);
            return;
        }

        UUID studentId = resolveOrCreateStudentId(ctx);

        List<Question> questions = questionService.getQuestionsByTestId(testIdUuid);
        Map<String, String> selected = loadSelectedAnswers(studentId, testIdUuid, LocalDate.now());

        int points = 0;
        int maxPossible = 0;
        int answeredCount = 0;

        List<Map<String, Object>> rows = new ArrayList<>();

        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            String selectedAid = selected.get(qid.toString());
            List<Answer> answers = answerService.getAnswersByQuestionId(qid);

            int qMax = 0;
            for (Answer a : answers) qMax += Math.max(0, extractAnswerScore(a));
            maxPossible += qMax;

            int earned = 0;
            String selectedText = null;

            if (selectedAid != null) {
                answeredCount++;
                for (Answer a : answers) {
                    if (String.valueOf(a.getId()).equals(selectedAid)) {
                        earned = Math.max(0, extractAnswerScore(a));
                        selectedText = a.getText();
                        break;
                    }
                }
                points += earned;
            }

            Map<String, Object> row = new HashMap<>();
            row.put("order", q.getOrder());
            row.put("text", q.getTextOfQuestion());
            row.put("answered", selectedAid != null);
            row.put("selectedText", selectedText);
            row.put("earned", earned);
            row.put("maxPoints", qMax);
            rows.add(row);
        }

        Integer minPoint = extractMinPoint(testDto);
        boolean passed = (minPoint == null) ? false : points >= minPoint;

        Map<String, Object> model = new HashMap<>();
        model.put("test", testDto);
        model.put("student_id", studentId.toString());

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

    // =========================================================
    // Helpers
    // =========================================================

    private Map<String, String> loadSelectedAnswers(UUID studentId, UUID testId, LocalDate date) {
        Map<String, String> map = new HashMap<>();
        try {
            Optional<String> jsonOpt = testAttemptService.getAttemptVersion(studentId, testId, date);
            if (jsonOpt.isEmpty() || jsonOpt.get() == null || jsonOpt.get().isBlank()) return map;

            JsonNode root = OBJECT_MAPPER.readTree(jsonOpt.get());
            JsonNode answers = root.path("answers");

            // формат: answers: [{question, answer}, ...]
            if (answers.isArray()) {
                for (JsonNode item : answers) {
                    String q = item.path("question").asText(null);
                    String a = item.path("answer").asText(null);
                    if (q != null && a != null) map.put(q, a);
                }
                return map;
            }

            // fallback на старый объект-словарь
            if (answers.isObject()) {
                Iterator<Map.Entry<String, JsonNode>> it = answers.fields();
                while (it.hasNext()) {
                    var e = it.next();
                    map.put(e.getKey(), e.getValue() == null ? null : e.getValue().asText());
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

    /**
     * Достаём score максимально безопасно (на случай если DTO названо иначе).
     */
    private int extractAnswerScore(Answer a) {
        try {
            Object v = a.getClass().getMethod("getScore").invoke(a);
            if (v == null) return 0;
            return ((Number) v).intValue();
        } catch (Exception ignored) {
            // если нет метода getScore — будет 0 (и не упадёт компиляция)
            return 0;
        }
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
}