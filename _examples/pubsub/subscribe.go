package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

// If multiple receivers are launched, the requests are sent to them, respectively.
// In this way, every receiver gets all of client requests.
func main() {
	pb, err := NewPubSub("amqp://guest:guest@localhost:5672", "test_pubsub")
	if err != nil {
		fmt.Printf("NewPubSub failed. err=%s\n", err)
		return
	}
	defer pb.Conn.Close()
	defer pb.Ch.Close()

	del, err := pb.Subscribe(false)
	for d := range del {
		fmt.Printf("receive:%s\n", d.Body)
		if d.ReplyTo != "" {
			resp := fmt.Sprintf("resp to %s", d.Body)
			err := pb.Reply(d.ReplyTo, amqp.Publishing{
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



