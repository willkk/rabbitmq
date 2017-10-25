package main

import (
	. "rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

// If multiple receivers are launched, the requests are sent to them in round-robin mode.
// In this way, every receiver only gets part of all requests.
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
	queue := "test_queue"
	_, err = DeclareQueue(ch, queue)
	if err != nil {
		fmt.Println("open channel failed.")
	}
	del, err := Consume(ch, queue, false)
	if err != nil {
		fmt.Printf("consume failed. err=%s", err)
		return
	}

	for d := range del {
		fmt.Printf("receive:%s\n", d.Body)
		if d.ReplyTo != "" {
			resp := fmt.Sprintf("resp to %s", d.Body)
			err := Publish(ch, "", d.ReplyTo, amqp.Publishing{
				Type: "plain/text",
				CorrelationId: d.CorrelationId,
				Body: []byte(resp),
			})
			if err != nil {
				fmt.Printf("resp failed. err=%s\n", err)
				return
			}

			// send ack when all works is done by consumer.
			d.Ack(false)
		}
	}
}
