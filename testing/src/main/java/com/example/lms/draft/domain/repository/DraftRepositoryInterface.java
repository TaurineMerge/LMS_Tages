package com.example.lms.draft.domain.repository;

import com.example.lms.draft.domain.model.DraftModel;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Контракт репозитория для работы с черновиками тестов (draft_b).
 * <p>
 * Интерфейс находится в domain-слое, чтобы сервисы не зависели от конкретной
 * реализации хранения данных (JDBC/JPA/внешний сервис и т.д.).
 */
public interface DraftRepositoryInterface {

    /**
     * Создать новый черновик.
     *
     * @param draftModel доменная модель черновика
     * @return созданный черновик (обычно с проставленным id)
     */
    DraftModel create(DraftModel draftModel);

    /**
     * Найти черновик по id.
     *
     * @param id идентификатор черновика
     * @return optional с черновиком, если найден
     */
    Optional<DraftModel> findById(UUID id);

    /**
     * Найти черновик по testId.
     * <p>
     * У тебя в таблице есть test_id not null — часто логично, что
     * на один тест максимум один черновик. Но интерфейс оставим гибким:
     * можно вернуть Optional, если гарантия "1:1".
     *
     * @param testId идентификатор теста
     * @return optional с черновиком, если найден
     */
    Optional<DraftModel> findByTestId(UUID testId);

    /**
     * Получить список всех черновиков.
     *
     * @return список черновиков
     */
    List<DraftModel> findAll();

    /**
     * Обновить существующий черновик.
     *
     * @param draftModel доменная модель с заполненным id
     * @return обновлённый черновик
     */
    DraftModel update(DraftModel draftModel);

    /**
     * Удалить черновик по id.
     *
     * @param id идентификатор черновика
     * @return true если удалено, false если записи не было
     */
    boolean deleteById(UUID id);
}