package main

import (
	. "rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	q, err := NewQueue("amqp://guest:guest@localhost:5672", "")
	if err != nil {
		fmt.Printf("NewQueue failed. err=%s\n", err)
		return
	}
	defer q.Close()

	key := "test_queue"
	pub := amqp.Publishing{
		Type: "plain/text",
		CorrelationId: "0_0_1",
		ReplyTo: q.Q.Name,
		Body: []byte("hello"),
	}

	del, err := q.Rpc(key, pub)
	if err != nil {
		fmt.Printf("consume failed. err=%s", err)
		return
	}

	fmt.Printf("rpc:%s\n", del.Body)
}

