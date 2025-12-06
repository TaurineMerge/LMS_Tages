package com.example.lms.test.infrastructure.database;

import java.sql.*;

public class TestDb {

    public static void main(String[] args) {
        String url = "jdbc:postgresql://localhost:5432/appdb";
        String user = "appuser";
        String password = "password";

        try (Connection connection = DriverManager.getConnection(url, user, password);
                Statement statement = connection.createStatement();
                ResultSet resultSet = statement.executeQuery("SELECT title FROM tests")) {

            System.out.println("Employee Data:");

            // StringBuilder title = new StringBuilder();
            // while (titleSet.next()) {
            // String title = titleSet.getString("title");
            // result.append(resultSet.getString("title")).append("\n");
            // System.out.println(result.toString());
            // }

        } catch (SQLException e) {
            System.err.println("Error executing query: " + e.getMessage());
            e.printStackTrace();
        }
    }
}