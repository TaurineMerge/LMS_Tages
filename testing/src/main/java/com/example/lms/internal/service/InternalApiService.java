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
import java.util.List;

/**
 * Facade-сервис для Internal API.
 * Агрегирует данные тестов и попыток,
 * не раскрывая доменную модель.
 */
public class InternalApiService {

	private static final Logger logger = LoggerFactory.getLogger(InternalApiService.class);

	private final TestAttemptService testAttemptService;
	private final TestService testService;
	private final DraftService draftService;

	public InternalApiService(
			TestAttemptService testAttemptService,
			TestService testService,
			DraftService draftService) {
		this.testAttemptService = testAttemptService;
		this.testService = testService;
		this.draftService = draftService;
	}

	// ------------------------------------------------------------------
	// PUBLIC API
	// ------------------------------------------------------------------

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

	public CourseTestResponse getTestByCourseId(UUID courseId) {
		logger.debug("Получение теста для курса: {}", courseId);

		List<Test> tests = testService.getTestsByCourseId(courseId.toString());
		if (tests.isEmpty()) {
			return new CourseTestResponse(null,"not_found");
		}

		// Assuming one test per course, take the first one
		Test test = tests.get(0);
		CourseTestResponse.TestData testData = new CourseTestResponse.TestData(
				UUID.fromString(test.getId()),
				courseId,
				test.getTitle(),
				test.getMin_point(),
				test.getDescription());

		return new CourseTestResponse(testData, "success");
	}

	public CourseDraftResponse getDraftByCourseId(UUID courseId) {
		logger.debug("Получение черновика для курса: {}", courseId);

		List<Draft> drafts = draftService.getDraftsByCourseId(courseId);
		if (drafts.isEmpty()) {
			return new CourseDraftResponse(null, "not_found");
		}

		// Assuming one draft per course, take the first one
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

	// ------------------------------------------------------------------
	// AGGREGATION HELPERS
	// ------------------------------------------------------------------

	private int countPassed(
			List<TestAttempt> attempts,
			Map<UUID, Test> testCache) {
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

	private List<PerTestStats> buildPerTestStats(
			List<TestAttempt> attempts,
			Map<UUID, Test> testCache) {
		Map<UUID, List<TestAttempt>> byTest = attempts.stream()
				.collect(Collectors.groupingBy(
						a -> UUID.fromString(a.getTest_id())));

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
					.filter(a -> Boolean.TRUE.equals(
							calculatePassed(a, test)))
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

	private String getLastCompletedAttemptDate(
			List<TestAttempt> attempts) {
		return attempts.stream()
				.filter(a -> Boolean.TRUE.equals(a.getCompleted()))
				.map(TestAttempt::getDate_of_attempt)
				.filter(Objects::nonNull)
				.map(LocalDate::parse)
				.max(LocalDate::compareTo)
				.map(d -> d.format(DateTimeFormatter.ISO_DATE))
				.orElse(null);
	}

	// ------------------------------------------------------------------
	// MAPPERS
	// ------------------------------------------------------------------

	private AttemptDetail mapToAttemptDetail(
			TestAttempt attempt,
			Boolean passed) {
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

	private AttemptsListItem mapToAttemptsListItem(
			TestAttempt attempt,
			Boolean passed) {
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

	// ------------------------------------------------------------------
	// DOMAIN ACCESS
	// ------------------------------------------------------------------

	private Test getCachedTest(
			Map<UUID, Test> cache,
			UUID testId) {
		return cache.computeIfAbsent(testId, this::getTestOrNull);
	}

	private Test getTestOrNull(UUID testId) {
		try {
			return testService.getTestById(testId.toString());
		} catch (Exception e) {
			logger.warn("Тест {} не найден", testId);
			return null;
		}
	}

	private TestAttempt safeGetAttempt(UUID attemptId) {
		try {
			return testAttemptService.getTestAttemptById(attemptId);
		} catch (Exception e) {
			return null;
		}
	}

	// ------------------------------------------------------------------
	// BUSINESS RULE
	// ------------------------------------------------------------------

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

	// ------------------------------------------------------------------
	// UTIL
	// ------------------------------------------------------------------

	private UUID parseUuid(String value) {
		return value != null ? UUID.fromString(value) : null;
	}
}
