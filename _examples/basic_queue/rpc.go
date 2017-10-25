package main

import (
	. "rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Printf("dail failed.")
		return
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		fmt.Printf("open channel failed.")
		return
	}

	q, err := DeclareQueue(ch, "")
	if err != nil {
		fmt.Println("open channel failed.")
		return
	}
	queue := "test_queue"
	pub := amqp.Publishing{
		Type: "plain/text",
		CorrelationId: "0_0_1",
		ReplyTo: q.Name,
		Body: []byte("hello"),
	}
	del, err := Rpc(ch, queue, pub)
	if err != nil {
		fmt.Printf("consume failed. err=%s", err)
		return
	}

	fmt.Printf("rpc:%s\n", del.Body)
}

