package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.repository.TestAttemptRepositoryInterface;

public class TestAttemptService {

	private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);

	private final TestAttemptRepositoryInterface testAttemptRepository;

	public TestAttemptService(TestAttemptRepositoryInterface testAttemptRepository) {
		this.testAttemptRepository = testAttemptRepository;
	}

	// Существующие методы...
	public TestAttemptModel createTestAttempt(UUID studentId, UUID testId, String attemptVersion) {
		TestAttemptModel testAttempt = new TestAttemptModel(studentId, testId);
		testAttempt.validate();
		return testAttemptRepository.save(testAttempt);
	}

	public TestAttemptModel getTestAttemptById(UUID id) {
		return testAttemptRepository.findById(id)
				.orElseThrow(() -> new IllegalArgumentException("No test attempt found with id: " + id));
	}

	public TestAttemptModel updateTestAttempt(TestAttemptModel testAttempt) {
		UUID id = testAttempt.getId();
		if (id == null) {
			throw new IllegalArgumentException("Test Attempt ID cannot be null for update");
		}

		TestAttemptModel existingAttempt = testAttemptRepository.findById(id)
				.orElseThrow(() -> new IllegalArgumentException("No test attempt found with id: " + id));

		existingAttempt.setPoint(testAttempt.getPoint());
		existingAttempt.validate();

		return testAttemptRepository.update(existingAttempt);
	}

	public void completeTestAttempt(UUID id, int points) {
		TestAttemptModel testAttempt = getTestAttemptById(id);
		testAttempt.complete(points);
		testAttempt.validate();
		testAttemptRepository.update(testAttempt);
	}

	public void attachCertificate(UUID id, UUID certificateId) {
		TestAttemptModel testAttempt = getTestAttemptById(id);
		testAttempt.validate();
		testAttemptRepository.update(testAttempt);
	}

	public void deleteTestAttempt(UUID id) {
		testAttemptRepository.deleteById(id);
	}

	// public Optional<TestAttemptModel> getBestAttempt(UUID studentId, UUID testId)
	// {
	// return testAttemptRepository.findBestAttemptByStudentAndTest(studentId,
	// testId);
	// }

	// Новые методы для поддержки контроллера:

	public List<TestAttemptModel> getAllTestAttempts() {
		return testAttemptRepository.findAll();
	}

	public Optional<TestAttemptModel> findById(UUID id) {
		return testAttemptRepository.findById(id);
	}

	public List<TestAttemptModel> getTestAttemptsByStudentId(UUID studentId) {
		return testAttemptRepository.findByStudentId(studentId);
	}

	public List<TestAttemptModel> getTestAttemptsByTestId(UUID testId) {
		return testAttemptRepository.findByTestId(testId);
	}

	public List<TestAttemptModel> getTestAttemptsByDate(LocalDate date) {
		return testAttemptRepository.findByDate(date);
	}

	public List<TestAttemptModel> getCompletedTestAttempts() {
		return testAttemptRepository.findCompletedAttempts();
	}

	public List<TestAttemptModel> getIncompleteTestAttempts() {
		return testAttemptRepository.findIncompleteAttempts();
	}

	public List<TestAttemptModel> getAttemptsByStudentAndTest(UUID studentId, UUID testId) {
		return testAttemptRepository.findByStudentAndTest(studentId, testId);
	}

	public int countAttemptsByStudentAndTest(UUID studentId, UUID testId) {
		return testAttemptRepository.countAttemptsByStudentAndTest(studentId, testId);
	}

	public boolean existsById(UUID id) {
		return testAttemptRepository.existsById(id);
	}

	// public Optional<TestAttemptModel> findBestAttemptByStudentAndTest(UUID
	// studentId, UUID testId) {
	// return testAttemptRepository.findBestAttemptByStudentAndTest(studentId,
	// testId);
	// }
}