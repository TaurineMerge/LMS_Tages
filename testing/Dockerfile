FROM maven:3.9-eclipse-temurin-17 AS builder
WORKDIR /app

# кэш зависимостей
COPY pom.xml .
RUN mvn dependency:go-offline

# исходники и ресурсы
COPY src ./src
COPY opentelemetry.properties ./opentelemetry.properties
COPY .env ./.env

# сборка без тестов и javadoc
RUN mvn clean package -DskipTests -Dmaven.javadoc.skip=true

FROM eclipse-temurin:17-jre
WORKDIR /app

# нужен curl для скачивания агента
RUN apt-get update && apt-get install -y curl ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/app.jar ./app.jar
COPY --from=builder /app/opentelemetry.properties ./opentelemetry.properties
COPY --from=builder /app/.env ./.env

# скачать Otel Java Agent
RUN curl -L -o /app/opentelemetry-javaagent.jar \
  https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v1.32.0/opentelemetry-javaagent.jar

ENV JAVA_TOOL_OPTIONS="-javaagent:/app/opentelemetry-javaagent.jar"
ENV OTEL_CONFIG_FILE=/app/opentelemetry.properties
ENV OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
ENV OTEL_METRICS_EXPORTER=none

CMD ["java", "-jar", "/app/app.jar"]