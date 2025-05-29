## Структура проекта

```
├── bin/ # Скомпилированные бинарники
├── cmd/ # Точка входа: основное приложение и мигратор
│   ├── eduhelper/
│   │   └── main.go
│   └── migrator/
│       └── main.go
├── config/ # Конфигурационные файлы
├── internal/ # Основная бизнес-логика, хранилища, HTTP сервер
├── migrations/ # SQL-миграции для БД
├── Makefile # Автоматизация сборки и миграций
└── go.mod, go.sum # Go-модули и зависимости
```

## Быстрый старт

### 1. Клонирование репозитория

```sh

git clone https://github.com/devhuiw/eduhelper

cd eduhelper

```

### 2. Настройка переменных окружения/конфига

Измени параметры в `config/local.yaml` или передавай переменные в команды Make напрямую.

### 3. Миграции базы данных

**Выполнить миграции вверх:**

```sh

make migrate-up user=<db_user> password=<db_password> db_name=<db_name> host=<db_host> port=<db_port>

```

**Откатить миграции вниз:**

```sh

make migrate-down user=<db_user> password=<db_password> db_name=<db_name> host=<db_host> port=<db_port>

```

- Все параметры можно не указывать, тогда используются значения по умолчанию (`root`, `localhost`, `3306`, `test`).
- Путь к миграциям: `./migrations`

### 4. Сборка и запуск приложения

**Собрать проект:**

```sh

make build

```

**Запустить сервис:**

```sh

make run

```

## Тестирование

```sh

make test

```

## Линтинг

```sh

make lint

```

> Не забудь установить [golangci-lint](https://golangci-lint.run/usage/install/), если он ещё не установлен.

## Структура Makefile

- `build` — сборка бинарника
- `run` — запуск приложения
- `test` — тесты
- `lint` — линтинг
- `tidy` — обновить зависимости
- `clean` — очистить сборку
- `migrate-up` — миграции вверх
- `migrate-down` — миграции вниз
