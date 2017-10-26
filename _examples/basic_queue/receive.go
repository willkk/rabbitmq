package main

import (
	. "rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

// If multiple receivers are launched, the requests are sent to them in round-robin mode.
// In this way, every receiver only gets part of all requests.
func main() {
	q, err := NewQueue("amqp://guest:guest@localhost:5672", "test_queue")
	if err != nil {
		fmt.Printf("NewQueue failed. err=%s\n", err)
		return
	}
	defer q.Close()

	del, err := q.Consume(false)
	for d := range del {
		fmt.Printf("receive:%s\n", d.Body)
		if d.ReplyTo != "" {
			resp := fmt.Sprintf("resp to %s", d.Body)
			err := q.Reply(d.ReplyTo, amqp.Publishing{
				Type: "plain/text",
				Body: []byte(resp)})

			if err != nil {
				fmt.Printf("resp failed. err=%s\n", err)
				return
			}

			// send ack when all works is done by consumer.
			d.Ack(false)
		}
	}
}
