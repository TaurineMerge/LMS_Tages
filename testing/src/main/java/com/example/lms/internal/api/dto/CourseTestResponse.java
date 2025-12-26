package com.example.lms.internal.api.dto;

import java.io.Serializable;
import java.util.UUID;

/**
 * Data Transfer Object (DTO) для ответа на запрос получения теста по идентификатору курса.
 * <p>
 * Используется во внутреннем API для передачи данных о тесте, связанного с конкретным курсом.
 * Класс реализует интерфейс {@link Serializable} для поддержки сериализации объектов.
 * <p>
 * Структура ответа включает данные теста и статус выполнения запроса, что позволяет
 * единообразно обрабатывать как успешные ответы, так и случаи, когда тест не найден.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see Serializable
 * @see CourseDraftResponse
 */
public class CourseTestResponse implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * Данные теста.
     * Содержит детальную информацию о тесте, если он найден.
     * Может быть {@code null}, если тест не найден для указанного курса.
     */
    private TestData data;

    /**
     * Статус выполнения запроса.
     * Возможные значения:
     * <ul>
     *   <li>"success" - тест успешно найден</li>
     *   <li>"not_found" - тест не найден для указанного курса</li>
     *   <li>"error" - произошла ошибка при обработке запроса</li>
     * </ul>
     */
    private String status;

    /**
     * Создает пустой экземпляр CourseTestResponse.
     * <p>
     * Используется фреймворками для десериализации JSON.
     */
    public CourseTestResponse() {
    }

    /**
     * Создает полностью инициализированный экземпляр CourseTestResponse.
     *
     * @param data данные теста (может быть {@code null})
     * @param status статус выполнения запроса (не может быть {@code null} или пустым)
     * @throws IllegalArgumentException если {@code status} равен {@code null} или пустой строке
     */
    public CourseTestResponse(TestData data, String status) {
        this.data = data;
        this.status = status;
    }

    /**
     * Возвращает данные теста.
     *
     * @return объект {@link TestData} с данными теста или {@code null}, если тест не найден
     */
    public TestData getData() {
        return data;
    }

    /**
     * Устанавливает данные теста.
     *
     * @param data объект {@link TestData} с данными теста (может быть {@code null})
     */
    public void setData(TestData data) {
        this.data = data;
    }

    /**
     * Возвращает статус выполнения запроса.
     *
     * @return строковый статус запроса
     */
    public String getStatus() {
        return status;
    }

    /**
     * Устанавливает статус выполнения запроса.
     *
     * @param status строковый статус запроса (не может быть {@code null} или пустым)
     * @throws IllegalArgumentException если {@code status} равен {@code null} или пустой строке
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Data Transfer Object (DTO) для передачи данных теста.
     * <p>
     * Содержит основные атрибуты теста, связанного с курсом.
     * Класс реализует интерфейс {@link Serializable} для поддержки сериализации объектов.
     *
     * @author Команда разработки LMS
     * @version 1.0
     */
    public static class TestData implements Serializable {

        private static final long serialVersionUID = 1L;

        /**
         * Уникальный идентификатор теста.
         */
        private UUID id;

        /**
         * Уникальный идентификатор курса, для которого предназначен тест.
         */
        private UUID courseId;

        /**
         * Название теста.
         * Пример: "Финальный экзамен по курсу Java Advanced"
         */
        private String title;

        /**
         * Минимальное количество баллов для успешного прохождения теста.
         * Если установлен, тест считается пройденным только при наборе баллов >= min_point.
         * Может быть {@code null}, если порог не установлен.
         */
        private Integer min_point;

        /**
         * Описание теста.
         * Может содержать информацию о целях теста, темах, длительности и других требованиях.
         */
        private String description;

        /**
         * Создает пустой экземпляр TestData.
         * <p>
         * Используется фреймворками для десериализации JSON.
         */
        public TestData() {
        }

        /**
         * Создает полностью инициализированный экземпляр TestData.
         *
         * @param id уникальный идентификатор теста
         * @param courseId уникальный идентификатор курса
         * @param title название теста
         * @param min_point минимальный порог баллов (может быть {@code null})
         * @param description описание теста
         */
        public TestData(UUID id, UUID courseId, String title, Integer min_point, String description) {
            this.id = id;
            this.courseId = courseId;
            this.title = title;
            this.min_point = min_point;
            this.description = description;
        }

        /**
         * Возвращает уникальный идентификатор теста.
         *
         * @return идентификатор теста
         */
        public UUID getId() {
            return id;
        }

        /**
         * Устанавливает уникальный идентификатор теста.
         *
         * @param id идентификатор теста
         */
        public void setId(UUID id) {
            this.id = id;
        }

        /**
         * Возвращает уникальный идентификатор курса.
         *
         * @return идентификатор курса
         */
        public UUID getCourseId() {
            return courseId;
        }

        /**
         * Устанавливает уникальный идентификатор курса.
         *
         * @param courseId идентификатор курса
         */
        public void setCourseId(UUID courseId) {
            this.courseId = courseId;
        }

        /**
         * Возвращает название теста.
         *
         * @return название теста
         */
        public String getTitle() {
            return title;
        }

        /**
         * Устанавливает название теста.
         *
         * @param title название теста
         * @throws IllegalArgumentException если {@code title} равен {@code null} или пустой строке
         */
        public void setTitle(String title) {
            this.title = title;
        }

        /**
         * Возвращает минимальное количество баллов для успешного прохождения.
         *
         * @return минимальный порог баллов или {@code null}, если порог не установлен
         */
        public Integer getMin_point() {
            return min_point;
        }

        /**
         * Устанавливает минимальное количество баллов для успешного прохождения.
         *
         * @param min_point минимальный порог баллов (может быть {@code null})
         */
        public void setMin_point(Integer min_point) {
            this.min_point = min_point;
        }

        /**
         * Возвращает описание теста.
         *
         * @return описание теста или {@code null}, если описание отсутствует
         */
        public String getDescription() {
            return description;
        }

        /**
         * Устанавливает описание теста.
         *
         * @param description описание теста (может быть {@code null})
         */
        public void setDescription(String description) {
            this.description = description;
        }
    }
}