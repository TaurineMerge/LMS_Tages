# Вариант A: Используйте более общий тег
FROM maven:3.9-eclipse-temurin-17

# Вариант B: Или используйте openjdk
# FROM openjdk:17-jdk-slim

WORKDIR /app
COPY . .
EXPOSE 8085
CMD ["mvn", "clean", "compile", "exec:java", "-Dexec.mainClass=com.example.lms.Main"]