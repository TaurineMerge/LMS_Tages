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
    private final Handlebars handlebars;
    
    public TestFormController(TestService testService,
                             QuestionService questionService,
                             AnswerService answerService,
                             Handlebars handlebars) {
        this.testService = testService;
        this.questionService = questionService;
        this.answerService = answerService;
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
            
            // Получаем тест
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
            model.put("title", "Редактирование теста: " + test.getTitle());
            model.put("test", test);
            model.put("questions", questionsWithAnswers);
            model.put("page", "test-builder");
            model.put("isNew", false);
            model.put("testId", testId);
            model.put("success", ctx.queryParam("success") != null);
            
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
     * POST /web/tests/save - Сохранение теста из формы (ОБНОВЛЕННЫЙ)
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
            saveQuestionsAndAnswers(savedTest.getId(), formData);
            
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
     * POST /web/tests/{id}/draft - Сохранение черновика
     */
    public void saveTestDraft(Context ctx) {
        try {
            String testId = ctx.pathParam("id");
            String title = ctx.formParam("title");
            String description = ctx.formParam("description");
            
            int minPoint;
            try {
                minPoint = Integer.parseInt(ctx.formParam("min_point"));
            } catch (NumberFormatException e) {
                minPoint = 60;
            }
            
            Test test = new Test(testId, null, title, minPoint, description);
            Test savedTest;
            
            if (testId == null || testId.isEmpty() || "new".equals(testId)) {
                savedTest = testService.createTest(test);
            } else {
                savedTest = testService.updateTest(test);
            }
            
            // Также сохраняем вопросы и ответы для черновика
            TestFormData formData = parseTestFormData(ctx);
            saveQuestionsAndAnswers(savedTest.getId(), formData);
            
            // Перезагружаем данные для отображения
            List<Question> questions = questionService.getQuestionsByTestId(UUID.fromString(savedTest.getId()));
            List<QuestionWithAnswers> questionsWithAnswers = questions.stream()
                .map(q -> {
                    List<Answer> answers = answerService.getAnswersByQuestionId(q.getId());
                    return new QuestionWithAnswers(q, answers);
                })
                .collect(Collectors.toList());
            
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Черновик сохранен");
            model.put("success", "Черновик успешно сохранен!");
            model.put("test", savedTest);
            model.put("questions", questionsWithAnswers);
            model.put("page", "test-builder");
            model.put("isNew", testId == null || testId.isEmpty() || "new".equals(testId));
            
            renderTemplateWithBody(ctx, "test-builder", model);
            
        } catch (Exception e) {
            Map<String, Object> model = new HashMap<>();
            model.put("title", "Ошибка сохранения черновика");
            model.put("error", "Ошибка при сохранении черновика: " + e.getMessage());
            model.put("test", getTestFromForm(ctx));
            model.put("questions", getQuestionsFromForm(ctx));
            model.put("page", "test-builder");
            model.put("isNew", true);
            
            renderTemplateWithBody(ctx, "test-builder", model);
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
        
        // 2. Собираем вопросы
        Map<Integer, TestFormData.QuestionFormData> questionsMap = new TreeMap<>();
        
        // Регулярные выражения для парсинга имен полей
        Pattern questionPattern = Pattern.compile("question\\[(\\d+)\\]\\[text\\]");
        Pattern answerPattern = Pattern.compile("question\\[(\\d+)\\]\\[answers\\]\\[(\\d+)\\]\\[text\\]");
        Pattern correctAnswerPattern = Pattern.compile("question\\[(\\d+)\\]\\[correct\\]");
        
        // Сначала собираем тексты вопросов
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
        
        // Затем собираем ответы для каждого вопроса
        for (Map.Entry<String, List<String>> entry : formParams.entrySet()) {
            String paramName = entry.getKey();
            Matcher answerMatcher = answerPattern.matcher(paramName);
            
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
                        // УБРАНО: answer.setOrder(answerIndex);
                        answer.setScore(0); // По умолчанию 0 баллов
                    }
                }
            }
        }
        
        // Собираем правильные ответы
        for (Map.Entry<String, List<String>> entry : formParams.entrySet()) {
            String paramName = entry.getKey();
            Matcher correctMatcher = correctAnswerPattern.matcher(paramName);
            
            if (correctMatcher.matches()) {
                int questionIndex = Integer.parseInt(correctMatcher.group(1));
                String correctAnswerValue = entry.getValue().get(0);
                
                if (correctAnswerValue != null && !correctAnswerValue.isEmpty()) {
                    try {
                        int correctAnswerIndex = Integer.parseInt(correctAnswerValue);
                        TestFormData.QuestionFormData question = questionsMap.get(questionIndex);
                        if (question != null && question.getAnswers() != null) {
                            if (question.getAnswers().size() > correctAnswerIndex) {
                                TestFormData.AnswerFormData answer = question.getAnswers().get(correctAnswerIndex);
                                answer.setIsCorrect(true);
                                // Правильный ответ получает 1 балл
                                answer.setScore(1);
                            }
                        }
                    } catch (NumberFormatException e) {
                        // Игнорируем некорректное значение
                    }
                }
            }
        }
        
        // 3. Преобразуем TreeMap в List (сохраняя порядок)
        List<TestFormData.QuestionFormData> questions = new ArrayList<>(questionsMap.values());
        formData.setQuestions(questions);
        
        return formData;
    }
    
    /**
     * Сохраняет вопросы и ответы в базу данных
     */
    private void saveQuestionsAndAnswers(String testId, TestFormData formData) {
        // Удаляем старые вопросы (если редактируем существующий тест)
        try {
            List<Question> existingQuestions = questionService.getQuestionsByTestId(UUID.fromString(testId));
            for (Question question : existingQuestions) {
                questionService.deleteQuestion(question.getId());
            }
        } catch (Exception e) {
            // Игнорируем ошибки при удалении (если вопросов еще нет)
        }
        
        // Сохраняем новые вопросы
        for (TestFormData.QuestionFormData questionData : formData.getQuestions()) {
            if (questionData.getTextOfQuestion() != null && 
                !questionData.getTextOfQuestion().trim().isEmpty()) {
                
                // Создаем вопрос
                Question question = new Question();
                question.setTestId(UUID.fromString(testId));
                question.setTextOfQuestion(questionData.getTextOfQuestion());
                question.setOrder(questionData.getOrder()); // Сохраняем order для вопроса
                
                Question savedQuestion = questionService.createQuestion(question);
                
                // Сохраняем ответы для этого вопроса
                if (questionData.getAnswers() != null) {
                    for (TestFormData.AnswerFormData answerData : questionData.getAnswers()) {
                        if (answerData.getText() != null && !answerData.getText().trim().isEmpty()) {
                            Answer answer = new Answer();
                            answer.setQuestionId(savedQuestion.getId());
                            answer.setText(answerData.getText());
                            
                            // Устанавливаем баллы
                            Integer score = answerData.getScore() != null ? answerData.getScore() : 0;
                            answer.setScore(score);
                            
                            // УБРАНО: answer.setOrder(answerData.getOrder());
                            
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
                question.setOrder(questionData.getOrder()); // Сохраняем order для вопроса
                
                List<Answer> answers = new ArrayList<>();
                if (questionData.getAnswers() != null) {
                    for (TestFormData.AnswerFormData answerData : questionData.getAnswers()) {
                        Answer answer = new Answer();
                        answer.setText(answerData.getText());
                        answer.setScore(answerData.getScore() != null ? answerData.getScore() : 0);
                        // УБРАНО: answer.setOrder(answerData.getOrder());
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