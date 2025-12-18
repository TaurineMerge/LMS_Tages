package com.example.lms.draft.api.dto;

import java.io.Serializable;

/**
 * DTO для представления черновика теста (draft).
 * <p>
 * Используется для передачи данных о черновике между слоями приложения.
 * Содержит информацию:
 * <ul>
 * <li><b>id</b> — идентификатор черновика</li>
 * <li><b>title</b> — название черновика (обычно название теста)</li>
 * <li><b>min_point</b> — минимальный проходной балл</li>
 * <li><b>description</b> — описание</li>
 * <li><b>test_id</b> — идентификатор теста, к которому относится черновик</li>
 * </ul>
 *
 * Реализует {@link Serializable}, что позволяет передавать объект
 * через сеть или сохранять в файлы, если это необходимо.
 */
public class Draft implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * Идентификатор черновика.
     */
    private String id;

    /**
     * Название черновика.
     */
    private String title;

    /**
     * Минимальный балл, необходимый для прохождения теста.
     */
    private Integer min_point;

    /**
     * Текстовое описание черновика.
     */
    private String description;

    /**
     * Идентификатор теста, к которому привязан черновик.
     */
    private String test_id;

    /**
     * Пустой конструктор для сериализации/десериализации.
     */
    public Draft() {
    }

    /**
     * Основной конструктор.
     *
     * @param id          идентификатор черновика
     * @param title       название
     * @param min_point   минимальный проходной балл
     * @param description описание
     * @param test_id     идентификатор теста
     */
    public Draft(String id, String title, Integer min_point, String description, String test_id) {
        this.id = id;
        this.title = title;
        this.min_point = min_point;
        this.description = description;
        this.test_id = test_id;
    }

    /** @return идентификатор черновика */
    public String getId() {
        return id;
    }

    /** @param id новый идентификатор черновика */
    public void setId(String id) {
        this.id = id;
    }

    /** @return название черновика */
    public String getTitle() {
        return title;
    }

    /** @param title новое название черновика */
    public void setTitle(String title) {
        this.title = title;
    }

    /** @return минимальный проходной балл */
    public Integer getMin_point() {
        return min_point;
    }

    /** @param min_point новый минимальный проходной балл */
    public void setMin_point(Integer min_point) {
        this.min_point = min_point;
    }

    /** @return описание черновика */
    public String getDescription() {
        return description;
    }

    /** @param description новое описание черновика */
    public void setDescription(String description) {
        this.description = description;
    }

    /** @return идентификатор теста */
    public String getTest_id() {
        return test_id;
    }

    /** @param test_id новый идентификатор теста */
    public void setTest_id(String test_id) {
        this.test_id = test_id;
    }

    @Override
    public String toString() {
        return "Draft{" +
                "id=" + id +
                ", title='" + title + '\'' +
                ", min_point=" + min_point +
                ", description='" + description + '\'' +
                ", test_id=" + test_id +
                '}';
    }
}
