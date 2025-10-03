-- Создаём схемы
CREATE SCHEMA IF NOT EXISTS custom;

-- Даём права пользователю на схемы
GRANT USAGE ON SCHEMA custom TO paldev;

-- Даём права на создание таблиц в схемах
GRANT CREATE ON SCHEMA custom TO paldev;
