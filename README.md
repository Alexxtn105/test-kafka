# Тестовый проект для Kafka

По материалам https://www.youtube.com/watch?v=4EdrCc29vXY

Клиент для Go:
```bash
go get github.com/IBM/sarama
```


# Запуск продюсера (producer):

```bash
go run main.go
```

# Запуск потребителя (consumer):

```bash
go run main.go
```



Тело тестового сообщение для Postman:
localhost:8400/orders
{
"customer_name":"John",
"coffee_type":"Latte"
}