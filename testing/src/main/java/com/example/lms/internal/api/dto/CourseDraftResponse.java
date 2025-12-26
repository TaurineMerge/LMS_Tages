package com.example.lms.internal.api.dto;

import java.io.Serializable;
import java.util.UUID;

/**
 * Data Transfer Object (DTO) для ответа на запрос получения черновика теста по идентификатору курса.
 * <p>
 * Используется во внутреннем API для передачи данных о черновике теста, связанного с конкретным курсом.
 * Класс реализует интерфейс {@link Serializable} для поддержки сериализации объектов.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see Serializable
 */
public class CourseDraftResponse implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * Данные черновика теста.
     * Может быть {@code null}, если черновик не найден.
     */
    private DraftData data;

    /**
     * Статус выполнения запроса.
     * Возможные значения:
     * <ul>
     *   <li>"success" - черновик успешно найден</li>
     *   <li>"not_found" - черновик не найден для указанного курса</li>
     *   <li>"error" - произошла ошибка при обработке запроса</li>
     * </ul>
     */
    private String status;

    /**
     * Создает пустой экземпляр CourseDraftResponse.
     * <p>
     * Используется фреймворками для десериализации JSON.
     */
    public CourseDraftResponse() {
    }

    /**
     * Создает полностью инициализированный экземпляр CourseDraftResponse.
     *
     * @param data данные черновика теста (может быть {@code null})
     * @param status статус выполнения запроса (не может быть {@code null} или пустым)
     * @throws IllegalArgumentException если {@code status} равен {@code null} или пустой строке
     */
    public CourseDraftResponse(DraftData data, String status) {
        this.data = data;
        this.status = status;
    }

    /**
     * Возвращает данные черновика теста.
     *
     * @return объект {@link DraftData} с данными черновика или {@code null}, если черновик не найден
     */
    public DraftData getData() {
        return data;
    }

    /**
     * Устанавливает данные черновика теста.
     *
     * @param data объект {@link DraftData} с данными черновика (может быть {@code null})
     */
    public void setData(DraftData data) {
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
     * Data Transfer Object (DTO) для передачи данных черновика теста.
     * <p>
     * Содержит основные атрибуты черновика теста, связанного с курсом.
     * Класс реализует интерфейс {@link Serializable} для поддержки сериализации объектов.
     *
     * @author Команда разработки LMS
     * @version 1.0
     */
    public static class DraftData implements Serializable {

        private static final long serialVersionUID = 1L;

        /**
         * Уникальный идентификатор черновика теста.
         */
        private UUID id;

        /**
         * Уникальный идентификатор теста, связанного с черновиком.
         * Может быть {@code null}, если черновик еще не связан с тестом.
         */
        private UUID testId;

        /**
         * Уникальный идентификатор курса, для которого создан черновик теста.
         */
        private UUID courseId;

        /**
         * Название черновика теста.
         * Пример: "Промежуточный тест по Java Core"
         */
        private String title;

        /**
         * Минимальное количество баллов для успешного прохождения теста.
         * Может быть {@code null}, если порог не установлен.
         */
        private Integer min_point;

        /**
         * Описание черновика теста.
         * Может содержать информацию о целях теста, темах и требованиях.
         */
        private String description;

        /**
         * Создает пустой экземпляр DraftData.
         * <p>
         * Используется фреймворками для десериализации JSON.
         */
        public DraftData() {
        }

        /**
         * Создает полностью инициализированный экземпляр DraftData.
         *
         * @param id уникальный идентификатор черновика теста
         * @param testId уникальный идентификатор связанного теста (может быть {@code null})
         * @param courseId уникальный идентификатор курса
         * @param title название черновика теста
         * @param min_point минимальный порог баллов (может быть {@code null})
         * @param description описание черновика теста
         */
        public DraftData(UUID id, UUID testId, UUID courseId, String title, Integer min_point, String description) {
            this.id = id;
            this.testId = testId;
            this.courseId = courseId;
            this.title = title;
            this.min_point = min_point;
            this.description = description;
        }

        /**
         * Возвращает уникальный идентификатор черновика теста.
         *
         * @return идентификатор черновика
         */
        public UUID getId() {
            return id;
        }

        /**
         * Устанавливает уникальный идентификатор черновика теста.
         *
         * @param id идентификатор черновика
         */
        public void setId(UUID id) {
            this.id = id;
        }

        /**
         * Возвращает уникальный идентификатор связанного теста.
         *
         * @return идентификатор теста или {@code null}, если тест не связан
         */
        public UUID getTestId() {
            return testId;
        }

        /**
         * Устанавливает уникальный идентификатор связанного теста.
         *
         * @param testId идентификатор теста (может быть {@code null})
         */
        public void setTestId(UUID testId) {
            this.testId = testId;
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
         * Возвращает название черновика теста.
         *
         * @return название черновика
         */
        public String getTitle() {
            return title;
        }

        /**
         * Устанавливает название черновика теста.
         *
         * @param title название черновика
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
         * Возвращает описание черновика теста.
         *
         * @return описание черновика или {@code null}, если описание отсутствует
         */
        public String getDescription() {
            return description;
        }

        /**
         * Устанавливает описание черновика теста.
         *
         * @param description описание черновика (может быть {@code null})
         */
        public void setDescription(String description) {
            this.description = description;
        }
    }
}