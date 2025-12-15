FROM maven:3.9-eclipse-temurin-17 AS build
WORKDIR /app
COPY . .
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jre
WORKDIR /app

# Скачиваем OpenTelemetry Java Agent версии 1.31.0
ADD https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v1.31.0/opentelemetry-javaagent.jar /opentelemetry-javaagent.jar

# Копируем собранный JAR
COPY --from=build /app/target/app.jar app.jar

# Копируем конфигурационный файл OpenTelemetry
COPY opentelemetry.properties .

EXPOSE 8085

# Запускаем с Java Agent
ENTRYPOINT ["java", \
    "-javaagent:/opentelemetry-javaagent.jar", \
    "-Dotel.javaagent.configuration-file=opentelemetry.properties", \
    "-jar", "app.jar"]