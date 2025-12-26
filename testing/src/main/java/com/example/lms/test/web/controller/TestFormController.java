package com.example.lms.test.web.controller;

import java.io.StringWriter;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import com.example.lms.answer.api.dto.Answer;
import com.example.lms.answer.domain.service.AnswerService;
import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.service.DraftService;
import com.example.lms.question.api.dto.Question;
import com.example.lms.question.domain.service.QuestionService;
import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test.web.dto.QuestionWithAnswers;
import com.example.lms.test.web.dto.TestFormData;
import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.Template;

import io.javalin.http.Context;
import io.javalin.http.NotFoundResponse;

/**
 * Контроллер для веб-интерфейса тестов (HTML формы)
 */
public class TestFormController {
    
    private final TestService testService;
    private final QuestionService questionService;
    private final AnswerService answerService;
    private final DraftService draftService;
    private final Handlebars handlebars;
    
    public TestFormController(TestService testService,
                             QuestionService questionService,
                             AnswerService answerService,
                             DraftService draftService,
                             Handlebars handlebars) {
        this.testService = testService;
        this.questionService = questionService;
        this.answerService = answerService;
        this.draftService = draftService;
        this.handlebars = handlebars;
    }
    
    /**
     * GET / - Главная страница (перенаправление на создание теста)
     */
    public void showHomePage(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        String redirectParam = redirect != null ? "?redirect=" + redirect : "";
        ctx.redirect("/web/tests/new" + redirectParam);
    }
    
    /**
     * GET /web/tests/new - Форма создания нового теста
     */
    public void showNewTestForm(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        Map<String, Object> model = new HashMap<>();
        model.put("title", "Создание нового теста");
        String courseId = ctx.queryParam("course_id");
        model.put("courseId", courseId);
        model.put("test", new Test(null, courseId, "", 60, ""));
        model.put("questions", new ArrayList<>());
        model.put("page", "test-builder");
        model.put("isNew", true);
        model.put("isDraft", false);
        model.put("success", ctx.queryParam("success") != null);
        model.put("redirect", redirect);
        
        // Добавляем все query параметры в скрытые поля формы
        addAllQueryParamsToModel(ctx, model);
        
        renderTemplateWithBody(ctx, "test-builder", model);
    }
    
    /**
     * GET /web/tests/{id}/edit - Универсальный редактор:
     * - Если id начинается с "draft-" - редактирование черновика
     * - Иначе - редактирование теста
     */
    public void showEditTestForm(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String id = ctx.pathParam("id");
            
            if (id.startsWith("draft-")) {
                // Редактирование черновика
                String draftId = id.substring(6); // Убираем "draft-" префикс
                UUID draftUuid = UUID.fromString(draftId);
                
                showEditDraftForm(ctx, draftUuid, redirect);
            } else {
                // Редактирование теста
                UUID testUuid = UUID.fromString(id);
                showEditPublishedTestForm(ctx, testUuid, redirect);
            }
            
        } catch (Exception e) {
            e.printStackTrace();
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка загрузки");
            model.put("error", "Тест или черновик не найден. Создайте новый тест.");
            model.put("test", new Test(null, null, "", 60, ""));
            model.put("questions", new ArrayList<>());
            model.put("page", "test-editor");
            model.put("isNew", true);
            model.put("isDraft", false);
            model.put("redirect", redirect);
            
            // Добавляем все query параметры в скрытые поля формы
            addAllQueryParamsToModel(ctx, model);
            
            renderTemplateWithBody(ctx, "test-editor", model);
        }
    }
    
    /**
     * Редактирование опубликованного теста
     */
    private void showEditPublishedTestForm(Context ctx, UUID testUuid, String redirect) throws Exception {
        String testId = testUuid.toString();
        
        // Получаем тест
        Test test = testService.getTestById(testId);
        if (test == null) {
            throw new NotFoundResponse("Тест не найден");
        }
        
        // Получаем вопросы для этого теста
        List<Question> questions = questionService.getQuestionsByTestId(testUuid);
        
        // Для каждого вопроса получаем ответы
        List<QuestionWithAnswers> questionsWithAnswers = questions.stream()
            .map(question -> {
                List<Answer> answers = answerService.getAnswersByQuestionId(question.getId());
                return new QuestionWithAnswers(question, answers);
            })
            .collect(Collectors.toList());
        
        Map<String, Object> model = new HashMap<>();

        // Проверяем, есть ли черновик для этого теста
        boolean hasDraft = false;
        Draft draft = null;
        try {
            draft = draftService.getDraftByTestId(testUuid);
            if (draft != null) {
                hasDraft = true;
                model.put("hasDraft", true);
                model.put("draft", draft);
                model.put("draftId", "draft-" + draft.getId()); // С префиксом
            }
        } catch (Exception e) {
            System.out.println("DEBUG: No draft found for test: " + e.getMessage());
        }
        
        model.put("title", "Редактирование теста: " + test.getTitle());
        model.put("test", test);
        model.put("courseId", test.getCourseId());
        model.put("questions", questionsWithAnswers);
        model.put("page", "test-editor");
        model.put("isNew", false);
        model.put("isDraft", false); // Это опубликованный тест
        model.put("entityId", testId);
        model.put("testId", testId);
        model.put("hasDraft", hasDraft);
        model.put("success", ctx.queryParam("success") != null);
        model.put("redirect", redirect);
        
        // Добавляем все query параметры в скрытые поля формы
        addAllQueryParamsToModel(ctx, model);
        
        renderTemplateWithBody(ctx, "test-editor", model);
    }
    
    /**
     * Редактирование черновика
     */
    private void showEditDraftForm(Context ctx, UUID draftUuid, String redirect) throws Exception {
        String draftId = draftUuid.toString();
        
        // Получаем черновик
        Draft draft = draftService.getDraftById(draftUuid);
        if (draft == null) {
            throw new NotFoundResponse("Черновик не найден");
        }
        
        // Получаем вопросы черновика
        List<Question> questions = questionService.getQuestionsByDraftId(draftUuid);
        
        // Для каждого вопроса получаем ответы
        List<QuestionWithAnswers> questionsWithAnswers = questions.stream()
            .map(question -> {
                List<Answer> answers = answerService.getAnswersByQuestionId(question.getId());
                return new QuestionWithAnswers(question, answers);
            })
            .collect(Collectors.toList());
        
        Map<String, Object> model = new HashMap<>();
        
        // Создаем временный объект теста для отображения
        Test testForDisplay = new Test(
            draft.getTestId() != null ? draft.getTestId().toString() : null,
            draft.getCourseId() != null ? draft.getCourseId().toString() : null,
            draft.getTitle(),
            draft.getMin_point(),
            draft.getDescription()
        );
        
        String displayId = "draft-" + draftId;
        
        model.put("title", "Редактирование черновика: " + draft.getTitle());
        model.put("test", testForDisplay);
        model.put("courseId", testForDisplay.getCourseId());
        model.put("questions", questionsWithAnswers);
        model.put("page", "test-editor");
        model.put("isNew", false);
        model.put("isDraft", true); // Это черновик
        model.put("entityId", displayId);
        model.put("draftId", draftId);
        model.put("testId", displayId); // Для совместимости
        model.put("success", ctx.queryParam("success") != null);
        model.put("redirect", redirect);
        
        // Если черновик привязан к тесту, показываем ссылку на него
        if (draft.getTestId() != null) {
            model.put("hasPublishedTest", true);
            model.put("publishedTestId", draft.getTestId().toString());
        }
        
        // Добавляем все query параметры в скрытые поля формы
        addAllQueryParamsToModel(ctx, model);
        
        renderTemplateWithBody(ctx, "test-editor", model);
    }
    
    /**
     * POST /web/tests/save - Универсальное сохранение:
     * - Если это черновик (isDraft=true) - сохраняем как черновик
     * - Если это тест - сохраняем как тест
     */
    public void saveTestFromForm(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String entityId = ctx.formParam("entityId");
            Boolean isDraft = Boolean.parseBoolean(ctx.formParam("isDraft"));
            Boolean isNew = Boolean.parseBoolean(ctx.formParam("isNew"));
            
            // Валидация формы перед сохранением
            String validationError = validateForm(ctx);
            if (validationError != null) {
                // Если валидация не прошла, показываем форму снова с ошибками
                Map<String, Object> model = new HashMap<>();
                model.put("title", "Ошибка валидации");
                model.put("error", validationError);
                model.put("test", getTestFromForm(ctx));
                model.put("questions", getQuestionsFromForm(ctx));
                model.put("page", "test-editor");
                model.put("isNew", isNew);
                model.put("isDraft", isDraft);
                model.put("entityId", entityId);
                model.put("redirect", redirect);
                
                // Добавляем все query параметры в скрытые поля формы
                addAllQueryParamsToModel(ctx, model);
                
                renderTemplateWithBody(ctx, "test-editor", model);
                return;
            }
            
            if (isDraft) {
                // Сохранение черновика
                saveDraftFromForm(ctx, entityId, isNew);
            } else {
                // Сохранение теста
                savePublishedTestFromForm(ctx, entityId, isNew);
            }
            
        } catch (Exception e) {
            e.printStackTrace();
            // При ошибке показываем форму снова
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка сохранения");
            model.put("error", "Ошибка при сохранении: " + e.getMessage());
            model.put("test", getTestFromForm(ctx));
            model.put("questions", getQuestionsFromForm(ctx));
            model.put("page", "test-editor");
            model.put("isNew", Boolean.parseBoolean(ctx.formParam("isNew")));
            model.put("isDraft", Boolean.parseBoolean(ctx.formParam("isDraft")));
            model.put("redirect", redirect);
            
            // Добавляем все query параметры в скрытые поля формы
            addAllQueryParamsToModel(ctx, model);
            
            renderTemplateWithBody(ctx, "test-editor", model);
        }
    }
    
    /**
     * Валидация данных формы на сервере (полная проверка, как на клиенте)
     */
    private String validateForm(Context ctx) {
        // 1. Проверка названия теста
        String title = ctx.formParam("title");
        if (title == null || title.trim().isEmpty()) {
            return "Название теста обязательно для заполнения";
        }
        
        // 2. Проверка минимального балла
        String minPointStr = ctx.formParam("min_point");
        if (minPointStr == null || minPointStr.trim().isEmpty()) {
            return "Минимальный балл обязательно для заполнения";
        }
        
        int minPoint;
        try {
            minPoint = Integer.parseInt(minPointStr);
        } catch (NumberFormatException e) {
            return "Минимальный балл должен быть числом";
        }
        
        if (minPoint < 0) {
            return "Минимальный балл не может быть отрицательным";
        }
        
        // 3. Проверка вопросов
        TestFormData formData = parseTestFormData(ctx);
        if (formData.getQuestions() == null || formData.getQuestions().isEmpty()) {
            return "Тест должен содержать хотя бы один вопрос";
        }
        
        // 4. Проверка каждого вопроса
        int questionIndex = 0;
        for (TestFormData.QuestionFormData question : formData.getQuestions()) {
            questionIndex++;
            
            // Проверка текста вопроса
            if (question.getTextOfQuestion() == null || 
                question.getTextOfQuestion().trim().isEmpty()) {
                return "Вопрос " + questionIndex + " не может быть пустым";
            }
            
            // Проверка ответов
            if (question.getAnswers() == null || question.getAnswers().size() < 2) {
                return "Вопрос " + questionIndex + " должен содержать минимум 2 ответа";
            }
            
            // Проверка, что есть хотя бы один ответ с баллами > 0
            boolean hasPositiveScore = false;
            boolean hasEmptyAnswerText = false;
            int emptyAnswerCount = 0;
            
            for (TestFormData.AnswerFormData answer : question.getAnswers()) {
                if (answer.getText() == null || answer.getText().trim().isEmpty()) {
                    hasEmptyAnswerText = true;
                    emptyAnswerCount++;
                }
                
                if (answer.getScore() != null && answer.getScore() > 0) {
                    hasPositiveScore = true;
                }
            }
            
            // Если есть пустые ответы
            if (hasEmptyAnswerText) {
                return "Вопрос " + questionIndex + ": некоторые ответы пустые";
            }
            
            // Если нет правильного ответа
            if (!hasPositiveScore) {
                return "Вопрос " + questionIndex + ": должен быть хотя бы один правильный ответ (баллы > 0)";
            }
        }
        
        // 5. Проверка, что минимальный балл не превышает максимальный балл за тест
        // Рассчитываем максимальный балл за тест
        int maxTotalPoints = 0;
        for (TestFormData.QuestionFormData question : formData.getQuestions()) {
            if (question.getAnswers() != null) {
                for (TestFormData.AnswerFormData answer : question.getAnswers()) {
                    Integer score = answer.getScore() != null ? answer.getScore() : 0;
                    if (score > 0) {
                        maxTotalPoints += score;
                    }
                }
            }
        }
        
        if (minPoint > maxTotalPoints) {
            return "Минимальный порог (" + minPoint + ") не может быть больше максимального количества баллов за тест (" + maxTotalPoints + ")";
        }
        
        return null; // Валидация пройдена
    }
    
    /**
     * Сохранение опубликованного теста
     */
    private void savePublishedTestFromForm(Context ctx, String entityId, boolean isNew) throws Exception {
        String redirect = getRedirectUrl(ctx);
        
        // 1. Получаем и валидируем данные теста
        Test test = getTestFromForm(ctx);
        if (test.getTitle() == null || test.getTitle().trim().isEmpty()) {
            throw new IllegalArgumentException("Название теста обязательно");
        }
        
        // 2. Получаем вопросы и ответы из формы
        TestFormData formData = parseTestFormData(ctx);
        
        // 3. Сохраняем тест в базу
        Test savedTest;
        if (isNew || entityId == null || entityId.isEmpty() || entityId.equals("new")) {
            savedTest = testService.createTest(test);
        } else {
            savedTest = testService.updateTest(test);
        }
        
        // 4. Сохраняем вопросы и ответы
        UUID testUuid = UUID.fromString(savedTest.getId());
        saveQuestionsAndAnswers(testUuid, formData);
        
        // 5. Делаем редирект на указанный URL с сохранением всех параметров
        performRedirect(ctx, redirect, "/web/tests/" + savedTest.getId() + "/edit?success=true");
    }
    
    /**
     * Сохранение черновика
     */
    private void saveDraftFromForm(Context ctx, String entityId, boolean isNew) throws Exception {
        String redirect = getRedirectUrl(ctx);
        
        // Получаем данные из формы
        String title = ctx.formParam("title");
        String description = ctx.formParam("description");
        
        int minPoint;
        try {
            minPoint = Integer.parseInt(ctx.formParam("min_point"));
        } catch (NumberFormatException e) {
            minPoint = 60;
        }
        
        // Парсим вопросы из формы
        TestFormData formData = parseTestFormData(ctx);
        
        if (isNew || entityId == null || entityId.isEmpty() || entityId.equals("new")) {
            // Новый черновик (с формы создания нового теста)
            saveNewTestDraft(ctx);
        } else {
            // Редактирование существующего черновика
            if (entityId.startsWith("draft-")) {
                String draftId = entityId.substring(6);
                UUID draftUuid = UUID.fromString(draftId);
                
                // Получаем существующий черновик
                Draft existingDraft = draftService.getDraftById(draftUuid);
                if (existingDraft == null) {
                    throw new NotFoundResponse("Черновик не найден");
                }
                
                // Обновляем черновик
                Draft draft = new Draft();
                draft.setId(draftUuid);
                draft.setTitle(title);
                draft.setMin_point(minPoint);
                draft.setDescription(description);
                draft.setTestId(existingDraft.getTestId()); // Сохраняем связь с тестом
                draft.setCourseId(existingDraft.getCourseId()); // Сохраняем курс
                
                Draft savedDraft = draftService.updateDraft(draft);
                
                // Удаляем старые вопросы черновика
                deleteDraftQuestions(draftUuid);
                
                // Сохраняем новые вопросы черновика
                if (formData.getQuestions() != null && !formData.getQuestions().isEmpty()) {
                    saveDraftQuestionsAndAnswers(savedDraft.getId(), existingDraft.getTestId(), formData);
                }
                
                // Делаем редирект на указанный URL с сохранением всех параметров
                performRedirect(ctx, redirect, "/web/tests/draft-" + savedDraft.getId() + "/edit?success=true");
            }
        }
    }
    
    /**
     * POST /web/tests/draft - Сохранение черновика нового теста
     */
    public void saveNewTestDraft(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        saveDraftInternal(ctx, null, true, redirect);
    }

    /**
     * POST /web/tests/{id}/draft - Сохранение черновика существующего теста
     */
    public void saveExistingTestDraft(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        String testId = ctx.pathParam("id");
        UUID testUuid = null;
        if (testId != null && !testId.trim().isEmpty()) {
            try {
                testUuid = UUID.fromString(testId);
            } catch (IllegalArgumentException e) {
                // Невалидный UUID, сохраняем как черновик без test_id
            }
        }
        saveDraftInternal(ctx, testUuid, false, redirect);
    }
    
    /**
     * Внутренний метод для сохранения черновика
     */
    private void saveDraftInternal(Context ctx, UUID testId, boolean isNewDraft, String redirect) {
        try {
            System.out.println("DEBUG: saveDraftInternal called with testId = '" + testId + "', isNewDraft = " + isNewDraft + ", redirect = " + redirect);
            
            // Получаем данные из формы
            String title = ctx.formParam("title");
            String description = ctx.formParam("description");
            
            int minPoint;
            try {
                minPoint = Integer.parseInt(ctx.formParam("min_point"));
            } catch (NumberFormatException e) {
                minPoint = 60;
            }
            
            System.out.println("DEBUG: title = '" + title + "', minPoint = " + minPoint);
            
            // Парсим вопросы из формы
            TestFormData formData = parseTestFormData(ctx);
            System.out.println("DEBUG: Parsed form data, questions count: " + 
                              (formData.getQuestions() != null ? formData.getQuestions().size() : 0));
            
            // 1. Сохраняем черновик в таблицу draft_b
            Draft draft = new Draft();
            draft.setTitle(title);
            draft.setMin_point(minPoint);
            draft.setDescription(description);
            draft.setTestId(testId); // ← UUID или null
            
            // Проверяем, существует ли уже черновик для этого testId
            Draft existingDraft = null;
            if (testId != null) {
                try {
                    existingDraft = draftService.getDraftByTestId(testId);
                    System.out.println("DEBUG: Existing draft found for testId " + testId + ": " + (existingDraft != null));
                } catch (Exception e) {
                    System.out.println("DEBUG: No existing draft found for testId " + testId);
                }
            }
            
            Draft savedDraft;
            if (existingDraft != null) {
                // Обновляем существующий черновик
                draft.setId(existingDraft.getId());
                savedDraft = draftService.updateDraft(draft);
                System.out.println("DEBUG: Draft updated, id = " + savedDraft.getId());
                
                // Удаляем старые вопросы этого черновика
                deleteDraftQuestions(savedDraft.getId());
            } else {
                // Создаем новый черновик
                savedDraft = draftService.createDraft(draft);
                System.out.println("DEBUG: Draft created with id = " + savedDraft.getId());
            }
            
            // 2. ВАЖНО: Сохраняем вопросы и ответы в таблицу question_d с draft_id
            if (formData.getQuestions() != null && !formData.getQuestions().isEmpty()) {
                System.out.println("DEBUG: Saving questions to question_d with draft_id = " + savedDraft.getId());
                saveDraftQuestionsAndAnswers(savedDraft.getId(), testId, formData);
                System.out.println("DEBUG: Questions saved as draft");
            }
            
            // 3. Делаем редирект на указанный URL с сохранением всех параметров
            performRedirect(ctx, redirect, "/web/tests/draft-" + savedDraft.getId() + "/edit?success=true");
            
        } catch (Exception e) {
            System.out.println("ERROR in saveDraftInternal: " + e.getMessage());
            e.printStackTrace();
            
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка сохранения черновика");
            model.put("error", "Ошибка при сохранении черновика: " + e.getMessage());
            
            try {
                Test testFromForm = getTestFromForm(ctx);
                model.put("test", testFromForm);
                model.put("questions", getQuestionsFromForm(ctx));
            } catch (Exception ex) {
                model.put("test", new Test(null, null, "", 60, ""));
                model.put("questions", new ArrayList<>());
            }
            
            model.put("page", "test-editor");
            model.put("isNew", true);
            model.put("redirect", redirect);
            
            // Добавляем все query параметры в скрытые поля формы
            addAllQueryParamsToModel(ctx, model);
            
            renderTemplateWithBody(ctx, "test-editor", model);
        }
    }

    /**
     * Сохраняет вопросы и ответы для черновика
     */
    private void saveDraftQuestionsAndAnswers(UUID draftId, UUID testId, TestFormData formData) throws Exception {
        // СОХРАНЯЕМ ВОПРОСЫ С DRAFT_ID
        for (int i = 0; i < formData.getQuestions().size(); i++) {
            TestFormData.QuestionFormData questionData = formData.getQuestions().get(i);
            
            if (questionData.getTextOfQuestion() != null && 
                !questionData.getTextOfQuestion().trim().isEmpty()) {
                
                // Создаем вопрос ДЛЯ ЧЕРНОВИКА
                Question question = new Question();
                question.setDraftId(draftId); // ← UUID!
                
                // Если редактируем существующий тест, можно установить связь
                if (testId != null) {
                    question.setTestId(testId);
                }
                
                question.setTextOfQuestion(questionData.getTextOfQuestion());
                question.setOrder(questionData.getOrder() != null ? questionData.getOrder() : i);
                
                Question savedQuestion = questionService.createQuestion(question);
                
                // Сохраняем ответы
                if (questionData.getAnswers() != null) {
                    for (TestFormData.AnswerFormData answerData : questionData.getAnswers()) {
                        if (answerData.getText() != null && !answerData.getText().trim().isEmpty()) {
                            Answer answer = new Answer();
                            answer.setQuestionId(savedQuestion.getId());
                            answer.setText(answerData.getText());
                            
                            Integer score = answerData.getScore() != null ? answerData.getScore() : 0;
                            answer.setScore(score);
                            
                            answerService.createAnswer(answer);
                        }
                    }
                }
            }
        }
    }

    /**
     * Удаляет вопросы черновика по draft_id
     */
    private void deleteDraftQuestions(UUID draftId) {
        try {
            List<Question> draftQuestions = questionService.getQuestionsByDraftId(draftId);
            for (Question question : draftQuestions) {
                try {
                    // Удаляем ответы
                    List<Answer> answers = answerService.getAnswersByQuestionId(question.getId());
                    for (Answer answer : answers) {
                        answerService.deleteAnswer(answer.getId());
                    }
                } catch (Exception e) {
                    System.out.println("DEBUG: Error deleting answers: " + e.getMessage());
                }
                
                try {
                    // Удаляем вопрос
                    questionService.deleteQuestion(question.getId());
                } catch (Exception e) {
                    System.out.println("DEBUG: Error deleting question: " + e.getMessage());
                }
            }
            System.out.println("DEBUG: Deleted " + draftQuestions.size() + " draft questions");
        } catch (Exception e) {
            System.out.println("DEBUG: Error deleting draft questions: " + e.getMessage());
        }
    }
    
    /**
     * PUT /web/tests/{id} - Обновление существующего теста
     */
    public void updateTestFromForm(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String testId = ctx.pathParam("id");
            Test test = getTestFromForm(ctx);
            test.setId(testId);
            
            // Валидация
            if (test.getTitle() == null || test.getTitle().trim().isEmpty()) {
                throw new IllegalArgumentException("Название теста обязательно");
            }
            
            TestFormData formData = parseTestFormData(ctx);
            Test updatedTest = testService.updateTest(test);
            
            // Обновляем вопросы и ответы
            UUID testUuid = UUID.fromString(updatedTest.getId());
            saveQuestionsAndAnswers(testUuid, formData);
            
            // Делаем редирект на указанный URL с сохранением всех параметров
            performRedirect(ctx, redirect, "/web/tests/" + updatedTest.getId() + "/edit?success=true");
            
        } catch (Exception e) {
            e.printStackTrace();
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка обновления");
            model.put("error", "Ошибка при обновлении теста: " + e.getMessage());
            model.put("test", getTestFromForm(ctx));
            model.put("questions", getQuestionsFromForm(ctx));
            model.put("page", "test-editor");
            model.put("isNew", false);
            model.put("isDraft", false);
            model.put("redirect", redirect);
            
            // Добавляем все query параметры в скрытые поля формы
            addAllQueryParamsToModel(ctx, model);
            
            renderTemplateWithBody(ctx, "test-editor", model);
        }
    }
    
    /**
     * PUT /web/tests/draft-{id} - Обновление черновика
     */
    public void updateDraftFromForm(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String id = ctx.pathParam("id");
            
            if (!id.startsWith("draft-")) {
                throw new IllegalArgumentException("ID должен быть черновиком");
            }
            
            String draftId = id.substring(6);
            UUID draftUuid = UUID.fromString(draftId);
            
            // Получаем данные из формы
            String title = ctx.formParam("title");
            String description = ctx.formParam("description");
            
            int minPoint;
            try {
                minPoint = Integer.parseInt(ctx.formParam("min_point"));
            } catch (NumberFormatException e) {
                minPoint = 60;
            }
            
            // Получаем существующий черновик
            Draft existingDraft = draftService.getDraftById(draftUuid);
            if (existingDraft == null) {
                throw new NotFoundResponse("Черновик не найден");
            }
            
            // Обновляем черновик
            Draft draft = new Draft();
            draft.setId(draftUuid);
            draft.setTitle(title);
            draft.setMin_point(minPoint);
            draft.setDescription(description);
            draft.setTestId(existingDraft.getTestId());
            draft.setCourseId(existingDraft.getCourseId());
            
            Draft updatedDraft = draftService.updateDraft(draft);
            
            // Обновляем вопросы и ответы
            TestFormData formData = parseTestFormData(ctx);
            deleteDraftQuestions(draftUuid);
            
            if (formData.getQuestions() != null && !formData.getQuestions().isEmpty()) {
                saveDraftQuestionsAndAnswers(updatedDraft.getId(), existingDraft.getTestId(), formData);
            }
            
            // Делаем редирект на указанный URL с сохранением всех параметров
            performRedirect(ctx, redirect, "/web/tests/" + id + "/edit?success=true");
            
        } catch (Exception e) {
            e.printStackTrace();
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка обновления черновика");
            model.put("error", "Ошибка при обновлении черновика: " + e.getMessage());
            model.put("test", getTestFromForm(ctx));
            model.put("questions", getQuestionsFromForm(ctx));
            model.put("page", "test-editor");
            model.put("isNew", false);
            model.put("isDraft", true);
            model.put("redirect", redirect);
            
            // Добавляем все query параметры в скрытые поля формы
            addAllQueryParamsToModel(ctx, model);
            
            renderTemplateWithBody(ctx, "test-editor", model);
        }
    }
    
    /**
     * DELETE /web/tests/{id} - Удаление теста
     */
    public void deleteTest(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String testId = ctx.pathParam("id");
            testService.deleteTest(testId);
            
            // Делаем редирект на указанный URL с сохранением всех параметров
            performRedirect(ctx, redirect, "/web/tests/new?deleted=true");
        } catch (Exception e) {
            e.printStackTrace();
            // При ошибке делаем редирект с параметром ошибки
            performRedirect(ctx, redirect, "/web/tests/" + ctx.pathParam("id") + "/edit?error=delete_failed");
        }
    }
    
    /**
     * DELETE /web/tests/draft-{id} - Удаление черновика
     */
    public void deleteDraft(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String id = ctx.pathParam("id");
            
            if (!id.startsWith("draft-")) {
                throw new IllegalArgumentException("ID должен быть черновиком");
            }
            
            String draftId = id.substring(6);
            UUID draftUuid = UUID.fromString(draftId);
            
            // Получаем черновик для проверки связей
            Draft draft = draftService.getDraftById(draftUuid);
            
            // Удаляем черновик
            draftService.deleteDraft(draftUuid);
            
            // Делаем редирект на указанный URL с сохранением всех параметров
            if (draft != null && draft.getTestId() != null) {
                performRedirect(ctx, redirect, "/web/tests/" + draft.getTestId() + "/edit?draft_deleted=true");
            } else {
                performRedirect(ctx, redirect, "/web/tests/new?draft_deleted=true");
            }
            
        } catch (Exception e) {
            e.printStackTrace();
            // При ошибке делаем редирект с параметром ошибки
            performRedirect(ctx, redirect, "/web/tests/" + ctx.pathParam("id") + "/edit?error=draft_delete_failed");
        }
    }
    
    /**
     * POST /web/tests/{id}/publish - Публикация черновика в тест
     */
    public void publishDraft(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String id = ctx.pathParam("id");
            
            // Проверяем, является ли это черновиком
            if (!id.startsWith("draft-")) {
                throw new IllegalArgumentException("ID должен быть черновиком (начинаться с 'draft-')");
            }
            
            String draftId = id.substring(6);
            UUID draftUuid = UUID.fromString(draftId);
            
            // Получаем черновик
            Draft draft = draftService.getDraftById(draftUuid);
            if (draft == null) {
                throw new NotFoundResponse("Черновик не найден");
            }
            
            Test test;
            Test savedTest;
            
            if (draft.getTestId() != null) {
                // Черновик привязан к существующему тесту - обновляем его
                test = new Test(
                    draft.getTestId().toString(),
                    draft.getCourseId() != null ? draft.getCourseId().toString() : null,
                    draft.getTitle(),
                    draft.getMin_point(),
                    draft.getDescription()
                );
                savedTest = testService.updateTest(test);
            } else {
                // Создаем новый тест из черновика
                test = new Test(
                    null,
                    draft.getCourseId() != null ? draft.getCourseId().toString() : null,
                    draft.getTitle(),
                    draft.getMin_point(),
                    draft.getDescription()
                );
                savedTest = testService.createTest(test);
            }
            
            UUID testUuid = UUID.fromString(savedTest.getId());
            
            // Удаляем старые вопросы теста (если обновляем существующий)
            if (draft.getTestId() != null) {
                deleteTestQuestions(testUuid);
            }
            
            // Копируем вопросы из черновика в тест
            List<Question> draftQuestions = questionService.getQuestionsByDraftId(draftUuid);
            for (Question draftQuestion : draftQuestions) {
                Question question = new Question();
                question.setTestId(testUuid);
                question.setTextOfQuestion(draftQuestion.getTextOfQuestion());
                question.setOrder(draftQuestion.getOrder());
                
                Question savedQuestion = questionService.createQuestion(question);
                
                // Копируем ответы
                List<Answer> draftAnswers = answerService.getAnswersByQuestionId(draftQuestion.getId());
                for (Answer draftAnswer : draftAnswers) {
                    Answer answer = new Answer();
                    answer.setQuestionId(savedQuestion.getId());
                    answer.setText(draftAnswer.getText());
                    answer.setScore(draftAnswer.getScore());
                    
                    answerService.createAnswer(answer);
                }
            }
            
            // Удаляем черновик после публикации
            draftService.deleteDraft(draftUuid);
            
            // Делаем редирект на указанный URL с сохранением всех параметров
            performRedirect(ctx, redirect, "/web/tests/" + savedTest.getId() + "/edit?success=true");
            
        } catch (Exception e) {
            System.out.println("ERROR in publishDraft: " + e.getMessage());
            e.printStackTrace();
            
            // При ошибке делаем редирект с параметром ошибки
            performRedirect(ctx, redirect, "/web/tests/" + ctx.pathParam("id") + "/edit?error=true");
        }
    }
    
    /**
     * POST /web/tests/{id}/create-draft - Создание черновика из теста
     */
    public void createDraftFromTest(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String testId = ctx.pathParam("id");
            UUID testUuid = UUID.fromString(testId);
            
            // Получаем тест
            Test test = testService.getTestById(testId);
            if (test == null) {
                throw new NotFoundResponse("Тест не найден");
            }
            
            // Проверяем, нет ли уже черновика для этого теста
            Draft existingDraft = draftService.getDraftByTestId(testUuid);
            if (existingDraft != null) {
                // Если черновик уже есть, делаем редирект на его редактирование
                performRedirect(ctx, redirect, "/web/tests/draft-" + existingDraft.getId() + "/edit");
                return;
            }
            
            // Создаем черновик на основе теста
            Draft draft = new Draft();
            draft.setTitle(test.getTitle());
            draft.setMin_point(test.getMin_point());
            draft.setDescription(test.getDescription());
            draft.setTestId(testUuid);
            draft.setCourseId(test.getCourseId() != null ? UUID.fromString(test.getCourseId()) : null);
            
            Draft savedDraft = draftService.createDraft(draft);
            
            // Копируем вопросы из теста в черновик
            List<Question> testQuestions = questionService.getQuestionsByTestId(testUuid);
            for (Question testQuestion : testQuestions) {
                Question draftQuestion = new Question();
                draftQuestion.setDraftId(savedDraft.getId());
                draftQuestion.setTextOfQuestion(testQuestion.getTextOfQuestion());
                draftQuestion.setOrder(testQuestion.getOrder());
                
                Question savedDraftQuestion = questionService.createQuestion(draftQuestion);
                
                // Копируем ответы
                List<Answer> testAnswers = answerService.getAnswersByQuestionId(testQuestion.getId());
                for (Answer testAnswer : testAnswers) {
                    Answer draftAnswer = new Answer();
                    draftAnswer.setQuestionId(savedDraftQuestion.getId());
                    draftAnswer.setText(testAnswer.getText());
                    draftAnswer.setScore(testAnswer.getScore());
                    
                    answerService.createAnswer(draftAnswer);
                }
            }
            
            // Делаем редирект на указанный URL с сохранением всех параметров
            performRedirect(ctx, redirect, "/web/tests/draft-" + savedDraft.getId() + "/edit");
            
        } catch (Exception e) {
            System.out.println("ERROR in createDraftFromTest: " + e.getMessage());
            e.printStackTrace();
            
            // При ошибке делаем редирект с параметром ошибки
            performRedirect(ctx, redirect, "/web/tests/" + ctx.pathParam("id") + "/edit?error=true");
        }
    }
    
    /**
     * GET /web/tests/{id}/preview - Предпросмотр теста
     */
    public void previewTest(Context ctx) {
        String redirect = getRedirectUrl(ctx);
        
        try {
            String testId = ctx.pathParam("id");
            Test test = testService.getTestById(testId);
            
            if (test == null) {
                throw new NotFoundResponse("Тест не найден");
            }
            
            // Получаем вопросы для этого теста
            List<Question> questions = questionService.getQuestionsByTestId(UUID.fromString(testId));
            
            // Для каждого вопроса получаем ответы
            List<QuestionWithAnswers> questionsWithAnswers = questions.stream()
                .map(question -> {
                    List<Answer> answers = answerService.getAnswersByQuestionId(question.getId());
                    return new QuestionWithAnswers(question, answers);
                })
                .collect(Collectors.toList());
            
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Предпросмотр: " + test.getTitle());
            model.put("test", test);
            model.put("questions", questionsWithAnswers);
            model.put("page", "test-preview");
            model.put("redirect", redirect);
            
            // Добавляем все query параметры в скрытые поля формы
            addAllQueryParamsToModel(ctx, model);
            
            renderTemplateWithBody(ctx, "test-preview", model);
        } catch (Exception e) {
            e.printStackTrace();
            // При ошибке делаем редирект
            performRedirect(ctx, redirect, "/web/tests/new");
        }
    }
    
    /**
     * PATCH /web/tests/{id} - Частичное обновление теста
     */
    public void partialUpdateTest(Context ctx) {
        try {
            String testId = ctx.pathParam("id");
            // TODO: Реализовать частичное обновление теста
            ctx.status(200).result("Частичное обновление теста ID: " + testId + " - пока не реализовано");
        } catch (Exception e) {
            e.printStackTrace();
            ctx.status(500).result("Ошибка при частичном обновлении теста: " + e.getMessage());
        }
    }
    
    /**
     * PATCH /web/tests/draft-{id} - Частичное обновление черновика
     */
    public void partialUpdateDraft(Context ctx) {
        try {
            String id = ctx.pathParam("id");
            if (!id.startsWith("draft-")) {
                throw new IllegalArgumentException("ID должен быть черновиком");
            }
            String draftId = id.substring(6);
            // TODO: Реализовать частичное обновление черновика
            ctx.status(200).result("Частичное обновление черновика ID: " + draftId + " - пока не реализовано");
        } catch (Exception e) {
            e.printStackTrace();
            ctx.status(500).result("Ошибка при частичном обновлении черновика: " + e.getMessage());
        }
    }
    
    /**
     * Вспомогательный метод для рендеринга с main.hbs
     */
    private void renderTemplateWithBody(Context ctx, String contentTemplate, Map<String, Object> model) {
        try {
            // Сначала рендерим контентный шаблон
            Template contentTemplateObj = handlebars.compile(contentTemplate);
            StringWriter contentWriter = new StringWriter();
            contentTemplateObj.apply(model, contentWriter);
            
            // Добавляем сгенерированный контент в модель
            model.put("body", contentWriter.toString());
            
            // Теперь рендерим main.hbs с контентом в body
            Template mainTemplate = handlebars.compile("main");
            StringWriter mainWriter = new StringWriter();
            mainTemplate.apply(model, mainWriter);
            
            ctx.contentType("text/html; charset=utf-8");
            ctx.result(mainWriter.toString());
            
        } catch (Exception e) {
            e.printStackTrace();
            String errorHtml = "<html><body><h1>Ошибка</h1><p>" + e.getMessage() + "</p></body></html>";
            ctx.contentType("text/html; charset=utf-8");
            ctx.result(errorHtml);
        }
    }
    
    /**
     * Получает данные теста из формы
     */
    private Test getTestFromForm(Context ctx) {
        String title = ctx.formParam("title");
        String description = ctx.formParam("description");
        String testId = ctx.formParam("testId");
        String courseId = ctx.formParam("courseId");
        
        int minPoint;
        try {
            minPoint = Integer.parseInt(ctx.formParam("min_point"));
        } catch (NumberFormatException e) {
            minPoint = 60;
        }
        
        return new Test(testId, courseId, title, minPoint, description);
    }
    
    /**
     * Парсит данные формы в TestFormData DTO
     */
    private TestFormData parseTestFormData(Context ctx) {
        TestFormData formData = new TestFormData();
        Map<String, List<String>> formParams = ctx.formParamMap();
        
        // Собираем вопросы
        Map<Integer, TestFormData.QuestionFormData> questionsMap = new TreeMap<>();
        
        // Регулярные выражения для парсинга имен полей
        Pattern questionPattern = Pattern.compile("question\\[(\\d+)\\]\\[text\\]");
        Pattern answerTextPattern = Pattern.compile("question\\[(\\d+)\\]\\[answers\\]\\[(\\d+)\\]\\[text\\]");
        Pattern answerPointsPattern = Pattern.compile("question\\[(\\d+)\\]\\[answers\\]\\[(\\d+)\\]\\[points\\]");
        
        // 1. Собираем тексты вопросов
        for (Map.Entry<String, List<String>> entry : formParams.entrySet()) {
            String paramName = entry.getKey();
            Matcher questionMatcher = questionPattern.matcher(paramName);
            
            if (questionMatcher.matches()) {
                int questionIndex = Integer.parseInt(questionMatcher.group(1));
                String questionText = entry.getValue().get(0);
                
                if (questionText != null && !questionText.trim().isEmpty()) {
                    TestFormData.QuestionFormData question = questionsMap.getOrDefault(
                        questionIndex, 
                        new TestFormData.QuestionFormData()
                    );
                    question.setTextOfQuestion(questionText);
                    question.setOrder(questionIndex);
                    questionsMap.put(questionIndex, question);
                }
            }
        }
        
        // 2. Собираем тексты ответов
        for (Map.Entry<String, List<String>> entry : formParams.entrySet()) {
            String paramName = entry.getKey();
            Matcher answerMatcher = answerTextPattern.matcher(paramName);
            
            if (answerMatcher.matches()) {
                int questionIndex = Integer.parseInt(answerMatcher.group(1));
                int answerIndex = Integer.parseInt(answerMatcher.group(2));
                String answerText = entry.getValue().get(0);
                
                if (answerText != null && !answerText.trim().isEmpty()) {
                    TestFormData.QuestionFormData question = questionsMap.get(questionIndex);
                    if (question != null) {
                        List<TestFormData.AnswerFormData> answers = question.getAnswers();
                        if (answers == null) {
                            answers = new ArrayList<>();
                            question.setAnswers(answers);
                        }
                        
                        // Ищем или создаем ответ
                        while (answers.size() <= answerIndex) {
                            answers.add(new TestFormData.AnswerFormData());
                        }
                        
                        TestFormData.AnswerFormData answer = answers.get(answerIndex);
                        answer.setText(answerText);
                        answer.setScore(0); // Значение по умолчанию
                    }
                }
            }
        }
        
        // 3. Собираем баллы для ответов
        for (Map.Entry<String, List<String>> entry : formParams.entrySet()) {
            String paramName = entry.getKey();
            Matcher pointsMatcher = answerPointsPattern.matcher(paramName);
            
            if (pointsMatcher.matches()) {
                int questionIndex = Integer.parseInt(pointsMatcher.group(1));
                int answerIndex = Integer.parseInt(pointsMatcher.group(2));
                String pointsValue = entry.getValue().get(0);
                
                TestFormData.QuestionFormData question = questionsMap.get(questionIndex);
                if (question != null && question.getAnswers() != null) {
                    if (question.getAnswers().size() > answerIndex) {
                        TestFormData.AnswerFormData answer = question.getAnswers().get(answerIndex);
                        
                        try {
                            Integer score = Integer.parseInt(pointsValue);
                            answer.setScore(score);
                        } catch (NumberFormatException e) {
                            answer.setScore(0);
                        }
                    }
                }
            }
        }
        
        // 4. Преобразуем TreeMap в List
        List<TestFormData.QuestionFormData> questions = new ArrayList<>(questionsMap.values());
        formData.setQuestions(questions);
        
        return formData;
    }
    
    /**
     * Сохраняет вопросы и ответы в базу данных для теста
     */
    private void saveQuestionsAndAnswers(UUID testId, TestFormData formData) {
        // Удаляем старые вопросы
        deleteTestQuestions(testId);
        
        // Сохраняем новые вопросы
        for (TestFormData.QuestionFormData questionData : formData.getQuestions()) {
            if (questionData.getTextOfQuestion() != null && 
                !questionData.getTextOfQuestion().trim().isEmpty()) {
                
                // Создаем вопрос
                Question question = new Question();
                question.setTestId(testId);
                question.setTextOfQuestion(questionData.getTextOfQuestion());
                question.setOrder(questionData.getOrder());
                
                Question savedQuestion = questionService.createQuestion(question);
                
                // Сохраняем ответы
                if (questionData.getAnswers() != null) {
                    for (TestFormData.AnswerFormData answerData : questionData.getAnswers()) {
                        if (answerData.getText() != null && !answerData.getText().trim().isEmpty()) {
                            Answer answer = new Answer();
                            answer.setQuestionId(savedQuestion.getId());
                            answer.setText(answerData.getText());
                            
                            Integer score = answerData.getScore() != null ? answerData.getScore() : 0;
                            answer.setScore(score);
                            
                            answerService.createAnswer(answer);
                        }
                    }
                }
            }
        }
    }
    
    /**
     * Удаляет вопросы теста
     */
    private void deleteTestQuestions(UUID testId) {
        try {
            List<Question> existingQuestions = questionService.getQuestionsByTestId(testId);
            for (Question question : existingQuestions) {
                try {
                    // Удаляем ответы
                    List<Answer> answers = answerService.getAnswersByQuestionId(question.getId());
                    for (Answer answer : answers) {
                        answerService.deleteAnswer(answer.getId());
                    }
                } catch (Exception e) {
                    System.out.println("DEBUG: Error deleting test answers: " + e.getMessage());
                }
                
                try {
                    // Удаляем вопрос
                    questionService.deleteQuestion(question.getId());
                } catch (Exception e) {
                    System.out.println("DEBUG: Error deleting test question: " + e.getMessage());
                }
            }
            System.out.println("DEBUG: Deleted " + existingQuestions.size() + " test questions");
        } catch (Exception e) {
            System.out.println("DEBUG: Error deleting test questions: " + e.getMessage());
        }
    }
    
    /**
     * Получает вопросы из формы для отображения в случае ошибки
     */
    private List<QuestionWithAnswers> getQuestionsFromForm(Context ctx) {
        try {
            TestFormData formData = parseTestFormData(ctx);
            List<QuestionWithAnswers> result = new ArrayList<>();
            
            for (TestFormData.QuestionFormData questionData : formData.getQuestions()) {
                Question question = new Question();
                question.setTextOfQuestion(questionData.getTextOfQuestion());
                question.setOrder(questionData.getOrder());
                
                List<Answer> answers = new ArrayList<>();
                if (questionData.getAnswers() != null) {
                    for (TestFormData.AnswerFormData answerData : questionData.getAnswers()) {
                        Answer answer = new Answer();
                        answer.setText(answerData.getText());
                        answer.setScore(answerData.getScore() != null ? answerData.getScore() : 0);
                        answers.add(answer);
                    }
                }
                
                result.add(new QuestionWithAnswers(question, answers));
            }
            
            return result;
        } catch (Exception e) {
            return new ArrayList<>();
        }
    }
    
    /**
     * Получает URL для редиректа из атрибутов контекста
     */
    private String getRedirectUrl(Context ctx) {
        // Сначала проверяем атрибут, установленный в before фильтре
        String redirect = ctx.attribute("redirect");
        if (redirect == null) {
            // Если нет в атрибутах, проверяем query параметр
            redirect = ctx.queryParam("redirect");
        }
        return redirect;
    }
    
    /**
     * Выполняет редирект с сохранением всех параметров
     */
    private void performRedirect(Context ctx, String externalRedirect, String internalRedirect) {
        if (externalRedirect != null && !externalRedirect.trim().isEmpty()) {
            // Если есть внешний redirect URL, делаем редирект на него
            // Добавляем все query параметры, которые были в оригинальном запросе
            String fullRedirect = buildRedirectUrlWithParams(externalRedirect, ctx);
            ctx.redirect(fullRedirect);
        } else {
            // Если нет внешнего redirect, используем внутренний
            ctx.redirect(internalRedirect);
        }
    }
    
    /**
     * Строит URL для редиректа с сохранением всех параметров
     */
    private String buildRedirectUrlWithParams(String baseUrl, Context ctx) {
        StringBuilder urlBuilder = new StringBuilder(baseUrl);
        
        // Получаем все query параметры из оригинального запроса
        Map<String, List<String>> queryParams = ctx.queryParamMap();
        boolean firstParam = true;
        
        // Проверяем, есть ли уже query параметры в baseUrl
        if (baseUrl.contains("?")) {
            firstParam = false;
        }
        
        // Добавляем все параметры, кроме "redirect" (чтобы избежать циклических редиректов)
        for (Map.Entry<String, List<String>> entry : queryParams.entrySet()) {
            String paramName = entry.getKey();
            if (!"redirect".equals(paramName)) {
                for (String paramValue : entry.getValue()) {
                    if (firstParam) {
                        urlBuilder.append("?");
                        firstParam = false;
                    } else {
                        urlBuilder.append("&");
                    }
                    urlBuilder.append(paramName).append("=").append(paramValue);
                }
            }
        }
        
        return urlBuilder.toString();
    }
    
    /**
     * Добавляет все query параметры в модель для скрытых полей формы
     */
    private void addAllQueryParamsToModel(Context ctx, Map<String, Object> model) {
        Map<String, List<String>> queryParams = ctx.queryParamMap();
        Map<String, String> hiddenParams = new HashMap<>();
        
        for (Map.Entry<String, List<String>> entry : queryParams.entrySet()) {
            String paramName = entry.getKey();
            List<String> values = entry.getValue();
            if (!values.isEmpty()) {
                hiddenParams.put(paramName, values.get(0));
            }
        }
        
        model.put("hiddenParams", hiddenParams);
    }
}