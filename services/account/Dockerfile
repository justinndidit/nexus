# STAGE 1: Build the JAR
FROM gradle:9.2-jdk25 AS builder

WORKDIR /app

# Copy gradle files first to cache dependencies
COPY settings.gradle build.gradle ./
COPY services/account/build.gradle ./services/account/

# Copy source code
COPY services/account/src ./services/account/src

# Build the fat JAR, skipping tests for speed (run tests in CI instead)
RUN gradle :services:account:bootJar -x test

# STAGE 2: Runtime
FROM eclipse-temurin:25-jre-alpine

WORKDIR /app

# Create a dedicated group and user
RUN addgroup -S nexus && adduser -S nexus -G nexus

COPY --from=builder /app/services/account/build/libs/*.jar app.jar

USER nexus

# Tune JVM for container environments
ENTRYPOINT ["java", "-XX:+UseContainerSupport", "-XX:MaxRAMPercentage=75.0", "-jar", "app.jar"]