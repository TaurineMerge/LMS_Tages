package com.example.lms.config;

/**
 * Конфигурационный класс для параметров подключения к базе данных.
 * <p>
 * Хранит:
 * <ul>
 *     <li>URL подключения</li>
 *     <li>имя пользователя</li>
 *     <li>пароль</li>
 * </ul>
 * Используется сервисами и репозиториями для создания соединений с БД.
 */
public class DatabaseConfig {

    /**
     * Строка подключения к базе данных (JDBC URL).
     * Например: jdbc:postgresql://localhost:5432/lms
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
     * @param url      строка подключения JDBC
     * @param user     имя пользователя
     * @param password пароль пользователя
     */
    public DatabaseConfig(String url, String user, String password) {
        this.url = url;
        this.user = user;
        this.password = password;
    }

    /**
     * @return строку подключения (JDBC URL)
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