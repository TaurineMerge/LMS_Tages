package com.example.lms.config;

/**
 * Конфигурационный класс для параметров подключения к базе данных.
 * <p>
 * Хранит основные параметры JDBC-подключения:
 * <ul>
 *     <li>{@code url} — строка подключения к базе данных</li>
 *     <li>{@code user} — имя пользователя базы данных</li>
 *     <li>{@code password} — пароль пользователя</li>
 * </ul>
 *
 * Класс является immutable-объектом: все поля финальные и задаются только в конструкторе.
 * Используется в репозиториях и сервисах, где требуется создание JDBC-соединений.
 */
public class DatabaseConfig {

    /**
     * JDBC URL базы данных.
     * <p>
     * Пример:
     * <pre>
     * jdbc:postgresql://localhost:5432/lms
     * </pre>
     */
    private final String url;

    /**
     * Имя пользователя базы данных.
     */
    private final String user;

    /**
     * Пароль пользователя базы данных.
     */
    private final String password;

    /**
     * Создаёт объект конфигурации подключения к базе данных.
     *
     * @param url      JDBC URL
     * @param user     имя пользователя
     * @param password пароль пользователя
     */
    public DatabaseConfig(String url, String user, String password) {
        this.url = url;
        this.user = user;
        this.password = password;
    }

    /**
     * @return JDBC URL базы данных
     */
    public String getUrl() {
        return url;
    }

    /**
     * @return имя пользователя базы данных
     */
    public String getUser() {
        return user;
    }

    /**
     * @return пароль пользователя базы данных
     */
    public String getPassword() {
        return password;
    }
}