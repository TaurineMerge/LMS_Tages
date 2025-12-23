package com.example.lms.draft.domain.service;

import java.util.List;
import java.util.UUID;

import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.model.DraftModel;
import com.example.lms.draft.domain.repository.DraftRepositoryInterface;

public class DraftService {
    private final DraftRepositoryInterface repository;
    
    /** Временный UUID для черновиков без курса. */
    private static final UUID TEMP_COURSE_ID = 
        UUID.fromString("11111111-1111-1111-1111-111111111111");

    public DraftService(DraftRepositoryInterface repository) {
        this.repository = repository;
    }

    /**
     * Преобразует DTO в доменную модель {@link DraftModel}.
     */
    private DraftModel toModel(Draft dto) {
        UUID courseUuid = TEMP_COURSE_ID; // Дефолтное значение
        
        // Если courseId передан и валиден, используем его
        if (dto.getCourseId() != null && !dto.getCourseId().equals(TEMP_COURSE_ID)) {
            // У Draft DTO courseId уже UUID, не нужно конвертировать
            courseUuid = dto.getCourseId();
        }
        
        if (dto.getId() == null) {
            // Новый черновик
            return new DraftModel(
                dto.getTestId(),  // МОЖЕТ БЫТЬ NULL
                courseUuid,       // course_id (temp или реальный)
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription()
            );
        } else {
            // Существующий черновик
            return new DraftModel(
                dto.getId(),
                dto.getTestId(),  // МОЖЕТ БЫТЬ NULL
                courseUuid,       // course_id (temp или реальный)
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
        dto.setTestId(model.getTestId());  // МОЖЕТ БЫТЬ NULL
        
        // Если courseId равен временному UUID, возвращаем null
        if (model.getCourseId() != null && !model.getCourseId().equals(TEMP_COURSE_ID)) {
            dto.setCourseId(model.getCourseId());
        } else {
            dto.setCourseId(null);
        }
        
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
        model.validate(); // Важно: вызываем валидацию
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
     * Возвращает null если черновик не найден (вместо исключения).
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
     */
    public List<Draft> getDraftsByCourseId(UUID courseId) {
        if (courseId == null) {
            // Возвращаем черновики с временным courseId
            return repository.findByCourseId(TEMP_COURSE_ID).stream()
                    .map(this::toDto)
                    .toList();
        }
        
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
     * Удаляет черновик по testId.
     */
    public boolean deleteDraftByTestId(UUID testId) {
        if (testId == null) {
            return false;
        }
        
        return repository.findByTestId(testId)
            .map(draft -> repository.deleteById(draft.getId()))
            .orElse(false);
    }
    
    /**
     * Удаляет черновики по courseId.
     */
    public boolean deleteDraftsByCourseId(UUID courseId) {
        if (courseId == null) {
            return repository.deleteByCourseId(TEMP_COURSE_ID);
        }
        
        return repository.deleteByCourseId(courseId);
    }
    
    // ---------------------------------------------------------------------
    // ДОПОЛНИТЕЛЬНЫЕ МЕТОДЫ
    // ---------------------------------------------------------------------
    
    /**
     * Проверяет, является ли UUID временным courseId.
     *
     * @param courseId UUID для проверки
     * @return true если это временный UUID, false если реальный курс
     */
    public boolean isTempCourseId(UUID courseId) {
        return TEMP_COURSE_ID.equals(courseId);
    }
    
    /**
     * Возвращает временный UUID для черновиков без курса.
     *
     * @return временный UUID
     */
    public static UUID getTempCourseId() {
        return TEMP_COURSE_ID;
    }
}