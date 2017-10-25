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

	keys := []string{"will.info.develop", "tom.warn.production", "lucy.error.develop",
	"will.test", "tom", "will.test.dev.pro"}
	for _, key := range keys {
		body := fmt.Sprintf("hello, test route, key:%s", key)
		err = Publish(ch, exchange, key, amqp.Publishing{
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



