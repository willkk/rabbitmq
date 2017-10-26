package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	d, err := NewDirector("amqp://guest:guest@localhost:5672", "test_route")
	if err != nil {
		fmt.Printf("NewDirector failed. err=%s\n", err)
		return
	}
	defer d.Close()

	keys := []string{"info", "warn", "error"}
	for _, key := range keys {
		body := fmt.Sprintf("hello, test route, key:%s", key)
		err = d.Publish(key, amqp.Publishing{
			Type: "plain/text",
			Body: []byte(body),
		})

		if err != nil {
			fmt.Printf("send failed. err=%s\n", body, err)
			return
		}
		fmt.Printf("send:%s\n", body)
	}

}



