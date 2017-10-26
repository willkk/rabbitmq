# Rabbitmq API Wrapper
A simple rabbitmq API wrapper library. It handles four common MQ types: Queue, Direct, PubSub and Topic.

## 4 Types
### 1. Queue  
**send:**  
```go
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
	defer q.Close()

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
```
**receive:**  
```go
package main

import (
	. "rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

// If multiple receivers are launched, the requests are sent to them in round-robin mode.
// In this way, every receiver only gets part of all requests.
func main() {
	q, err := NewQueue("amqp://guest:guest@localhost:5672", "test_queue")
	if err != nil {
		fmt.Printf("NewQueue failed. err=%s\n", err)
		return
	}
	defer q.Close()

	del, err := q.Consume(false)
	for d := range del {
		fmt.Printf("receive:%s\n", d.Body)
		if d.ReplyTo != "" {
			resp := fmt.Sprintf("resp to %s", d.Body)
			err := q.Reply(d.ReplyTo, amqp.Publishing{
				Type: "plain/text",
				Body: []byte(resp)})

			if err != nil {
				fmt.Printf("resp failed. err=%s\n", err)
				return
			}

			// send ack when all works is done by consumer.
			d.Ack(false)
		}
	}
}
```
**rpc:**  
```go
package main

import (
	. "rabbitmq"
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	q, err := NewQueue("amqp://guest:guest@localhost:5672", "")
	if err != nil {
		fmt.Printf("NewQueue failed. err=%s\n", err)
		return
	}
	defer q.Close()

	key := "test_queue"
	pub := amqp.Publishing{
		Type: "plain/text",
		CorrelationId: "0_0_1",
		ReplyTo: q.Q.Name,
		Body: []byte("hello"),
	}

	del, err := q.Rpc(key, pub)
	if err != nil {
		fmt.Printf("consume failed. err=%s", err)
		return
	}

	fmt.Printf("rpc:%s\n", del.Body)
}


```
### 2. Director  
**producer:**  
```go
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
```
**consumer:**
```go
package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	di, err := NewDirector("amqp://guest:guest@localhost:5672", "test_route")
	if err != nil {
		fmt.Printf("NewDirector failed. err=%s\n", err)
		return
	}
	defer di.Close()

	keys := []string{"info"}
	del, err := di.Consume(keys, false)

	for d := range del {
		fmt.Printf("receive:%s\n", d.Body)
		if d.ReplyTo != "" {
			resp := fmt.Sprintf("resp to %s", d.Body)
			err := di.Reply(d.ReplyTo, amqp.Publishing{
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
``` 
### 3. PubSub  
**publish:**  
```go
package main

import (
	. "rabbitmq"
	"github.com/streadway/amqp"
	"fmt"
)

func main() {
	pb, err := NewPubSub("amqp://guest:guest@localhost:5672", "test_pubsub")
	if err != nil {
		fmt.Printf("NewPubSub failed. err=%s\n", err)
		return
	}
	defer pb.Close()

	body := "hello, test pubsub"
	err = pb.Publish(amqp.Publishing{
		Type: "plain/text",
		Body: []byte(body),
	})

	if err != nil {
		fmt.Printf("send failed. err=%s\n", body, err)
		return
	}
	fmt.Printf("send:%s\n", body)
}
```
**subscribe:**
```go
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
	defer pb.Close()

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
```
### 4. Topic  
**producer:**
```go
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
```
**consumer:**
```go
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
```
