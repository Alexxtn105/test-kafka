package main

// Продюсер. Его задача - получать заказы от пользователей и отправлять их в топик кафки

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"net/http"
)

// Order структура заказа
type Order struct {
	CustomerName string `json:"customer_name"`
	CoffeeType   string `json:"coffee_type"`
}

func main() {
	//	fmt.Println("Hello, World!")
	http.HandleFunc("/order", placeOrder)
	log.Fatal(http.ListenAndServe(":8400", nil))
}

// placeOrder Разместить POST-запрос в кафку
func placeOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 1. Парсим тело запроса в заказ
	order := new(Order)
	err := json.NewDecoder(r.Body).Decode(order)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Printf("%#v\n", order)
	// 2. Конвертируем тело в слайс байтов
	orderInBytes, err := json.Marshal(order)
	//fmt.Println(orderInBytes)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Отправляем байты в кафку, в параметрах указываем тему "coffee_orders" и сообщение в байтах
	err = PushOrderToQueue("coffee_orders", orderInBytes)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Направляем пользователю респонс - мапу с ключами success и msg
	response := map[string]any{
		"success": true,
		"msg":     "Order for " + order.CustomerName + " placed successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Error placing order", http.StatusInternalServerError)
		return
	}

}

// PushOrderToQueue функция отправки сообщения в топик кафки
func PushOrderToQueue(topic string, message []byte) error {
	// Настраиваем брокеры. В нашем случае он один localhost
	brokers := []string{"localhost:9092"}

	// Создаем новое подключение к брокеру
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	// В конце обязательно закрываем продюсера
	defer producer.Close()

	// Создаем сообщение для кафки
	msg := &sarama.ProducerMessage{
		Topic: topic,                         // топик
		Value: sarama.StringEncoder(message), // отправляемое сообщение (преобразуем слайс байтов в тип Encoder)
	}

	// Отправляем сообщение. Функция возвращает партицию, положение в очереди и ошибку
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Заказ сохранен в topic(%s)/partition(%d)/offset(%d)\n",
		topic,
		partition,
		offset,
	)
	return nil
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	// создаем конфигурацию
	config := sarama.NewConfig()

	//задаем параметры
	config.Producer.Return.Successes = true          // в случае успешной доставки сообщения возвращаются в канал Successes
	config.Producer.RequiredAcks = sarama.WaitForAll //Ожидаем все подтверждения
	config.Producer.Retry.Max = 5                    // максимальное количество повторов

	// создаем нового продюсера и возвращаем его
	return sarama.NewSyncProducer(brokers, config)
}
