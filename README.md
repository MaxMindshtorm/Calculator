# HTTP Сервер подсчёта арифметических выражений

## Что это вообще за проект

Этот мини-проект был написан по заданию из курса Яндекс лицея. Это калькулятьр, который просто вычисляет значение по предоставленному выражению

## Как этим пользоваться?

### Запуск

Для запуска нужно ввести команду:

```
go run cmd/main.go
```

### Использование

Этот сервер принимает http POST запросы на адрес `/api/v1/calculate`.
Тип запросов:

```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "2+2"
}'
```

Сервер поддерживает, рацианальные числа арифмитические операции (-, +, \*, /), а так же скобки и унарный минус

### Примеры использования

Тут будут представлены только тела запросов, подставив их в ранне указанную команду, вы сможете получить те же результаты
| Name | Response status | Method | Path | Body |
| --- | --- | --- | --- | --- |
| OK | 200 | POST | /api/v1/calculate | `{"expression:"2+2"}` |
| Method is not allowed | 405 | Not POST | /api/v1/calculate | `{"expression:"2+2"}` |
| Expression is not valid | 422 | POST | /api/v1/calculate | `{"expression:"2*(2+2"}"` |
| Wrong Path | 404 | POST | /any | `{"expression:"2+2"}` |
| Bad request | 400 | POST | /api/v1/calculate | `invalid body` |
