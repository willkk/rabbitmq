package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	conn, err := Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Println("dail failed.")
		return
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		fmt.Println("open channel failed.")
		return
	}
	exchange := "test_topic"
	kind := "topic"
	err = DeclareExchange(ch, exchange, kind)
	if err != nil {
		fmt.Println("declare exchange failed.")
		return
	}

	q, err := DeclareQueue(ch, "")
	if err != nil {
		fmt.Println("open channel failed.")
		return
	}
	key := "*.info.*"
	err = QueueBind(ch, exchange, q.Name, key)
	if err != nil {
		fmt.Println("declare exchange failed.")
		return
	}

	del, err := Consume(ch, q.Name, false)
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



