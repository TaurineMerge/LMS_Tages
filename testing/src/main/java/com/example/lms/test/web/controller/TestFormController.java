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
        ctx.redirect("/web/tests/new");
    }
    
    /**
     * GET /web/tests/new - Форма создания нового теста
     */
    public void showNewTestForm(Context ctx) {
        Map<String, Object> model = new HashMap<>();
        model.put("title", "Создание нового теста");
        model.put("test", new Test(null, null, "", 60, ""));
        model.put("questions", new ArrayList<>());
        model.put("page", "test-builder");
        model.put("isNew", true);
        model.put("success", ctx.queryParam("success") != null);
        
        renderTemplateWithBody(ctx, "test-builder", model);
    }
    
    /**
     * GET /web/tests/{id}/edit - Форма редактирования теста
     */
    public void showEditTestForm(Context ctx) {
        try {
            String testId = ctx.pathParam("id");
            UUID testUuid = UUID.fromString(testId);
            
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

            // Пытаемся получить черновик для этого теста (опционально)
            boolean hasDraft = false;
            try {
                Draft draft = draftService.getDraftByTestId(testUuid);
                if (draft != null) {
                    hasDraft = true;
                    model.put("hasDraft", true);
                    model.put("draft", draft);
                }
            } catch (Exception e) {
                // Если черновика нет - игнорируем
                System.out.println("DEBUG: No draft found for test: " + e.getMessage());
            }
            
            model.put("title", "Редактирование теста: " + test.getTitle());
            model.put("test", test);
            model.put("questions", questionsWithAnswers);
            model.put("page", "test-builder");
            model.put("isNew", false);
            model.put("testId", testId);
            model.put("success", ctx.queryParam("success") != null);
            model.put("hasDraft", hasDraft);
            
            renderTemplateWithBody(ctx, "test-builder", model);
        } catch (Exception e) {
            e.printStackTrace();
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка загрузки");
            model.put("error", "Тест не найден. Создайте новый тест.");
            model.put("test", new Test(null, null, "", 60, ""));
            model.put("questions", new ArrayList<>());
            model.put("page", "test-builder");
            model.put("isNew", true);
            
            renderTemplateWithBody(ctx, "test-builder", model);
        }
    }
    
    /**
     * POST /web/tests/save - Сохранение теста из формы
     */
    public void saveTestFromForm(Context ctx) {
        try {
            // 1. Получаем и валидируем данные теста
            Test test = getTestFromForm(ctx);
            if (test.getTitle() == null || test.getTitle().trim().isEmpty()) {
                throw new IllegalArgumentException("Название теста обязательно");
            }
            
            // 2. Получаем вопросы и ответы из формы
            TestFormData formData = parseTestFormData(ctx);
            
            // 3. Сохраняем тест в базу
            Test savedTest;
            if (test.getId() == null || test.getId().isEmpty()) {
                savedTest = testService.createTest(test);
            } else {
                savedTest = testService.updateTest(test);
            }
            
            // 4. Сохраняем вопросы и ответы
            UUID testUuid = UUID.fromString(savedTest.getId());
            saveQuestionsAndAnswers(testUuid, formData);
            
            // 5. Перенаправляем на страницу редактирования
            ctx.redirect("/web/tests/" + savedTest.getId() + "/edit?success=true");
            
        } catch (Exception e) {
            e.printStackTrace();
            // При ошибке показываем форму снова
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка сохранения");
            model.put("error", "Ошибка при сохранении теста: " + e.getMessage());
            model.put("test", getTestFromForm(ctx));
            model.put("questions", getQuestionsFromForm(ctx));
            model.put("page", "test-builder");
            model.put("isNew", true);
            
            renderTemplateWithBody(ctx, "test-builder", model);
        }
    }
    
    /**
     * POST /web/tests/draft - Сохранение черновика нового теста
     */
    public void saveNewTestDraft(Context ctx) {
        saveDraftInternal(ctx, null, true);
    }

    /**
     * POST /web/tests/{id}/draft - Сохранение черновика существующего теста
     */
    public void saveExistingTestDraft(Context ctx) {
        String testId = ctx.pathParam("id");
        UUID testUuid = null;
        if (testId != null && !testId.trim().isEmpty()) {
            try {
                testUuid = UUID.fromString(testId);
            } catch (IllegalArgumentException e) {
                // Невалидный UUID, сохраняем как черновик без test_id
            }
        }
        saveDraftInternal(ctx, testUuid, false);
    }

    /**
     * Внутренний метод для сохранения черновика
     */
    private void saveDraftInternal(Context ctx, UUID testId, boolean isNewDraft) {
        try {
            System.out.println("DEBUG: saveDraftInternal called with testId = '" + testId + "', isNewDraft = " + isNewDraft);
            
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
                
                // СОХРАНЯЕМ ВОПРОСЫ С DRAFT_ID
                for (int i = 0; i < formData.getQuestions().size(); i++) {
                    TestFormData.QuestionFormData questionData = formData.getQuestions().get(i);
                    
                    if (questionData.getTextOfQuestion() != null && 
                        !questionData.getTextOfQuestion().trim().isEmpty()) {
                        
                        // Создаем вопрос ДЛЯ ЧЕРНОВИКА
                        Question question = new Question();
                        question.setDraftId(savedDraft.getId()); // ← UUID!
                        
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
                System.out.println("DEBUG: Questions saved as draft");
            }
            
            // 3. Подготавливаем данные для отображения
            Map<String, Object> model = new HashMap<>();
            
            // Создаем временный объект теста для отображения
            Test testForDisplay = new Test(null, null, title, minPoint, description);
            
            String displayTestId;
            boolean isNewForDisplay;
            
            if (savedDraft.getTestId() != null) {
                displayTestId = savedDraft.getTestId().toString();
                isNewForDisplay = false;
            } else {
                displayTestId = "draft-" + savedDraft.getId();
                isNewForDisplay = true;
            }
            
            // Загружаем вопросы черновика для отображения
            List<QuestionWithAnswers> questionsWithAnswers = new ArrayList<>();
            try {
                List<Question> questions = questionService.getQuestionsByDraftId(savedDraft.getId());
                for (Question question : questions) {
                    List<Answer> answers = answerService.getAnswersByQuestionId(question.getId());
                    questionsWithAnswers.add(new QuestionWithAnswers(question, answers));
                }
            } catch (Exception e) {
                System.out.println("DEBUG: Error loading draft questions: " + e.getMessage());
            }
            
            // Настраиваем модель
            model.put("title", "Черновик сохранен");
            model.put("success", "Черновик успешно сохранен! (ID черновика: " + savedDraft.getId() + ")");
            model.put("test", testForDisplay);
            model.put("questions", questionsWithAnswers);
            model.put("page", "test-builder");
            model.put("isNew", isNewForDisplay);
            model.put("testId", displayTestId);
            
            System.out.println("DEBUG: Rendering template with " + questionsWithAnswers.size() + " questions");
            renderTemplateWithBody(ctx, "test-builder", model);
            
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
            
            model.put("page", "test-builder");
            model.put("isNew", true);
            
            renderTemplateWithBody(ctx, "test-builder", model);
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
     * GET /web/tests/{id}/preview - Предпросмотр теста
     */
    public void previewTest(Context ctx) {
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
            
            renderTemplateWithBody(ctx, "test-preview", model);
        } catch (Exception e) {
            e.printStackTrace();
            ctx.redirect("/web/tests/new");
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
        
        int minPoint;
        try {
            minPoint = Integer.parseInt(ctx.formParam("min_point"));
        } catch (NumberFormatException e) {
            minPoint = 60;
        }
        
        return new Test(testId, null, title, minPoint, description);
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
     * Сохраняет вопросы и ответы в базу данных
     */
    private void saveQuestionsAndAnswers(UUID testId, TestFormData formData) {
        // Удаляем старые вопросы
        try {
            List<Question> existingQuestions = questionService.getQuestionsByTestId(testId);
            for (Question question : existingQuestions) {
                questionService.deleteQuestion(question.getId());
            }
        } catch (Exception e) {
            // Игнорируем ошибки при удалении
        }
        
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
}