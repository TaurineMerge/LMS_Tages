package com.example.lms.test.domain.repository;

import com.example.lms.test.domain.model.TestModel;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/**
 * Репозиторий для работы с тестами.
 *
 * Соответствует таблице test_d:
 *  - id          UUID    (PK, not null)
 *  - course_id   UUID    (FK → course_d.id, может быть null)
 *  - title       VARCHAR
 *  - min_point   INT
 *  - description TEXT
 *
 * Отвечает за доступ к данным:
 *  - создание / обновление / удаление тестов;
 *  - выборка по идентификатору и по курсу;
 *  - поиск по названию;
 *  - подсчёт тестов внутри курса.
 */
public interface TestRepositoryInterface {

    /**
     * Сохранить новый тест.
     *
     * @param test доменная модель теста (обычно без id до сохранения)
     * @return сохранённая модель с присвоенным идентификатором
     */
    TestModel save(TestModel test);

    /**
     * Обновить существующий тест.
     *
     * @param test модель теста с актуальными данными (id должен быть заполнен)
     * @return обновлённая модель
     */
    TestModel update(TestModel test);

    /**
     * Найти тест по его идентификатору.
     *
     * @param id идентификатор теста (TEST_D.id)
     * @return Optional с тестом, если найден
     */
    Optional<TestModel> findById(UUID id);

    /**
     * Получить список всех тестов.
     *
     * @return список всех тестов
     */
    List<TestModel> findAll();

    /**
     * Найти все тесты, привязанные к указанному курсу.
     *
     * @param courseId идентификатор курса (TEST_D.course_id)
     * @return список тестов данного курса
     */
    List<TestModel> findByCourseId(UUID courseId);

    /**
     * Удалить тест по его идентификатору.
     *
     * @param id идентификатор теста
     * @return true, если тест был удалён; false — если запись не найдена
     */
    boolean deleteById(UUID id);

    /**
     * Проверить существование теста по его идентификатору.
     *
     * @param id идентификатор теста
     * @return true, если тест существует
     */
    boolean existsById(UUID id);

    /**
     * Найти тесты по части названия (регистронезависимый поиск).
     *
     * Использует поле TEST_D.title.
     *
     * @param title часть названия для поиска
     * @return список найденных тестов
     */
    List<TestModel> findByTitleContaining(String title);

    /**
     * Получить количество тестов для курса.
     *
     * @param courseId идентификатор курса
     * @return число тестов, относящихся к курсу
     */
    int countByCourseId(UUID courseId);

    /**
     * Проверяет, существует ли тест по ID курса.
     *
     * @param course_id строковый UUID курса
     * @return true — если тест был существует; false — если не найден
     */
    boolean existsByCourseId(UUID course_id);
}