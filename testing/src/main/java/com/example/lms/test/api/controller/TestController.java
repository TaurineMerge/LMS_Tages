package com.example.lms.test.api.controller;

import io.javalin.http.Context;
import io.javalin.http.HttpStatus;
import java.util.Map;
import java.util.HashMap;
import java.util.List;
import java.util.ArrayList;
import java.time.LocalDateTime;

public class TestController {

    // Временное хранилище для заглушек
    private static final Map<String, Map<String, Object>> testStore = new HashMap<>();
    
    static {
        // Инициализируем несколько тестов для демонстрации
        Map<String, Object> test1 = new HashMap<>();
        test1.put("id", "1");
        test1.put("title", "Тест по математике");
        test1.put("description", "Базовые математические операции");
        test1.put("questionsCount", 20);
        test1.put("durationMinutes", 45);
        test1.put("createdAt", "2024-01-15T10:30:00");
        test1.put("updatedAt", "2024-01-20T14:45:00");
        
        Map<String, Object> test2 = new HashMap<>();
        test2.put("id", "2");
        test2.put("title", "Тест по программированию");
        test2.put("description", "Основы Java и ООП");
        test2.put("questionsCount", 25);
        test2.put("durationMinutes", 60);
        test2.put("createdAt", "2024-01-16T09:15:00");
        test2.put("updatedAt", "2024-01-22T11:20:00");
        
        testStore.put("1", test1);
        testStore.put("2", test2);
    }

    public static void getTests(Context ctx) {
        // Получение всех тестов
        List<Map<String, Object>> tests = new ArrayList<>(testStore.values());
        
        Map<String, Object> response = new HashMap<>();
        response.put("status", "success");
        response.put("message", "Список тестов получен");
        response.put("data", tests);
        response.put("count", tests.size());
        
        ctx.status(HttpStatus.OK);
        ctx.json(response);
    }

    public static void createTest(Context ctx) {
        // Создание нового теста
        try {
            // В реальном приложении здесь будет парсинг тела запроса
            // Для заглушки просто создаем новый тест
            String newId = String.valueOf(testStore.size() + 1);
            
            Map<String, Object> newTest = new HashMap<>();
            newTest.put("id", newId);
            newTest.put("title", "Новый тест");
            newTest.put("description", "Описание нового теста");
            newTest.put("questionsCount", 0);
            newTest.put("durationMinutes", 30);
            newTest.put("createdAt", LocalDateTime.now().toString());
            newTest.put("updatedAt", LocalDateTime.now().toString());
            
            testStore.put(newId, newTest);
            
            Map<String, Object> response = new HashMap<>();
            response.put("status", "success");
            response.put("message", "Тест успешно создан");
            response.put("data", newTest);
            
            ctx.status(HttpStatus.CREATED);
            ctx.json(response);
            
        } catch (Exception e) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "Ошибка при создании теста");
            errorResponse.put("error", e.getMessage());
            
            ctx.status(HttpStatus.INTERNAL_SERVER_ERROR);
            ctx.json(errorResponse);
        }
    }

    public static void getTestById(Context ctx) {
        // Получение теста по ID
        String testId = ctx.pathParam("id");
        
        if (testId == null || testId.trim().isEmpty()) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "ID теста не указан");
            
            ctx.status(HttpStatus.BAD_REQUEST);
            ctx.json(errorResponse);
            return;
        }
        
        Map<String, Object> test = testStore.get(testId);
        
        if (test == null) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "Тест с ID " + testId + " не найден");
            
            ctx.status(HttpStatus.NOT_FOUND);
            ctx.json(errorResponse);
        } else {
            Map<String, Object> response = new HashMap<>();
            response.put("status", "success");
            response.put("message", "Тест найден");
            response.put("data", test);
            
            ctx.status(HttpStatus.OK);
            ctx.json(response);
        }
    }

    public static void updateTest(Context ctx) {
        // Обновление теста
        String testId = ctx.pathParam("id");
        
        if (testId == null || testId.trim().isEmpty()) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "ID теста не указан");
            
            ctx.status(HttpStatus.BAD_REQUEST);
            ctx.json(errorResponse);
            return;
        }
        
        if (!testStore.containsKey(testId)) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "Тест с ID " + testId + " не найден");
            
            ctx.status(HttpStatus.NOT_FOUND);
            ctx.json(errorResponse);
            return;
        }
        
        try {
            // В реальном приложении здесь будет парсинг тела запроса
            // Для заглушки просто обновляем время изменения
            Map<String, Object> test = testStore.get(testId);
            test.put("updatedAt", LocalDateTime.now().toString());
            test.put("title", "Обновленный тест"); // Пример обновления поля
            
            Map<String, Object> response = new HashMap<>();
            response.put("status", "success");
            response.put("message", "Тест успешно обновлен");
            response.put("data", test);
            
            ctx.status(HttpStatus.OK);
            ctx.json(response);
            
        } catch (Exception e) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "Ошибка при обновлении теста");
            errorResponse.put("error", e.getMessage());
            
            ctx.status(HttpStatus.INTERNAL_SERVER_ERROR);
            ctx.json(errorResponse);
        }
    }

    public static void deleteTest(Context ctx) {
        // Удаление теста
        String testId = ctx.pathParam("id");
        
        if (testId == null || testId.trim().isEmpty()) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "ID теста не указан");
            
            ctx.status(HttpStatus.BAD_REQUEST);
            ctx.json(errorResponse);
            return;
        }
        
        if (!testStore.containsKey(testId)) {
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("status", "error");
            errorResponse.put("message", "Тест с ID " + testId + " не найден");
            
            ctx.status(HttpStatus.NOT_FOUND);
            ctx.json(errorResponse);
            return;
        }
        
        testStore.remove(testId);
        
        Map<String, Object> response = new HashMap<>();
        response.put("status", "success");
        response.put("message", "Тест с ID " + testId + " успешно удален");
        
        ctx.status(HttpStatus.OK);
        ctx.json(response);
    }
}