# Создание и запуск сервера кафки:
1. Устанавливаем образ в Docker
```bash
docker pull apache/kafka:3.7.0
```

2. Запускаем контейнер
```bash
docker run -p 9092:9092 apache/kafka:3.7.0
```


Web-client for Kafka (https://github.com/provectus/kafka-ui):
```bash
# Установка в Docker
docker pull provectuslabs/kafka-ui

# Запуск образа
docker run -it -p 8080:8080 -e DYNAMIC_CONFIG_ENABLED=true provectuslabs/kafka-ui
```



