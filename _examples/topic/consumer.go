package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	t, err := NewTopic("amqp://guest:guest@localhost:5672", "test_topic")
	if err != nil {
		fmt.Printf("NewTopic failed. err=%s\n", err)
		return
	}
	defer t.Close()

	keys := []string{"*.info.*"}
	del, err := t.Consume(keys, false)
	for d := range del {
		fmt.Printf("receive:%s\n", d.Body)
		if d.ReplyTo != "" {
			resp := fmt.Sprintf("resp to %s", d.Body)
			err := t.Reply(d.ReplyTo, amqp.Publishing{
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



