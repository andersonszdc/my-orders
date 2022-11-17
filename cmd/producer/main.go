package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/andersonszdc/my-orders/internal/order/entity"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ch *amqp.Channel, order entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func GenerateOrders() entity.Order {
	return entity.Order{
		ID:    uuid.New().String(),
		Price: rand.Float64() * 100,
		Tax:   rand.Float64() * 10,
	}
}

func main() {
	var url string = os.Getenv("RABBITMQ_DEFAULT_URL")

	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	for i := 0; i < 100; i++ {
		Publish(ch, GenerateOrders())
		time.Sleep(300 * time.Millisecond)
	}
}
