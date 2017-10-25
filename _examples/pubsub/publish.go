package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	pb, err := NewPubSub("amqp://guest:guest@localhost:5672", "test_pubsub")
	if err != nil {
		fmt.Printf("NewPubSub failed. err=%s\n", err)
		return
	}

	body := "hello, test pubsub"
	err = pb.Publish(amqp.Publishing{
		Type: "plain/text",
		Body: []byte(body),
	})

	if err != nil {
		fmt.Printf("send failed. err=%s\n", body, err)
		return
	}
	fmt.Printf("send:%s\n", body)
}



