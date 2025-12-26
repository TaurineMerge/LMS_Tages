package com.example.lms.internal.service;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.*;
import java.util.stream.Collectors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.internal.api.dto.*;
import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test_attempt.api.dto.TestAttempt;
import com.example.lms.test_attempt.domain.service.TestAttemptService;
import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.service.DraftService;

/**
 * Фасад-сервис для внутреннего API, предоставляющий агрегированные данные о тестах и попытках их прохождения.
 * <p>
 * Этот сервис служит слоем абстракции между внешними потребителями внутреннего API и доменной логикой приложения.
 * Он скрывает детали реализации доменной модели и предоставляет оптимизированные представления данных.
 * </p>
 * <p>
 * Основные функции:
 * <ul>
 *   <li>Получение детальной информации о попытках прохождения тестов</li>
 *   <li>Агрегация статистики пользователей</li>
 *   <li>Предоставление данных о тестах и черновиках по курсам</li>
 *   <li>Кэширование часто запрашиваемых данных для оптимизации производительности</li>
 * </ul>
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see TestAttemptService
 * @see TestService
 * @see DraftService
 */
public class InternalApiService {

    private static final Logger logger = LoggerFactory.getLogger(InternalApiService.class);

    private final TestAttemptService testAttemptService;
    private final TestService testService;
    private final DraftService draftService;

    /**
     * Создает новый экземпляр InternalApiService с указанными зависимостями.
     *
     * @param testAttemptService сервис для работы с попытками прохождения тестов
     * @param testService сервис для работы с тестами
     * @param draftService сервис для работы с черновиками тестов
     * @throws NullPointerException если любой из параметров равен null
     */
    public InternalApiService(
            TestAttemptService testAttemptService,
            TestService testService,
            DraftService draftService) {
        this.testAttemptService = testAttemptService;
        this.testService = testService;
        this.draftService = draftService;
    }

    /**
     * Получает детальную информацию о конкретной попытке прохождения теста.
     *
     * @param attemptId уникальный идентификатор попытки
     * @return объект {@link AttemptDetail} с детальной информацией о попытке или {@code null},
     *         если попытка с указанным идентификатором не найдена
     */
    public AttemptDetail getAttemptDetail(UUID attemptId) {
        logger.debug("Получение деталей попытки: {}", attemptId);

        TestAttempt attempt = safeGetAttempt(attemptId);
        if (attempt == null) {
            return null;
        }
        UUID testId = UUID.fromString(attempt.getTest_id());
        Test test = getTestOrNull(testId);
        Boolean passed = calculatePassed(attempt, test);

        return mapToAttemptDetail(attempt, passed);
    }

    /**
     * Получает список всех попыток прохождения тестов для указанного пользователя.
     *
     * @param userId уникальный идентификатор пользователя
     * @return список объектов {@link AttemptsListItem}, представляющих попытки пользователя.
     *         Возвращает пустой список, если у пользователя нет попыток
     */
    public List<AttemptsListItem> getUserAttempts(UUID userId) {
        logger.debug("Получение попыток пользователя: {}", userId);

        List<TestAttempt> attempts = testAttemptService.getTestAttemptsByStudentId(userId.toString());
        Map<UUID, Test> testCache = new HashMap<>();

        return attempts.stream()
                .map(attempt -> {
                    UUID testId = UUID.fromString(attempt.getTest_id());
                    Test test = getCachedTest(testCache, testId);
                    Boolean passed = calculatePassed(attempt, test);
                    return mapToAttemptsListItem(attempt, passed);
                })
                .toList();
    }

    /**
     * Получает агрегированную статистику пользователя по всем попыткам прохождения тестов.
     *
     * @param userId уникальный идентификатор пользователя
     * @return объект {@link UserStats} со статистикой пользователя
     */
    public UserStats getUserStats(UUID userId) {
        logger.debug("Получение статистики пользователя: {}", userId);

        List<TestAttempt> attempts = testAttemptService.getTestAttemptsByStudentId(userId.toString());
        Map<UUID, Test> testCache = new HashMap<>();

        int attemptsTotal = attempts.size();
        int attemptsPassed = countPassed(attempts, testCache);

        Integer bestScore = attempts.stream()
                .map(TestAttempt::getPoint)
                .filter(Objects::nonNull)
                .max(Integer::compareTo)
                .orElse(null);

        String lastAttemptAt = getLastCompletedAttemptDate(attempts);
        List<PerTestStats> perTestStats = buildPerTestStats(attempts, testCache);

        return new UserStats(
                userId,
                attemptsTotal,
                attemptsPassed,
                bestScore,
                lastAttemptAt,
                perTestStats);
    }

    /**
     * Получает информацию о тесте, связанном с указанным курсом.
     * <p>
     * Предполагается, что с каждым курсом связан только один тест.
     * Если тест не найден, возвращается ответ со статусом "not_found".
     *
     * @param courseId уникальный идентификатор курса
     * @return объект {@link CourseTestResponse} с данными теста или статусом "not_found"
     */
    public CourseTestResponse getTestByCourseId(UUID courseId) {
        logger.debug("Получение теста для курса: {}", courseId);

        List<Test> tests = testService.getTestsByCourseId(courseId.toString());
        if (tests.isEmpty()) {
            return new CourseTestResponse(null, "not_found");
        }

        // Предполагаем, что с курсом связан только один тест, берем первый
        Test test = tests.get(0);
        CourseTestResponse.TestData testData = new CourseTestResponse.TestData(
                UUID.fromString(test.getId()),
                courseId,
                test.getTitle(),
                test.getMin_point(),
                test.getDescription());

        return new CourseTestResponse(testData, "success");
    }

    /**
     * Получает информацию о черновике теста, связанном с указанным курсом.
     * <p>
     * Предполагается, что с каждым курсом связан только один черновик теста.
     * Если черновик не найден, возвращается ответ со статусом "not_found".
     *
     * @param courseId уникальный идентификатор курса
     * @return объект {@link CourseDraftResponse} с данными черновика или статусом "not_found"
     */
    public CourseDraftResponse getDraftByCourseId(UUID courseId) {
        logger.debug("Получение черновика для курса: {}", courseId);

        List<Draft> drafts = draftService.getDraftsByCourseId(courseId);
        if (drafts.isEmpty()) {
            return new CourseDraftResponse(null, "not_found");
        }

        // Предполагаем, что с курсом связан только один черновик, берем первый
        Draft draft = drafts.get(0);
        CourseDraftResponse.DraftData draftData = new CourseDraftResponse.DraftData(
                draft.getId(),
                draft.getTestId(),
                courseId,
                draft.getTitle(),
                draft.getMin_point(),
                draft.getDescription());

        return new CourseDraftResponse(draftData, "success");
    }

    /**
     * Подсчитывает количество успешно пройденных попыток.
     *
     * @param attempts список попыток для анализа
     * @param testCache кэш объектов тестов для оптимизации производительности
     * @return количество успешно пройденных попыток
     */
    private int countPassed(List<TestAttempt> attempts, Map<UUID, Test> testCache) {
        int passed = 0;
        for (TestAttempt attempt : attempts) {
            UUID testId = UUID.fromString(attempt.getTest_id());
            Test test = getCachedTest(testCache, testId);

            if (Boolean.TRUE.equals(calculatePassed(attempt, test))) {
                passed++;
            }
        }
        return passed;
    }

    /**
     * Строит статистику по каждому тесту на основе попыток пользователя.
     *
     * @param attempts список всех попыток пользователя
     * @param testCache кэш объектов тестов для оптимизации производительности
     * @return список объектов {@link PerTestStats} со статистикой по каждому тесту
     */
    private List<PerTestStats> buildPerTestStats(List<TestAttempt> attempts, Map<UUID, Test> testCache) {
        Map<UUID, List<TestAttempt>> byTest = attempts.stream()
                .collect(Collectors.groupingBy(a -> UUID.fromString(a.getTest_id())));

        List<PerTestStats> stats = new ArrayList<>();

        for (var entry : byTest.entrySet()) {
            UUID testId = entry.getKey();
            List<TestAttempt> testAttempts = entry.getValue();

            Test test = getCachedTest(testCache, testId);
            String title = test != null ? test.getTitle() : "Unknown Test";

            Integer bestScore = testAttempts.stream()
                    .map(TestAttempt::getPoint)
                    .filter(Objects::nonNull)
                    .max(Integer::compareTo)
                    .orElse(null);

            int passedCount = (int) testAttempts.stream()
                    .filter(a -> Boolean.TRUE.equals(calculatePassed(a, test)))
                    .count();

            stats.add(new PerTestStats(
                    testId,
                    title,
                    testAttempts.size(),
                    bestScore,
                    passedCount));
        }

        return stats;
    }

    /**
     * Получает дату последней завершенной попытки в формате ISO.
     *
     * @param attempts список попыток для анализа
     * @return дата последней завершенной попытки в формате "YYYY-MM-DD" или {@code null},
     *         если завершенных попыток нет
     */
    private String getLastCompletedAttemptDate(List<TestAttempt> attempts) {
        return attempts.stream()
                .filter(a -> Boolean.TRUE.equals(a.getCompleted()))
                .map(TestAttempt::getDate_of_attempt)
                .filter(Objects::nonNull)
                .map(LocalDate::parse)
                .max(LocalDate::compareTo)
                .map(d -> d.format(DateTimeFormatter.ISO_DATE))
                .orElse(null);
    }

    /**
     * Преобразует объект попытки прохождения теста в детальное представление.
     *
     * @param attempt объект попытки прохождения теста
     * @param passed флаг, указывающий пройден ли тест успешно
     * @return объект {@link AttemptDetail} с детальной информацией о попытке
     */
    private AttemptDetail mapToAttemptDetail(TestAttempt attempt, Boolean passed) {
        return new AttemptDetail(
                UUID.fromString(attempt.getId()),
                UUID.fromString(attempt.getStudent_id()),
                UUID.fromString(attempt.getTest_id()),
                attempt.getDate_of_attempt(),
                attempt.getPoint(),
                attempt.getCompleted(),
                passed,
                parseUuid(attempt.getCertificate_id()),
                attempt.getAttempt_version(),
                attempt.getAttempt_snapshot(),
                null);
    }

    /**
     * Преобразует объект попытки прохождения теста в элемент списка попыток.
     *
     * @param attempt объект попытки прохождения теста
     * @param passed флаг, указывающий пройден ли тест успешно
     * @return объект {@link AttemptsListItem} с основной информацией о попытке
     */
    private AttemptsListItem mapToAttemptsListItem(TestAttempt attempt, Boolean passed) {
        return new AttemptsListItem(
                UUID.fromString(attempt.getId()),
                UUID.fromString(attempt.getTest_id()),
                attempt.getDate_of_attempt(),
                attempt.getPoint(),
                attempt.getCompleted(),
                passed,
                parseUuid(attempt.getCertificate_id()),
                attempt.getAttempt_snapshot());
    }

    /**
     * Получает тест из кэша или загружает его при отсутствии.
     *
     * @param cache кэш объектов тестов
     * @param testId уникальный идентификатор теста
     * @return объект теста или {@code null}, если тест не найден
     */
    private Test getCachedTest(Map<UUID, Test> cache, UUID testId) {
        return cache.computeIfAbsent(testId, this::getTestOrNull);
    }

    /**
     * Безопасно получает тест по идентификатору, возвращая {@code null} при ошибке.
     *
     * @param testId уникальный идентификатор теста
     * @return объект теста или {@code null}, если тест не найден или произошла ошибка
     */
    private Test getTestOrNull(UUID testId) {
        try {
            return testService.getTestById(testId.toString());
        } catch (Exception e) {
            logger.warn("Тест {} не найден", testId);
            return null;
        }
    }

    /**
     * Безопасно получает попытку прохождения теста по идентификатору.
     *
     * @param attemptId уникальный идентификатор попытки
     * @return объект попытки или {@code null}, если попытка не найдена
     */
    private TestAttempt safeGetAttempt(UUID attemptId) {
        try {
            return testAttemptService.getTestAttemptById(attemptId);
        } catch (Exception e) {
            return null;
        }
    }

    /**
     * Определяет, была ли попытка прохождения теста успешной.
     * <p>
     * Попытка считается успешной, если:
     * <ol>
     *   <li>Она была завершена ({@code completed == true})</li>
     *   <li>Набрано количество баллов не меньше минимального порога теста</li>
     * </ol>
     *
     * @param attempt попытка прохождения теста
     * @param test тест, для которого выполнена попытка
     * @return {@code true} если тест пройден успешно, {@code false} если не пройден,
     *         {@code null} если результат не может быть определен
     */
    private Boolean calculatePassed(TestAttempt attempt, Test test) {
        if (!Boolean.TRUE.equals(attempt.getCompleted())) {
            return null;
        }

        if (attempt.getPoint() == null || test == null) {
            return null;
        }

        Integer minPoint = test.getMin_point();
        return minPoint == null || attempt.getPoint() >= minPoint;
    }

    /**
     * Преобразует строковое представление UUID в объект UUID.
     *
     * @param value строковое представление UUID или {@code null}
     * @return объект UUID или {@code null}, если входная строка равна {@code null}
     * @throws IllegalArgumentException если строка не соответствует формату UUID
     */
    private UUID parseUuid(String value) {
        return value != null ? UUID.fromString(value) : null;
    }
}