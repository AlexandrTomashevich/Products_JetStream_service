# Products JetStream service

Проект "Products JetStream Service" - это web-сервис для просмотра и управления товарами. Он позволяет пользователю просматривать список продуктов, а также подробные характеристики каждого товара. Интеграция с JetStream обеспечивает надежное и эффективное взаимодействие с данными.

## Основные характеристики:

1. Подключение и подписка на канал JetStream для обработки данных.
2. Полученные данные пишутся в Postgres.
3. Так же полученные данные сохраняются in memory в сервисе (Кеш).
4. В случае падения сервиса восстанавливаем Кеш из Postgres.
5. Поднимаем http сервер и выдавать данные по id из кеша.
6. Интуитивно понятный web-интерфейс отображения полученных данных, для их запроса по id.
7. Возможность просмотра детальной информации о товаре нажав на его "id".

## Пример работы:

После запуска всех компонентов и перехода на главную страницу **localhost:8000/**, вы увидите два продукта. Каждый продукт будет представлен своим именем, ценой и производителем. После запуска авто-теса добавится еще один продукт (тестовый). Это результат работы сервиса, который извлекает данные из базы данных и отображает их в удобном для пользователя виде.

## Инструкция для запуска
### Предварительные условия:
- Установить docker (docker-desktop) и docker-compose
  - [инструкцию по установке](https://docs.docker.com/get-docker/)

### Шаги для запуска:
В процессе запуска периодически вы будете видеть сообщение "invalidate cache", которое является частью желаемого поведения нашего приложения, оно означает что кэш был очищен (по логике каждые 15 секунд).
- Запустить jetstream
```bash
docker-compose -f deployments/docker-compose.yaml up -d
```

- Запустить сервер
```bash
go run cmd/app/main.go
```

- Запустить consumer
```bash
go run cmd/consumer/main.go
```

- Запустить producer
```bash
go run cmd/producer/main.go
```

- Проверка:
  Перейти на страницу **localhost:8000/**. Вы должны увидеть страницу со списком товаров.
  На странице **localhost:8000/products**. Вы должны увидеть страницу со списком товаров в JSON-формате.
  На странице **localhost:8000/product**. Вы перейдете на страницу, где увидите сообщение, о необходимости добавить к адресу номер (id) товара например: **localhost:8000/product?id=1** и автоматически будете переадресованы на начальную страницу с товарами.

### Тестирование:
Убедитесь, что все компоненты запущены и работают корректно. Проверьте отображение продуктов на главной странице и удостоверьтесь, что все данные корректно отображаются.
- Для проведения авто-тестирования сервиса выполните следующую команду:
```bash
go test ./test
```

## Зависимости и используемые библиотеки

В данном проекте были использованы следующие зависимости и библиотеки:
- Go 
- NATS
  - github.com/nats-io/nats.go
- PostgreSQL
  - github.com/lib/pq
- HTML/Go templates
  - text/template 
- Другие утилиты и библиотеки Go:
  - database/sq
  - log 
  - net/http 
  - encoding/json 
  - os и syscall
  - sync 
  - strconv 
  - time 
  - testing 