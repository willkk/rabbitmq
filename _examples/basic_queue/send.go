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
	queue := "test_queue"
	_, err = DeclareQueue(ch, queue)
	if err != nil {
		fmt.Println("open channel failed.")
		return
	}
	body := "hello, test queue"
	err = Publish(ch, "", queue, amqp.Publishing{
		Type: "plain/text",
		Body: []byte(body),
	})

	if err != nil {
		fmt.Printf("send failed. err=%s\n", body, err)
		return
	}
	fmt.Printf("send:%s\n", body)
}



