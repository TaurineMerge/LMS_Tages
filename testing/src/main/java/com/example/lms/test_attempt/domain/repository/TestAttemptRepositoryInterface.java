package com.example.lms.test_attempt.domain.repository;

import com.example.lms.test_attempt.domain.model.TestAttemptModel;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с попытками тестов.
 *
 * Соответствует таблице test_attempt_b:
 *  - id                UUID    (PK, not null)
 *  - student_id        UUID    (not null)
 *  - test_id           UUID    (FK → test_d.id, not null)
 *  - date_of_attempt   DATE    (может быть null)
 *  - point             INTEGER (может быть null)
 *  - certificate_id    UUID    (может быть null)
 *  - attempt_version   JSON    (может быть null)
 *  - attempt_snapshot  VARCHAR(256) (может быть null)
 *  - completed         BOOLEAN (может быть null)
 *
 * Отвечает за доступ к данным:
 *  - создание / обновление / удаление попыток тестов;
 *  - выборка по идентификатору, студенту и тесту;
 *  - подсчёт попыток;
 *  - поиск завершенных попыток.
 */
public interface TestAttemptRepositoryInterface {

    /**
     * Сохранить новую попытку теста.
     *
     * @param testAttempt доменная модель попытки теста
     * @return сохранённая модель с присвоенным идентификатором
     */
    TestAttemptModel save(TestAttemptModel testAttempt);

    /**
     * Обновить существующую попытку теста.
     *
     * @param testAttempt модель попытки теста с актуальными данными
     * @return обновлённая модель
     */
    TestAttemptModel update(TestAttemptModel testAttempt);

    /**
     * Найти попытку теста по её идентификатору.
     *
     * @param id идентификатор попытки (test_attempt_b.id)
     * @return Optional с попыткой, если найдена
     */
    Optional<TestAttemptModel> findById(UUID id);

    /**
     * Получить список всех попыток тестов.
     *
     * @return список всех попыток
     */
    List<TestAttemptModel> findAll();

    /**
     * Удалить попытку теста по её идентификатору.
     *
     * @param id идентификатор попытки
     * @return true, если попытка была удалена; false — если запись не найдена
     */
    boolean deleteById(UUID id);

    /**
     * Проверить существование попытки теста по её идентификатору.
     *
     * @param id идентификатор попытки
     * @return true, если попытка существует
     */
    boolean existsById(UUID id);

    /**
     * Найти все попытки теста для указанного студента.
     *
     * @param studentId идентификатор студента (test_attempt_b.student_id)
     * @return список попыток данного студента
     */
    List<TestAttemptModel> findByStudentId(UUID studentId);

    /**
     * Найти все попытки для указанного теста.
     *
     * @param testId идентификатор теста (test_attempt_b.test_id)
     * @return список попыток данного теста
     */
    List<TestAttemptModel> findByTestId(UUID testId);

    /**
     * Найти все попытки для указанного студента и теста.
     *
     * @param studentId идентификатор студента
     * @param testId идентификатор теста
     * @return список попыток
     */
    List<TestAttemptModel> findByStudentIdAndTestId(UUID studentId, UUID testId);

    /**
     * Получить количество попыток для студента.
     *
     * @param studentId идентификатор студента
     * @return число попыток, относящихся к студенту
     */
    int countByStudentId(UUID studentId);

    /**
     * Получить количество попыток для теста.
     *
     * @param testId идентификатор теста
     * @return число попыток, относящихся к тесту
     */
    int countByTestId(UUID testId);

    /**
     * Найти все завершенные попытки тестов.
     *
     * @return список завершенных попыток
     */
    List<TestAttemptModel> findCompletedAttempts();
}