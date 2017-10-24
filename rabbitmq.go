package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"time"
	"errors"
)

// like amqp://guest:guest@localhost:5672
var mqAddr string


func MQInit(addr string) {
	mqAddr = addr
}

func Dial(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func OpenChannel(conn *amqp.Connection)(*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func DeclareQueue(ch *amqp.Channel, name string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// Receive
func Consume(ch *amqp.Channel, queue string) (<-chan amqp.Delivery, error) {
	d, err := ch.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// if rpc is false, Delivery in returned list should be ignored.
func Publish(ch *amqp.Channel, queue string, publishing amqp.Publishing) error {
	err := ch.Publish(
		"",
		queue,
		false,
		false,
		publishing)
	if err != nil {
		return err
	}

	return nil
}

func Rpc(ch *amqp.Channel, queue string, publishing amqp.Publishing) (*amqp.Delivery, error) {
	err := ch.Publish(
		"",
		queue,
		false,
		false,
		publishing)
	if err != nil {
		return nil, err
	}

	d, err := Consume(ch, publishing.ReplyTo)
	if err != nil {
		return nil, err
	}

	select {
	case <-time.After(time.Second*30):
		return nil, errors.New("time out")
	case resp := <-d:
		if resp.CorrelationId == publishing.CorrelationId {
			return &resp, nil
		}
	}

	return nil, errors.New("no CorrelationId matched")
}