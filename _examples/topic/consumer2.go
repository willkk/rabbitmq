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

	// If a key matches multiple key patterns, this request is sent only once.
	// For example, a request with key "will.error.develop", it matches "will.*.*",
	// "*.error.*" and "*.*.develop", but consumer only receive this request once.
	keys := []string{"will.*.*", "*.error.*", "*.*.develop"}
	for _, key := range keys {
		err = QueueBind(ch, exchange, q.Name, key)
		if err != nil {
			fmt.Println("declare exchange failed.")
			return
		}
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



