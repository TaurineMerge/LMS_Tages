package com.example.lms.test_attempt.domain.service;

import java.time.LocalDate;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test_attempt.domain.model.Test_AttemptModel;
import com.example.lms.test_attempt.domain.repository.Test_AttemptRepositoryInterface;

public class TestAttemptService {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptService.class);

    private final Test_AttemptRepositoryInterface testAttemptRepository;

    public TestAttemptService(Test_AttemptRepositoryInterface testAttemptRepository) {
        this.testAttemptRepository = testAttemptRepository;
    }

    // Существующие методы...
    public Test_AttemptModel createTestAttempt(UUID studentId, UUID testId, String attemptVersion) {
        Test_AttemptModel testAttempt = new Test_AttemptModel(studentId, testId);
        testAttempt.validate();
        return testAttemptRepository.save(testAttempt);
    }

    public Test_AttemptModel getTestAttemptById(UUID id) {
        return testAttemptRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("No test attempt found with id: " + id));
    }

    public Test_AttemptModel updateTestAttempt(Test_AttemptModel testAttempt) {
        UUID id = testAttempt.getId();
        if (id == null) {
            throw new IllegalArgumentException("Test Attempt ID cannot be null for update");
        }

        Test_AttemptModel existingAttempt = testAttemptRepository.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("No test attempt found with id: " + id));

        
        existingAttempt.setPoint(testAttempt.getPoint());
        existingAttempt.validate();

        return testAttemptRepository.update(existingAttempt);
    }

    public void completeTestAttempt(UUID id, int points) {
        Test_AttemptModel testAttempt = getTestAttemptById(id);
        testAttempt.complete(points);
        testAttempt.validate();
        testAttemptRepository.update(testAttempt);
    }

    public void attachCertificate(UUID id, UUID certificateId) {
        Test_AttemptModel testAttempt = getTestAttemptById(id);
        testAttempt.validate();
        testAttemptRepository.update(testAttempt);
    }

    public void deleteTestAttempt(UUID id) {
        testAttemptRepository.deleteById(id);
    }

    public Optional<Test_AttemptModel> getBestAttempt(UUID studentId, UUID testId) {
        return testAttemptRepository.findBestAttemptByStudentAndTest(studentId, testId);
    }

    // Новые методы для поддержки контроллера:
    
    public List<Test_AttemptModel> getAllTestAttempts() {
        return testAttemptRepository.findAll();
    }
    
    public Optional<Test_AttemptModel> findById(UUID id) {
        return testAttemptRepository.findById(id);
    }
    
    public List<Test_AttemptModel> getTestAttemptsByStudentId(UUID studentId) {
        return testAttemptRepository.findByStudentId(studentId);
    }
    
    public List<Test_AttemptModel> getTestAttemptsByTestId(UUID testId) {
        return testAttemptRepository.findByTestId(testId);
    }
    
    public List<Test_AttemptModel> getTestAttemptsByDate(LocalDate date) {
        return testAttemptRepository.findByDate(date);
    }
    
    public List<Test_AttemptModel> getCompletedTestAttempts() {
        return testAttemptRepository.findCompletedAttempts();
    }
    
    public List<Test_AttemptModel> getIncompleteTestAttempts() {
        return testAttemptRepository.findIncompleteAttempts();
    }
    
    public List<Test_AttemptModel> getAttemptsByStudentAndTest(UUID studentId, UUID testId) {
        return testAttemptRepository.findByStudentAndTest(studentId, testId);
    }
    
    public int countAttemptsByStudentAndTest(UUID studentId, UUID testId) {
        return testAttemptRepository.countAttemptsByStudentAndTest(studentId, testId);
    }
    
    public boolean existsById(UUID id) {
        return testAttemptRepository.existsById(id);
    }
    
    public Optional<Test_AttemptModel> findBestAttemptByStudentAndTest(UUID studentId, UUID testId) {
        return testAttemptRepository.findBestAttemptByStudentAndTest(studentId, testId);
    }
}