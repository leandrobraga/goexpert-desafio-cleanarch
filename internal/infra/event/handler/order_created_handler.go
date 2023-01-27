package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/leandrobraga/goexpert-desafio-cleanarch/pkg/events"
	ampq "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *ampq.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *ampq.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msqRabbitmq := ampq.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}
	h.RabbitMQChannel.PublishWithContext(
		context.Background(),
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msqRabbitmq,  // messge to publish
	)

}
