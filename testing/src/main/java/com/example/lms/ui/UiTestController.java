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
 * Актуальная логика (UI):
 * - одна форма на весь тест
 * - вверху кнопки: «Сохранить» и «Завершить»
 * - «Сохранить» сохраняет все текущие ответы, можно исправлять (перезапись)
 * - «Завершить» тоже сначала сохраняет всё, затем:
 *    - если есть пропуски -> редирект на /finish (страница подтверждения)
 *    - если нет пропусков -> завершает и ведёт на /results
 *
 * Важно для устойчивости результатов:
 * - результаты отрисовываем ТОЛЬКО по attempt_version (attempt_snapshot не трогаем)
 * - attempt_version хранит тексты вопросов/ответов + maxPoints/earnedPoints
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
    // COURSE ENTRYPOINT (category/course)
    // =========================================================

    /**
     * "Умная" точка входа для кнопки "Пройти тест" на странице курса.
     *
     * Логика:
     * - если есть незавершённая попытка (completed=false) по этому тесту -> редиректим на неё
     * - иначе если есть завершённая попытка (completed=true) -> открываем результат последней завершённой
     * - иначе создаём новую попытку и ведём в неё
     *
     * GET /ui/category/{categoryId}/course/{courseId}/test
     */
    public void startOrResumeFromCourse(Context ctx) {
        String categoryId = ctx.pathParam("categoryId");
        String courseIdStr = ctx.pathParam("courseId");

        // 1) Находим тест по courseId (в проекте курс может иметь только 1 тест)
        List<Test> tests;
        try {
            tests = testService.getTestsByCourseId(courseIdStr);
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Тест для курса не найден: courseId=" + courseIdStr);
            return;
        }

        if (tests == null || tests.isEmpty() || tests.get(0).getId() == null) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("У курса нет теста: courseId=" + courseIdStr);
            return;
        }

        UUID testId = tryParseUuid(String.valueOf(tests.get(0).getId()));
        if (testId == null) {
            ctx.status(500).contentType("text/plain; charset=utf-8")
                    .result("Некорректный testId у курса: courseId=" + courseIdStr);
            return;
        }

        // 2) studentId (пока cookie, позже интегрируем JWT)
        UUID studentId = resolveOrCreateStudentId(ctx);

        // 3) Если есть незавершённая попытка -> продолжаем её
        UUID attemptId = testAttemptService.getLatestIncompleteAttempt(studentId, testId)
                .map(a -> tryParseUuid(a.getId()))
                .orElse(null);

        if (attemptId != null) {
            ctx.redirect("/testing/ui/category/" + categoryId + "/course/" + courseIdStr + "/attempt/" + attemptId);
            return;
        }

        // 4) Иначе показываем результат последней завершённой (если есть)
        UUID lastCompletedAttemptId = testAttemptService.getLatestCompletedAttempt(studentId, testId)
                .map(a -> tryParseUuid(a.getId()))
                .orElse(null);

        if (lastCompletedAttemptId != null) {
            ctx.redirect("/testing/ui/tests/" + testId + "/results?attemptId=" + lastCompletedAttemptId);
            return;
        }

        // 5) Иначе создаём новую попытку
        UUID newAttemptId = testAttemptService.createNewAttempt(studentId, testId);
        ctx.redirect("/testing/ui/category/" + categoryId + "/course/" + courseIdStr + "/attempt/" + newAttemptId);
    }

    /**
     * Открыть конкретную попытку в UI.
     *
     * GET /ui/category/{categoryId}/course/{courseId}/attempt/{attemptId}
     */
    public void openAttemptFromCourse(Context ctx) {
        String attemptIdStr = ctx.pathParam("attemptId");
        UUID attemptId = parseUuidOr400(ctx, attemptIdStr, "attemptId");
        if (attemptId == null) return;

        try {
            // Берём testId из попытки, чтобы переиспользовать существующие страницы
            com.example.lms.test_attempt.api.dto.TestAttempt attempt = testAttemptService.getTestAttemptById(attemptId);
            UUID testId = tryParseUuid(attempt.getTest_id());
            if (testId == null) {
                ctx.status(500).contentType("text/plain; charset=utf-8")
                        .result("Некорректный testId у попытки: attemptId=" + attemptIdStr);
                return;
            }
            ctx.redirect("/testing/ui/tests/" + testId + "/take?attemptId=" + attemptId);
        } catch (Exception e) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Попытка не найдена: " + attemptIdStr);
        }
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

        List<Question> questions = questionService.getQuestionsByTestId(testId);

        int completedCount = testAttemptService.countCompletedAttempts(studentId, testId);
        int attemptNo = completedCount + 1;

        // init attempt_version, если пустой: сразу кладём текст вопроса + maxPoints + мета теста
        List<TestAttemptService.QuestionInit> init = new ArrayList<>();
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            List<Answer> answers = answerService.getAnswersByQuestionId(qid);
            int maxPoints = 0;
            for (Answer a : answers) maxPoints += Math.max(0, extractAnswerScore(a));

            init.add(new TestAttemptService.QuestionInit(qid, q.getOrder(), q.getTextOfQuestion(), maxPoints));
        }
        testAttemptService.initAttemptVersionIfEmpty(
                attemptId,
                attemptNo,
                init,
                (testDto != null ? testDto.getTitle() : null),
                extractMinPoint(testDto)
        );

        Map<String, List<String>> selectedAnswerByQuestion = loadSelectedAnswers(attemptId);
        Map<String, List<String>> selectedTextByQuestion = loadSelectedAnswerTexts(attemptId);

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
            // Для «сохранить всё» удобнее одинаковый нейм: q_{questionId}
            String inputName = "q_" + qIdStr;
            String hint = multi ? "Можно выбрать несколько вариантов" : "Выберите один вариант";

            // если answerIds не нашли (например, тест меняли), пытаемся подсветить по текстам
            List<String> selectedTextsFallback = selectedTextByQuestion.getOrDefault(qIdStr, Collections.emptyList());

            List<Map<String, Object>> answerViews = new ArrayList<>();
            List<String> selectedTexts = new ArrayList<>();

            for (Answer a : answers) {
                String aIdStr = String.valueOf(a.getId());
                boolean selected = (answered && selectedAnswerIds.contains(aIdStr))
                        || (!answered && selectedTextsFallback != null && selectedTextsFallback.contains(a.getText()));

                Map<String, Object> av = new HashMap<>();
                av.put("id", aIdStr);
                av.put("text", a.getText());
                av.put("selected", selected);
                av.put("disabled", false);

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
     * POST /ui/tests/{testId}/save
     *
     * Сохраняет все текущие ответы из формы. Можно нажимать сколько угодно раз,
     * ответы будут перезаписываться (можно исправлять).
     */
    public void saveAllAnswers(Context ctx) {
        UUID testId = parseUuidOr400(ctx, ctx.pathParam("testId"), "testId");
        if (testId == null) return;

        UUID attemptId = tryParseUuid(ctx.formParam("attempt_id"));
        if (attemptId == null) {
            ctx.status(400).contentType("text/plain; charset=utf-8").result("Нет attempt_id");
            return;
        }

        resolveStudentIdForPost(ctx);

        List<Question> questions = questionService.getQuestionsByTestId(testId);
        persistAllAnswersFromForm(ctx, attemptId, questions);

        ctx.redirect("/testing/ui/tests/" + testId + "/take?attemptId=" + attemptId);
    }

    /**
     * legacy endpoint (не используется в новом UI-шаблоне).
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

        List<Answer> allAnswers = answerService.getAnswersByQuestionId(questionId);
        Map<String, Integer> scoreByAnswerId = new HashMap<>();
        for (Answer a : allAnswers) {
            scoreByAnswerId.put(String.valueOf(a.getId()), Math.max(0, extractAnswerScore(a)));
        }

        List<Integer> answerPoints = new ArrayList<>();
        int earnedPoints = 0;
        for (UUID aid : answerIds) {
            int p = scoreByAnswerId.getOrDefault(aid.toString(), 0);
            answerPoints.add(p);
            earnedPoints += p;
        }

        try {
            testAttemptService.saveAnswers(attemptId, questionId, answerIds, answerPoints, earnedPoints);
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

        // На странице подтверждения тест может уже не существовать (например, админ удалил),
        // поэтому берём данные по возможности из attempt_version.
        Test testDto = null;
        try { testDto = testService.getTestById(testIdStr); } catch (Exception ignored) {}

        UUID studentId = resolveOrCreateStudentId(ctx);

        AttemptVersion av = readAttemptVersion(attemptId);
        List<Map<String, Object>> missing = new ArrayList<>();
        int total = 0;
        int missingCount = 0;
        int answeredCount = 0;

        if (av != null && av.items != null && !av.items.isEmpty()) {
            total = av.items.size();
            for (AttemptItem it : av.items) {
                boolean has = (it.answerIds != null && !it.answerIds.isEmpty()) || (it.answerTexts != null && !it.answerTexts.isEmpty());
                if (has) {
                    answeredCount++;
                } else {
                    missingCount++;
                    Map<String, Object> m = new HashMap<>();
                    m.put("id", (it.questionId == null ? "" : it.questionId));
                    m.put("order", it.order == null ? "" : it.order);
                    missing.add(m);
                }
            }
        } else {
            // fallback (если attempt_version старый/пустой)
            List<Question> questions = questionService.getQuestionsByTestId(testId);
            Map<String, List<String>> selected = loadSelectedAnswers(attemptId);
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
            total = questions.size();
            missingCount = missing.size();
            answeredCount = total - missingCount;
        }

        Map<String, Object> model = new HashMap<>();
        // шаблон ожидает test.title/test.id
        Map<String, Object> testModel = new HashMap<>();
        testModel.put("id", testIdStr);
        testModel.put("title", (testDto != null ? testDto.getTitle() : (av != null && av.testTitle != null ? av.testTitle : "Тест")));
        model.put("test", testModel);
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

        // 1) сначала сохраняем всё, что пришло из формы (можно исправлять)
        List<Question> questions = questionService.getQuestionsByTestId(testId);
        persistAllAnswersFromForm(ctx, attemptId, questions);

        // 2) затем проверяем пропуски уже по attempt_version
        AttemptVersion av = readAttemptVersion(attemptId);
        int missingCount = 0;
        if (av != null && av.items != null) {
            for (AttemptItem it : av.items) {
                boolean has = it.answerIds != null && !it.answerIds.isEmpty();
                boolean hasText = it.answerTexts != null && !it.answerTexts.isEmpty();
                if (!has && !hasText) missingCount++;
            }
        }

        if (missingCount > 0 && !force) {
            ctx.redirect("/testing/ui/tests/" + testId + "/finish?attemptId=" + attemptId);
            return;
        }

        int points = 0;
        if (av != null && av.items != null) {
            for (AttemptItem it : av.items) points += Math.max(0, it.earnedPoints);
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

                // fallback для UI: если answerIds не сохранены/пустые, но есть answerTexts
                JsonNode texts = item.get("answerTexts");
                if ((arr == null || !arr.isArray() || arr.size() == 0) && texts != null && texts.isArray() && texts.size() > 0) {
                    // здесь возвращаем пусто: id-ов всё равно нет, но мы подсветим по текстам в showTakePage
                    map.putIfAbsent(qid, Collections.emptyList());
                }
            }
        } catch (Exception ignored) {}
        return map;
    }

    private Map<String, List<String>> loadSelectedAnswerTexts(UUID attemptId) {
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

                JsonNode arr = item.get("answerTexts");
                if (arr != null && arr.isArray() && arr.size() > 0) {
                    List<String> texts = new ArrayList<>();
                    for (JsonNode n : arr) {
                        String v = n.asText(null);
                        if (v != null && !v.isBlank()) texts.add(v);
                    }
                    if (!texts.isEmpty()) map.put(qid, texts);
                }
            }
        } catch (Exception ignored) {}
        return map;
    }

    // =========================================================
    // attempt_version parsing (для результатов/проверок)
    // =========================================================

    private static class AttemptItem {
        String questionId;
        Integer order;
        String questionText;
        Integer maxPoints;
        List<String> answerIds = new ArrayList<>();
        List<String> answerTexts = new ArrayList<>();
        int earnedPoints;
    }

    private static class AttemptVersion {
        Integer attemptNo;
        String testTitle;
        Integer minPoint;
        List<AttemptItem> items = new ArrayList<>();
    }

    private AttemptVersion readAttemptVersion(UUID attemptId) {
        try {
            Optional<String> jsonOpt = testAttemptService.getAttemptVersionByAttemptId(attemptId);
            if (jsonOpt.isEmpty()) return null;
            String json = jsonOpt.get();
            if (json == null || json.isBlank()) return null;

            JsonNode root = OBJECT_MAPPER.readTree(json);
            AttemptVersion av = new AttemptVersion();
            JsonNode nAttemptNo = root.get("attemptNo");
            if (nAttemptNo != null && nAttemptNo.isInt()) av.attemptNo = nAttemptNo.asInt();
            else if (nAttemptNo != null && nAttemptNo.isTextual()) {
                try { av.attemptNo = Integer.parseInt(nAttemptNo.asText()); } catch (Exception ignored) {}
            }
            av.testTitle = root.path("testTitle").asText(null);
            JsonNode nMin = root.get("minPoint");
            if (nMin != null && nMin.isInt()) av.minPoint = nMin.asInt();

            JsonNode arr = root.path("answers");
            if (arr.isArray()) {
                for (JsonNode item : arr) {
                    if (item == null || !item.isObject()) continue;
                    AttemptItem it = new AttemptItem();
                    it.questionId = item.path("questionId").asText(null);
                    if (item.has("order")) it.order = item.get("order").asInt();
                    it.questionText = item.path("questionText").asText(null);
                    if (item.has("maxPoints")) it.maxPoints = item.get("maxPoints").asInt();
                    it.earnedPoints = item.path("earnedPoints").asInt(0);

                    JsonNode ids = item.get("answerIds");
                    if (ids != null && ids.isArray()) {
                        for (JsonNode v : ids) {
                            String s = v.asText(null);
                            if (s != null && !s.isBlank()) it.answerIds.add(s);
                        }
                    }
                    JsonNode texts = item.get("answerTexts");
                    if (texts != null && texts.isArray()) {
                        for (JsonNode v : texts) {
                            String s = v.asText(null);
                            if (s != null && !s.isBlank()) it.answerTexts.add(s);
                        }
                    }
                    av.items.add(it);
                }
            }
            return av;
        } catch (Exception ignored) {
            return null;
        }
    }

    private void persistAllAnswersFromForm(Context ctx, UUID attemptId, List<Question> questions) {
        for (Question q : questions) {
            UUID qid = safeUuid(q.getId(), "вопроса", ctx);
            if (qid == null) return;

            String key = "q_" + qid;
            List<String> chosenIdsRaw = ctx.formParams(key);
            if (chosenIdsRaw == null) chosenIdsRaw = Collections.emptyList();

            // чистим пустые
            List<String> cleaned = new ArrayList<>();
            for (String s : chosenIdsRaw) {
                if (s != null && !s.isBlank()) cleaned.add(s);
            }

            List<UUID> answerIds = new ArrayList<>();
            for (String s : cleaned) {
                UUID id = tryParseUuid(s);
                if (id != null) answerIds.add(id);
            }

            List<Answer> allAnswers = answerService.getAnswersByQuestionId(qid);
            Map<String, Answer> byId = new HashMap<>();
            int maxPoints = 0;
            for (Answer a : allAnswers) {
                byId.put(String.valueOf(a.getId()), a);
                maxPoints += Math.max(0, extractAnswerScore(a));
            }

            List<String> answerTexts = new ArrayList<>();
            List<Integer> answerPoints = new ArrayList<>();
            int earned = 0;
            for (String s : cleaned) {
                Answer a = byId.get(s);
                if (a == null) continue;
                int p = Math.max(0, extractAnswerScore(a));
                earned += p;
                answerTexts.add(a.getText());
                answerPoints.add(p);
            }

            testAttemptService.upsertAnswers(
                    attemptId,
                    qid,
                    q.getTextOfQuestion(),
                    maxPoints,
                    answerIds,
                    answerTexts,
                    answerPoints,
                    earned
            );
        }
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

        // Для результатов НЕ обязателен сам тест: можем отрисовать только по attempt_version
        Test testDto = null;
        try { testDto = testService.getTestById(testIdStr); } catch (Exception ignored) {}

        UUID studentId = resolveOrCreateStudentId(ctx);

        AttemptVersion av = readAttemptVersion(attemptId);
        if (av == null) {
            ctx.status(404).contentType("text/plain; charset=utf-8")
                    .result("Попытка не найдена или attempt_version пустой");
            return;
        }

        int points = 0;
        int maxPossible = 0;
        int answeredCount = 0;
        List<Map<String, Object>> rows = new ArrayList<>();

        for (AttemptItem it : av.items) {
            boolean answered = (it.answerIds != null && !it.answerIds.isEmpty()) || (it.answerTexts != null && !it.answerTexts.isEmpty());
            if (answered) answeredCount++;

            points += Math.max(0, it.earnedPoints);
            maxPossible += Math.max(0, it.maxPoints == null ? 0 : it.maxPoints);

            Map<String, Object> row = new HashMap<>();
            row.put("order", it.order == null ? "" : it.order);
            row.put("text", it.questionText == null ? "" : it.questionText);
            row.put("answered", answered);
            row.put("selectedText", (it.answerTexts == null || it.answerTexts.isEmpty()) ? null : String.join(", ", it.answerTexts));
            row.put("earned", Math.max(0, it.earnedPoints));
            row.put("maxPoints", Math.max(0, it.maxPoints == null ? 0 : it.maxPoints));
            rows.add(row);
        }

        Integer minPoint = (testDto != null ? extractMinPoint(testDto) : null);
        if (minPoint == null) minPoint = av.minPoint;
        boolean passed = (minPoint != null) && points >= minPoint;

        int attemptNo = (av.attemptNo == null ? readAttemptNoFromJson(attemptId) : av.attemptNo);

        // test в model нужен для шаблона (title/id). Если тест уже удалили, собираем «виртуальный».
        Map<String, Object> testModel = new HashMap<>();
        testModel.put("id", testIdStr);
        String title = (testDto != null ? testDto.getTitle() : null);
        if (title == null || title.isBlank()) title = (av.testTitle != null ? av.testTitle : "Тест");
        testModel.put("title", title);

        Map<String, Object> model = new HashMap<>();
        model.put("test", testModel);
        model.put("student_id", studentId.toString());
        model.put("attempt_id", attemptId.toString());
        model.put("attempt_no", attemptNo);

        model.put("points", points);
        model.put("maxPossible", maxPossible);

        model.put("totalQuestions", av.items.size());
        model.put("answeredCount", answeredCount);

        model.put("minPoint", minPoint);
        model.put("passed", passed);

        model.put("rows", rows);

        String html = HbsRenderer.render("test-results", model);
        ctx.contentType("text/html; charset=utf-8");
        ctx.result(html);
    }

}