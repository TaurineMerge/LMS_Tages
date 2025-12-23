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
        return new TestModel(
                dto.getId() != null ? UUID.fromString(dto.getId().toString()) : null,
                dto.getCourseId() != null ? UUID.fromString(dto.getCourseId().toString()) : null,
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription());
    }

    // ---------------------------------------------------------------------
    // MODEL → DTO
    // ---------------------------------------------------------------------

    /**
     * Преобразует доменную модель теста в DTO для API.
     *
     * @param model доменная модель теста
     * @return DTO, отправляемый наружу
     */
    private Test toDto(TestModel model) {
        return new Test(
                model.getId() != null ? model.getId().toString() : null,
                model.getCourseId() != null ? model.getCourseId().toString() : null,
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
        UUID uuid = UUID.fromString(id);
        return repository.deleteById(uuid);
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
}