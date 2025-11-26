package main

import (
	"common"
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

const (
	amqpUser     = "guest"
	amqpPassword = "guest"
	amqpHost     = "localhost"
	amqpPort     = "5672"
)

func main() {
	ch, close := common.ConnectRabbitAMQP(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		e := close()
		if e != nil {
			log.Fatal(e)
		}
	}()

	q, err := ch.QueueDeclare(
		common.OrderCreatedEvent,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	marshallerOrder, err := json.Marshal(common.Order{
		ID: "order-1",
		Items: []common.Item{
			{
				ID:       "item-1",
				Quantity: 1,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        marshallerOrder,
		})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Order created event published")

}
