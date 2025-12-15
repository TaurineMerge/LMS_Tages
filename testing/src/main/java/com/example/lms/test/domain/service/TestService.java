package com.example.lms.test.domain.service;

import java.util.List;
import java.util.UUID;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.model.TestModel;
import com.example.lms.test.domain.repository.TestRepositoryInterface;

public class TestService {
    private final TestRepositoryInterface repository;

    public TestService(TestRepositoryInterface repository) {
        this.repository = repository;
    }

    // DTO -> Model
    private TestModel toModel(Test dto) {
        return new TestModel(
                dto.getId() != null ? UUID.fromString(dto.getId().toString()) : null, // если у тебя UUID
                null, // например, courseId можно отдельно задавать
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription());
    }

    // Model -> DTO
    private Test toDto(TestModel model) {
        return new Test(
                model.getId() != null ? model.getId().getMostSignificantBits() : null,
                model.getTitle(),
                model.getMinPoint(),
                model.getDescription());
    }

    // Примеры методов
    public List<Test> getAllTests() {
        return repository.findAll().stream()
                .map(this::toDto)
                .toList();
    }

    public Test createTest(Test dto) {
        TestModel model = toModel(dto);
        TestModel saved = repository.save(model);
        return toDto(saved);
    }

    public Test getTestById(String id) {
        UUID uuid = UUID.fromString(id);
        TestModel model = repository.findById(uuid).orElseThrow();
        return toDto(model);
    }

    public Test updateTest(Test dto) {
        TestModel model = toModel(dto);
        TestModel updated = repository.update(model);
        return toDto(updated);
    }

    public boolean deleteTest(String id) {
        UUID uuid = UUID.fromString(id);
        return repository.deleteById(uuid);
    }
}
