FROM openjdk:21
COPY target/lazy-rabbit-reminder-1.0.0.jar /usr/src/reminder
WORKDIR /usr/src/reminder
EXPOSE 8080
CMD ["java" ,"-jar", "lazy-rabbit-reminder-1.0.0.jar" ]