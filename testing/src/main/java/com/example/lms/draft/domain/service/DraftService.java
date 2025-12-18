package com.example.lms.draft.domain.service;

import java.util.List;
import java.util.UUID;

import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.model.DraftModel;
import com.example.lms.draft.domain.repository.DraftRepositoryInterface;

/**
 * Сервис для работы с черновиками тестов (draft).
 * <p>
 * Отвечает за:
 * <ul>
 *   <li>конвертацию DTO ↔ Domain Model</li>
 *   <li>вызовы репозитория ({@link DraftRepositoryInterface})</li>
 *   <li>бизнес-логику CRUD-операций над черновиками</li>
 * </ul>
 *
 * Сервисный слой отделяет контроллеры от репозитория
 * и обеспечивает единое место для обработки бизнес-процессов.
 */
public class DraftService {

    /** Репозиторий черновиков (слой работы с базой данных). */
    private final DraftRepositoryInterface repository;

    /**
     * Создаёт сервис черновиков.
     *
     * @param repository репозиторий, выполняющий операции с БД
     */
    public DraftService(DraftRepositoryInterface repository) {
        this.repository = repository;
    }

    // ---------------------------------------------------------------------
    // DTO -> MODEL
    // ---------------------------------------------------------------------

    /**
     * Преобразует DTO в доменную модель {@link DraftModel}.
     *
     * @param dto объект API DTO
     * @return доменная модель черновика
     */
    private DraftModel toModel(Draft dto) {
        return new DraftModel(
                dto.getId() != null ? UUID.fromString(dto.getId().toString()) : null,
                dto.getTest_id() != null ? UUID.fromString(dto.getTest_id().toString()) : null,
                dto.getTitle(),
                dto.getMin_point(),
                dto.getDescription()
        );
    }

    // ---------------------------------------------------------------------
    // MODEL → DTO
    // ---------------------------------------------------------------------

    /**
     * Преобразует доменную модель черновика в DTO для API.
     *
     * @param model доменная модель
     * @return DTO, отправляемый наружу
     */
    private Draft toDto(DraftModel model) {
        return new Draft(
                model.getId() != null ? model.getId().toString() : null,
                model.getTitle(),
                model.getMinPoint(),
                model.getDescription(),
                model.getTestId() != null ? model.getTestId().toString() : null
        );
    }

    // ---------------------------------------------------------------------
    // PUBLIC API METHODS
    // ---------------------------------------------------------------------

    /**
     * Возвращает список всех черновиков, конвертируя их в DTO.
     *
     * @return список черновиков в формате DTO
     */
    public List<Draft> getAllDrafts() {
        return repository.findAll().stream()
                .map(this::toDto)
                .toList();
    }

    /**
     * Создаёт новый черновик.
     *
     * @param dto входные данные черновика
     * @return созданный черновик в виде DTO
     */
    public Draft createDraft(Draft dto) {
        DraftModel model = toModel(dto);
        DraftModel saved = repository.create(model);
        return toDto(saved);
    }

    /**
     * Получает черновик по ID.
     * <p>
     * Если черновик не найден — будет выброшено
     * {@link java.util.NoSuchElementException}.
     *
     * @param id строковый UUID черновика
     * @return DTO найденного черновика
     */
    public Draft getDraftById(String id) {
        UUID uuid = UUID.fromString(id);
        DraftModel model = repository.findById(uuid).orElseThrow();
        return toDto(model);
    }

    /**
     * Получает черновик по testId.
     * <p>
     * Частый сценарий: найти черновик для конкретного теста.
     *
     * @param testId строковый UUID теста
     * @return DTO найденного черновика
     */
    public Draft getDraftByTestId(String testId) {
        UUID uuid = UUID.fromString(testId);
        DraftModel model = repository.findByTestId(uuid).orElseThrow();
        return toDto(model);
    }

    /**
     * Обновляет существующий черновик.
     *
     * @param dto данные черновика с актуальными полями
     * @return DTO обновлённого черновика
     */
    public Draft updateDraft(Draft dto) {
        DraftModel model = toModel(dto);
        DraftModel updated = repository.update(model);
        return toDto(updated);
    }

    /**
     * Удаляет черновик по ID.
     *
     * @param id строковый UUID черновика
     * @return true — если удалён; false — если не найден
     */
    public boolean deleteDraft(String id) {
        UUID uuid = UUID.fromString(id);
        return repository.deleteById(uuid);
    }
}