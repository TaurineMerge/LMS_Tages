package com.example.lms.internal.service;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.internal.api.dto.AttemptDetail;
import com.example.lms.internal.api.dto.AttemptsListItem;
import com.example.lms.internal.api.dto.PerTestStats;
import com.example.lms.internal.api.dto.UserStats;
import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.service.TestAttemptService;

/**
 * Сервис для работы с Internal API.
 * Агрегирует данные из TestAttemptService и TestService
 * для предоставления данных Python-сервису.
 */
public class InternalApiService {

    private static final Logger logger = LoggerFactory.getLogger(InternalApiService.class);

    private final TestAttemptService testAttemptService;
    private final TestService testService;

    public InternalApiService(TestAttemptService testAttemptService, TestService testService) {
        this.testAttemptService = testAttemptService;
        this.testService = testService;
    }

    /**
     * Получить детальную информацию о попытке по ID.
     *
     * @param attemptId ID попытки
     * @return AttemptDetailDto или null если не найдено
     */
    public AttemptDetail getAttemptDetail(UUID attemptId) {
        logger.debug("Получение деталей попытки: {}", attemptId);

        TestAttemptModel attempt = testAttemptService.findById(attemptId).orElse(null);
        if (attempt == null) {
            logger.warn("Попытка {} не найдена", attemptId);
            return null;
        }

        // Получаем информацию о тесте для определения passed
        Test test = getTestOrNull(attempt.getTestId());
        Boolean passed = calculatePassed(attempt, test);

        return new AttemptDetail(
                attempt.getId(),
                attempt.getStudentId(),
                attempt.getTestId(),
                attempt.getDateOfAttempt().toString(),
                attempt.getPoint(),
                attempt.getCompleted(),
                passed,
                attempt.getCertificateId(), // всегда null
                null, // attempt_version - пропускаем
                attempt.getAttemptSnapshot(),
                null // meta - пока null
        );
    }

    /**
     * Получить список всех попыток пользователя.
     *
     * @param userId ID пользователя (студента)
     * @return список попыток
     */
    public List<AttemptsListItem> getUserAttempts(UUID userId) {
        logger.debug("Получение попыток пользователя: {}", userId);

        List<TestAttemptModel> attempts = testAttemptService.getTestAttemptsByStudentId(userId);

        // Кэш для тестов, чтобы не запрашивать каждый раз
        Map<UUID, Test> testCache = new HashMap<>();

        return attempts.stream()
                .map(attempt -> {
                    Test test = testCache.computeIfAbsent(
                            attempt.getTestId(),
                            this::getTestOrNull);

                    Boolean passed = calculatePassed(attempt, test);

                    return new AttemptsListItemDto(
                            attempt.getId(),
                            attempt.getTestId(),
                            attempt.getDateOfAttempt().toString(),
                            attempt.getPoint(),
                            attempt.getCompleted(),
                            passed,
                            attempt.getCertificateId(),
                            attempt.getAttemptSnapshot());
                })
                .collect(Collectors.toList());
    }

    /**
     * Получить статистику пользователя по всем попыткам.
     *
     * @param userId ID пользователя
     * @return статистика пользователя
     */
    public UserStats getUserStats(UUID userId) {
        logger.debug("Получение статистики пользователя: {}", userId);

        List<TestAttemptModel> allAttempts = testAttemptService.getTestAttemptsByStudentId(userId);

        // Общая статистика
        int attemptsTotal = allAttempts.size();
        Integer bestScore = allAttempts.stream()
                .map(TestAttemptModel::getPoint)
                .filter(point -> point != null)
                .max(Integer::compareTo)
                .orElse(null);

        // Последняя попытка
        String lastAttemptAt = allAttempts.stream()
                .filter(a -> a.getCompleted() != null && a.getCompleted())
                .map(a -> a.getDateOfAttempt().atStartOfDay())
                .max(LocalDateTime::compareTo)
                .map(dt -> dt.format(DateTimeFormatter.ISO_DATE_TIME))
                .orElse(null);

        // Группируем по тестам
        Map<UUID, List<TestAttemptModel>> attemptsByTest = allAttempts.stream()
                .collect(Collectors.groupingBy(TestAttemptModel::getTestId));

        // Кэш тестов
        Map<UUID, Test> testCache = new HashMap<>();

        // Считаем пройденные попытки
        int attemptsPassed = 0;
        for (TestAttemptModel attempt : allAttempts) {
            Test test = testCache.computeIfAbsent(attempt.getTestId(), this::getTestOrNull);
            if (calculatePassed(attempt, test)) {
                attemptsPassed++;
            }
        }

        // Статистика по каждому тесту
        List<PerTestStats> perTestStats = new ArrayList<>();

        for (Map.Entry<UUID, List<TestAttemptModel>> entry : attemptsByTest.entrySet()) {
            UUID testId = entry.getKey();
            List<TestAttemptModel> testAttempts = entry.getValue();

            Test test = testCache.computeIfAbsent(testId, this::getTestOrNull);

            String testTitle = test != null ? test.getTitle() : "Unknown Test";

            Integer testBestScore = testAttempts.stream()
                    .map(TestAttemptModel::getPoint)
                    .filter(point -> point != null)
                    .max(Integer::compareTo)
                    .orElse(null);

            int passedCount = (int) testAttempts.stream()
                    .filter(attempt -> calculatePassed(attempt, test))
                    .count();

            perTestStats.add(new PerTestStats(
                    testId,
                    testTitle,
                    testAttempts.size(),
                    testBestScore,
                    passedCount));
        }

        return new UserStats(
                userId,
                attemptsTotal,
                attemptsPassed,
                bestScore,
                lastAttemptAt,
                perTestStats);
    }

    // ------------------------------------------------------------------
    // ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
    // ------------------------------------------------------------------

    /**
     * Получить тест по ID или вернуть null если не найден.
     */
    private Test getTestOrNull(UUID testId) {
        try {
            return testService.getTestById(testId.toString());
        } catch (Exception e) {
            logger.warn("Тест {} не найден: {}", testId, e.getMessage());
            return null;
        }
    }

    /**
     * Вычислить значение passed для попытки.
     * 
     * @param attempt попытка
     * @param test    тест (может быть null)
     * @return true если тест пройден, false если нет, null если не завершён
     */
    private Boolean calculatePassed(TestAttemptModel attempt, Test test) {
        if (attempt.getCompleted() == null || !attempt.getCompleted()) {
            return null; // Не завершено
        }

        if (attempt.getPoint() == null) {
            return null; // Нет баллов
        }

        if (test == null) {
            // Тест не найден - не можем определить
            return null;
        }

        Integer minPoint = test.getMin_point();
        if (minPoint == null) {
            // Нет минимального порога - считаем пройденным
            return true;
        }

        return attempt.getPoint() >= minPoint;
    }
}