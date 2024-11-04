# Создание и запуск сервера кафки:
1. Устанавливаем образ в Docker
```bash
docker pull apache/kafka:3.7.0
```

2. Запускаем контейнер
```bash
docker run -p 9092:9092 apache/kafka:3.7.0
```