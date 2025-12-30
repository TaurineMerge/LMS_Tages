package com.example.lms.test.domain.service;

import java.util.List;
import java.util.UUID;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.model.TestModel;
import com.example.lms.test.domain.repository.TestRepositoryInterface;

/**
 * Сервис для работы с тестами.
 * <p>
 * Отвечает за:
 * <ul>
 * <li>конвертацию DTO ↔ Domain Model</li>
 * <li>вызовы репозитория ({@link TestRepositoryInterface})</li>
 * <li>бизнес-логику CRUD-операций над тестами</li>
 * </ul>
 *
 * Сервисный слой отделяет контроллеры от репозитория
 * и обеспечивает единое место для обработки бизнес-процессов.
 */
public class TestService {

    /** Репозиторий тестов (слой работы с базой данных). */
    private final TestRepositoryInterface repository;
    
    /** Временный UUID для тестов без курса. */
    private static final UUID TEMP_COURSE_ID = 
        UUID.fromString("11111111-1111-1111-1111-111111111111");

    /**
     * Создаёт сервис тестов.
     *
     * @param repository репозиторий, выполняющий операции с БД
     */
    public TestService(TestRepositoryInterface repository) {
        this.repository = repository;
    }

    // ---------------------------------------------------------------------
    // DTO -> MODEL
    // ---------------------------------------------------------------------

    /**
     * Преобразует DTO в доменную модель {@link TestModel}.
     * <p>
     * Доменная модель используется на сервисном и репозиторном уровнях.
     * DTO применяется только для API и не содержит бизнес-логики.
     *
     * @param dto объект API DTO
     * @return доменная модель теста
     */
    private TestModel toModel(Test dto) {
        UUID courseUuid = TEMP_COURSE_ID; // Дефолтное значение
        
        // Если courseId передан и валиден, используем его
        if (dto.getCourseId() != null && !dto.getCourseId().trim().isEmpty()) {
            try {
                courseUuid = UUID.fromString(dto.getCourseId());
            } catch (IllegalArgumentException e) {
                // Если невалидный UUID, используем временный
                // Можно добавить логгирование: logger.warn("Invalid courseId: {}, using temp UUID", dto.getCourseId());
                courseUuid = TEMP_COURSE_ID;
            }
        }
        
        return new TestModel(
                dto.getId() != null ? UUID.fromString(dto.getId()) : null,
                courseUuid, // Всегда будет значение (TEMP_COURSE_ID или переданный)
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription());
    }

    // ---------------------------------------------------------------------
    // MODEL → DTO
    // ---------------------------------------------------------------------

    /**
     * Преобразует доменную модель теста в DTO для API.
     * <p>
     * Если courseId равен временному UUID, возвращаем null в DTO
     * для обратной совместимости с API.
     *
     * @param model доменная модель теста
     * @return DTO, отправляемый наружу
     */
    private Test toDto(TestModel model) {
        // Если courseId равен временному UUID, возвращаем null
        String courseId = null;
        if (model.getCourseId() != null && !model.getCourseId().equals(TEMP_COURSE_ID)) {
            courseId = model.getCourseId().toString();
        }
        
        return new Test(
                model.getId() != null ? model.getId().toString() : null,
                courseId, // null для временного курса
                model.getTitle(),
                model.getMinPoint(),
                model.getDescription());
    }

    // ---------------------------------------------------------------------
    // PUBLIC API METHODS
    // ---------------------------------------------------------------------

    /**
     * Возвращает список всех тестов, конвертируя их в DTO.
     *
     * @return список тестов в формате DTO
     */
    public List<Test> getAllTests() {
        return repository.findAll().stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Создаёт новый тест.
     *
     * @param dto входные данные теста
     * @return созданный тест в виде DTO
     */
    public Test createTest(Test dto) {
        TestModel model = toModel(dto);
        TestModel saved = repository.save(model);
        return toDto(saved);
    }

    /**
     * Получает тест по ID.
     * <p>
     * Если тест не найден — будет выброшено
     * {@link java.util.NoSuchElementException}.
     *
     * @param id строковый UUID теста
     * @return DTO найденного теста
     */
    public Test getTestById(String id) {
        UUID uuid = UUID.fromString(id);
        TestModel model = repository.findById(uuid).orElseThrow();
        return toDto(model);
    }

    /**
     * Обновляет существующий тест.
     *
     * @param dto данные теста с актуальными полями
     * @return DTO обновлённого теста
     */
    public Test updateTest(Test dto) {
        TestModel model = toModel(dto);
        TestModel updated = repository.update(model);
        return toDto(updated);
    }

    /**
     * Удаляет тест по ID.
     *
     * @param id строковый UUID теста
     * @return true — если тест был удалён; false — если не найден
     */
    public boolean deleteTest(String id) {
    try {
        UUID uuid = UUID.fromString(id);
        
        // Проверяем, существует ли тест
        if (repository.findById(uuid).isEmpty()) {
            return false; // Тест не найден
        }
        
        // Пытаемся удалить (сначала вопросы, потом тест)
        return repository.deleteById(uuid);
        
    } catch (IllegalArgumentException e) {
        throw new RuntimeException("Неверный формат ID теста: " + id, e);
    } catch (RuntimeException e) {
        // Добавляем информацию об ID в сообщение об ошибке
        throw new RuntimeException("Ошибка при удалении теста " + id + ": " + e.getMessage(), e);
    }
}

    /**
     * Проверяет, существует ли тест по ID курса.
     *
     * @param course_id строковый UUID курса
     * @return true — если тест был существует; false — если не найден
     */
    public boolean existsByCourseId(String course_id) {
        UUID uuid = UUID.fromString(course_id);
        return repository.existsByCourseId(uuid);
    }

    /**
     * Получает все тесты по ID курса.
     *
     * @param courseId строковый UUID курса
     * @return список тестов в формате DTO
     */
    public List<Test> getTestsByCourseId(String courseId) {
        UUID uuid = UUID.fromString(courseId);
        List<TestModel> models = repository.findByCourseId(uuid);
        return models.stream()
                .map(this::toDto)
                .toList();
    }
    
    // ---------------------------------------------------------------------
    // ДОПОЛНИТЕЛЬНЫЕ МЕТОДЫ (опционально)
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
     * Возвращает временный UUID для тестов без курса.
     *
     * @return временный UUID
     */
    public static UUID getTempCourseId() {
        return TEMP_COURSE_ID;
    }
}