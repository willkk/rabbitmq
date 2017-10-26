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

	keys := []string{"will.info.develop", "tom.warn.production", "lucy.error.develop",
	"will.test", "tom", "will.test.dev.pro"}
	for _, key := range keys {
		body := fmt.Sprintf("hello, test route, key:%s", key)
		err = t.Publish(key, amqp.Publishing{
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



