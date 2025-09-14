## Структура 
<img width="626" height="371" alt="Untitled Diagram drawio(1)" src="https://github.com/user-attachments/assets/f199d884-f606-4f4a-a5f9-9ad3a06fb7e8" />

## Технологии

- Kafka
- PostgreSQL
- Redis
- Docker, Docker Compose

## Демонстрация
![ezgif-26a4a38df78c75](https://github.com/user-attachments/assets/887073a1-58a9-40a4-8462-9385e5784ed3)


## Запуск

1. **Клонируйте репозиторий**
```bash
git clone https://github.com/Hochmuch/wb-tech-L0
```
2. **Соберите образы проекта**
```bash
make build
```
3. **Запустите сервис**
```bash
make up
```
4. Сервис доступен в браузере по адресу [http://localhost:8080](http://localhost:8080)

## Тестирование

1. **Сгенерируйте моки**

```bash
make mocks
```

2. **Запустите тесты**

  Они находятся в internal/repository и internal/service

3. **Отправка своих сообщений**

  Вы можете отправить свои сообщения с помощью kafka-ui, он расположен по адресу [http://localhost:8090](http://localhost:8090)

