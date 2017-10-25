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

	// If a key matches multiple key patterns, this request is sent only once.
	// For example, a request with key "will.error.develop", it matches "will.*.*",
	// "*.error.*" and "*.*.develop", but consumer only receive this request once.
	keys := []string{"will.*.*", "*.error.*", "*.*.develop"}
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



