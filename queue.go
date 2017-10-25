package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Queue struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
	Q  *amqp.Queue
	Addr string
}

func NewQueue(addr, name string) (*Queue, error) {
	conn, err := Dial(addr)
	if err != nil {
		return nil, err
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		return nil, err
	}

	q, err := DeclareQueue(ch, name)
	if err != nil {
		return nil, err
	}

	return &Queue{conn, ch, q, addr}, nil
}

func (q *Queue)Publish(publishing amqp.Publishing) error {
	return Publish(q.Ch, "", q.Q.Name, publishing)
}

func (q *Queue)Reply(key string, publishing amqp.Publishing) error {
	return Publish(q.Ch, "", key, publishing)
}

func (q *Queue)Consume(autoAck bool)(<-chan amqp.Delivery, error) {
	return Consume(q.Ch, q.Q.Name, autoAck)
}

func (q *Queue)Rpc(key string, publishing amqp.Publishing) (*amqp.Delivery, error) {
	return Rpc(q.Ch, key, publishing)
}