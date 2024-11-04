package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// потребитель прослушивает топик заказов с названием "coffee_orders"
	topic := "coffee_orders"
	// количество принятых сообщений
	msgCount := 0

	// 1. Создаем нового потребителя и запустить его
	worker, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	// TODO: рассмотреть возможность такого подхода к закрытию воркера:
	//defer worker.Close()

	//читаем самую старую (OffsetOldest), нулевую партицию из топика
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer started")

	// 2. Обрабатываем сигналы ОС - для остановки процесса
	// создаем канал для сигналов ОС
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// 3. Создаем горутину для запуска потребителя (воркера)
	//канал для ощения с горутиной
	doneCh := make(chan struct{})
	//создаем горутину для обработки всех сообщений (работает бесконечно)
	go func() {
		for {
			select {

			case err := <-consumer.Errors(): // вывод ошибок, возникающих в потребителе
				fmt.Println(err)

			case msg := <-consumer.Messages(): // обработка сообщений для потребителя
				//пока что просто увеличим счетчик сообщений
				msgCount++
				//вывод принятого сообщения в лог
				fmt.Printf("Принято: Количество (%d) | Topic (%s) | Message (%s) \n", msgCount, string(msg.Topic), string(msg.Value))
				//извлекаем заказ из сообщения (сообщение из кафки представляет собой слайс байтов, образованный из строки)
				order := string(msg.Value)
				//симулируем подготовку на основе заказа
				fmt.Printf("Готовим кофе для заказа: %s\n", order)

			case <-sigchan: //обработка сигналов ОС
				fmt.Println("обнаружено прерывание")
				doneCh <- struct{}{}
			}
		}
	}()

	// ожидаем получения сообщения из канала doneCh
	<-doneCh

	//как только канал закончил работу, выводим общее количество обработанных сообщений
	fmt.Println("Обработано сообщений: ", msgCount)

	// 4. Закрываем потребителя по выходу (может лучше вынести в defer?)
	if err := worker.Close(); err != nil {
		panic(err)
	}
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	// создаем конфигурацию
	config := sarama.NewConfig()

	//задаем параметры
	config.Consumer.Return.Errors = true // возвращение ошибок

	// создаем нового продюсера и возвращаем его
	return sarama.NewConsumer(brokers, config)
}
