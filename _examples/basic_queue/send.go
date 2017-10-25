package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	q, err := NewQueue("amqp://guest:guest@localhost:5672", "test_queue")
	if err != nil {
		fmt.Printf("NewQueue failed. err=%s\n", err)
		return
	}
	defer q.Ch.Close()
	defer q.Conn.Close()

	body := "hello, test queue"
	err = q.Publish(amqp.Publishing{
		Type: "plain/text",
		Body: []byte(body)})

	if err != nil {
		fmt.Printf("send failed. err=%s\n", body, err)
		return
	}
	fmt.Printf("send:%s\n", body)
}



