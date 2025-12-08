package com.example.lms.config;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class DatabaseConfig {

    public static void main(String[] args) {
        String url = "jdbc:postgresql://localhost:5432/appdb";
        String user = "appuser";
        String password = "password";
        Connection connection = null;
        try {
            connection = DriverManager.getConnection(url, user, password);
            System.out.println("Соединение установлено!");

        } catch (SQLException e) {
            System.err.println("Ошибка при подключении к базе данных: " + e.getMessage());
        } finally {
            if (connection != null) {
                try {
                    connection.close();
                } catch (SQLException e) {
                    System.err.println("Ошибка при закрытии соединения: " + e.getMessage());
                }
            }
        }
    }
}
