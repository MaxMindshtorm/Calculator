# HTTP Сервер подсчёта арифметических выражений

## Что это вообще за проект

Этот мини-проект был написан по заданию из курса Яндекс лицея. Это калькулятор, который вычисляет значение по предоставленному выражению, по возможности делая это параллельно

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
go run .\orchestrator\cmd\main\main.go
go run .\agent\cmd\main\main.go
```

### Использование

В `About.md` пердаставленны примеры запросов и объяснение того, как работает данный сервис, прочитайте это перед использованием моего калькулятора.

В файле `Config.env` можно настраивать параметры работы сервисов

Сервер поддерживает, рацианальные числа арифмитические операции (-, +, \*, /), а так же скобки и унарный минус


## Тестирование

Если вам захочется проверить корректность работы данного калькулятор, то вы можете запустить тесты, которые я прописал для этого, так же в файлах "...\*_test.go" вы можете добавлять свои тесты

### Запуск тестов

Для того чтобы запустить тесты, с помощью команды `cd Calculator` перейдите в корневую папку проекта и используйте команду

```
go test ./...
```

Это запустит тестирование пакета с хэндлерами и пакета с калькулятором, так что вы сможете удостовериться в работоспособности проекта без запуска кода
По всем вопросам обращайтесь [telegram](https://t.me/YrikZh)
