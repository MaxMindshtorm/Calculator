# Приветствую тебя, проверяющий этого проекта, проект ещё не дописан, не мог бы та прийти сюда через пару дней и ты встетишь полностью функцианирующий проект, мне ведь действительно немного осталось, но произошел некий завал на учебе и ещё несколько неприятных ситуация, так что надеюсь на понимание и верусь сида с нормальным проектом и ридмишкой через пару дней!
# HTTP Сервер подсчёта арифметических выражений

## Что это вообще за проект

Этот мини-проект был написан по заданию из курса Яндекс лицея. Это калькулятор, который просто вычисляет значение по предоставленному выражению

## Как этим пользоваться?

### Установка

Клонируйте проект

```
git clone https://github.com/MaxMindshtorm/Calculator.git
cd Calculator
```

### Запуск

Для запуска нужно ввести команду:

```
go run cmd/main.go
```

### Использование

Этот сервер принимает http POST запросы на адрес `/api/v1/calculate`.
По умолчанию используется порт 8080, но его можно поменять на удобный вам при запуске.
К примеру вот так:

```
$env:PORT=8888; go run .\cmd\main.go
```

Тип запросов, которые можно отправлять:

```
curl --location "http://localhost:8080/api/v1/calculate" ^
--header "Content-Type: application/json" ^
--data "{\"expression\": \"2+2\"}"
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

## Тестирование

Если вам захочется проверить корректность работы данного калькулятор, то вы можете запустить тесты, которые я прописал для этого, так же в файлах "...\_test.go" вы можете добавлять свои тесты

### Запуск тестов

Для того чтобы запустить тесты, с помощью команды `cd Calculator` перейдите в корневую папку проекта и используйте команду

```
go test ./...
```

Это запустит тестирование пакета с хэндлерами и пакета с калькулятором, так что вы сможете удостовериться в работоспособности проекта без запуска кода
