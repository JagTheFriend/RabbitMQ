package main

import (
	"common"
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
			return
		}
		e = ch.Close()
		if e != nil {
			log.Fatal(e)
			return
		}
	}()

	listen(ch)
}

func listen(ch *amqp091.Channel) {
	q, err := ch.QueueDeclare(common.OrderCreatedEvent, true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			o := common.Order{}
			if err := json.Unmarshal(d.Body, &o); err != nil {
				if err := d.Nack(false, false); err != nil {
					log.Printf("Error nacking message: %s", err)
				}
				log.Printf("Error unmarshaling order: %s", err)
				continue
			}

			paymentLink, err := createPaymentLink(o.ID)
			if err != nil {
				log.Printf("Error creating payment link: %s", err)
				continue
			}

			log.Printf("Payment link created for order %s: %s", o.ID, paymentLink)
		}
	}()

	log.Println("AMQP Payments Service is listening for OrderCreated events")
	<-forever
}

func createPaymentLink(orderID string) (string, error) {
	return "http://payment-link.com/order/" + orderID, nil
}
