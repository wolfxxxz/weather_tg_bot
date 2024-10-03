# Weather Bot with MongoDB

# This project implements a weather forecast bot with the following key features:
    MongoDB for storing user data.
    Ticker to check weather information periodically.
    Go Routines for handling concurrent flows efficiently.
    Primitive Unit Testing for basic test coverage.
    Integration with third-party API to fetch weather forecasts.
    Docker Compose for deploying the project in the cloud.
    Mocks for testing the HTTP client used in API interactions.

# Packages Used:
    github.com/caarlos0/env - v3.5.0 (for environment variable management)
    github.com/go-telegram-bot-api/telegram-bot-api/v5 - v5.5.1 (Telegram Bot API)
    github.com/golang/mock - v1.6.0 (for mocking interfaces during testing)
    github.com/google/uuid - v1.3.1 (for generating unique user IDs)
    github.com/joho/godotenv - v1.5.1 (loading environment variables from .env files)
    github.com/sirupsen/logrus - v1.9.3 (for logging)
    github.com/stretchr/testify - v1.8.4 (for unit testing assertions)
    go.mongodb.org/mongo-driver - v1.12.1 (MongoDB Go driver)
    gopkg.in/mgo.v2 - v2.0.0 (MongoDB interaction)