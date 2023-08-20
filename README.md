# Url-Shortener
A trained project inspired by [@ntuzov](https://habr.com/en/companies/selectel/articles/747738/)

## Local run
   * `go build cmd/url-shortener/main.go`
   * `CONFIG_PATH=./config/local.yaml ./main`

## Test through the [Hurl](https://hurl.dev/)
   * `hurl --test **/*.hurl`

## The roadmap

Начало - `13.08.2023`

Step | Status | Done
--- | :---: | ---
Выбор библиотек | ✅ | `Day 2`
Конфигурация приложения | ✅ | `Day 3`
Настраиваем logger | ✅ | `Day 3`
Пишем Storage | ✅ | `Day 6`
Handlers — обработчики запросов | ... |
Авторизация
Функциональные тесты
Деплой проекта