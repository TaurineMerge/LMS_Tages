package com.example.lms.draft.domain.service;

import java.util.List;
import java.util.UUID;

import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.model.DraftModel;
import com.example.lms.draft.domain.repository.DraftRepositoryInterface;

public class DraftService {
    private final DraftRepositoryInterface repository;

    public DraftService(DraftRepositoryInterface repository) {
        this.repository = repository;
    }

    /**
     * Преобразует DTO в доменную модель {@link DraftModel}.
     * Убрана логика с временным courseId - используем null для черновиков без курса
     */
    private DraftModel toModel(Draft dto) {
        if (dto.getId() == null) {
            // Новый черновик
            return new DraftModel(
                dto.getTestId(),  // МОЖЕТ БЫТЬ NULL
                dto.getCourseId(), // МОЖЕТ БЫТЬ NULL - используем переданное значение
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription()
            );
        } else {
            // Существующий черновик
            return new DraftModel(
                dto.getId(),
                dto.getTestId(),  // МОЖЕТ БЫТЬ NULL
                dto.getCourseId(), // МОЖЕТ БЫТЬ NULL
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription()
            );
        }
    }

    /**
     * Преобразует доменную модель черновика в DTO для API.
     */
    private Draft toDto(DraftModel model) {
        Draft dto = new Draft();
        dto.setId(model.getId());
        dto.setTitle(model.getTitle());
        dto.setMin_point(model.getMinPoint());
        dto.setDescription(model.getDescription());
        dto.setTestId(model.getTestId());
        dto.setCourseId(model.getCourseId()); // Просто передаем как есть
        return dto;
    }

    /**
     * Возвращает список всех черновиков.
     */
    public List<Draft> getAllDrafts() {
        return repository.findAll().stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Создаёт новый черновик.
     */
    public Draft createDraft(Draft dto) {
        if (dto.getId() != null) {
            throw new IllegalArgumentException("При создании нового черновика ID должен быть null");
        }
        
        DraftModel model = toModel(dto);
        model.validate();
        DraftModel saved = repository.create(model);
        return toDto(saved);
    }

    /**
     * Получает черновик по ID.
     */
    public Draft getDraftById(UUID id) {
        if (id == null) {
            throw new IllegalArgumentException("ID черновика не может быть null");
        }
        
        DraftModel model = repository.findById(id)
            .orElseThrow(() -> new RuntimeException("Черновик с ID " + id + " не найден"));
        return toDto(model);
    }

    /**
     * Получает черновик по testId.
     */
    public Draft getDraftByTestId(UUID testId) {
        if (testId == null) {
            return null;
        }
        
        return repository.findByTestId(testId)
            .map(this::toDto)
            .orElse(null);
    }
    
    /**
     * Получает черновики по courseId.
     * Если courseId null, возвращает черновики без привязки к курсу
     */
    public List<Draft> getDraftsByCourseId(UUID courseId) {
        return repository.findByCourseId(courseId).stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Проверяет, существует ли черновик для указанного теста.
     */
    public boolean existsByTestId(UUID testId) {
        if (testId == null) {
            return false;
        }
        
        return repository.findByTestId(testId).isPresent();
    }

    /**
     * Обновляет существующий черновик.
     */
    public Draft updateDraft(Draft dto) {
        if (dto.getId() == null) {
            throw new IllegalArgumentException("ID черновика обязателен для обновления");
        }
        
        DraftModel model = toModel(dto);
        model.validate();
        DraftModel updated = repository.update(model);
        return toDto(updated);
    }

    /**
     * Удаляет черновик по ID.
     */
    public boolean deleteDraft(UUID id) {
        if (id == null) {
            throw new IllegalArgumentException("ID черновика не может быть null");
        }
        
        return repository.deleteById(id);
    }

    /**
     * Удаляет черновики по courseId.
     */
    public boolean deleteDraftsByCourseId(UUID courseId) {
        return repository.deleteByCourseId(courseId);
    }
}