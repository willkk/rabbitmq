package rabbitmq

import (
	"github.com/streadway/amqp"
	"errors"
)

type Topic struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
	Exchange string
	Addr string
	pool *Pool
}

func NewTopic(addr, exchange string) (*Topic, error) {
	pool, err := NewPool(addr, 1, 1)
	if err != nil {
		return nil, err
	}
	conn, err := pool.Conn()
	if err != nil {
		return nil, err
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		return nil, err
	}

	err = DeclareExchange(ch, exchange, "topic")
	if err != nil {
		return nil, err
	}

	return &Topic{conn, ch, exchange, addr, pool}, nil
}

func (t *Topic)Publish(key string, publishing amqp.Publishing) error {
	return Publish(t.Ch, t.Exchange, key, publishing)
}

func (t *Topic)Close() {
	t.Ch.Close()
	t.pool.Close()
}

func (t *Topic)Consume(keys []string, autoAck bool)(<-chan amqp.Delivery, error) {
	if len(keys) == 0 {
		return nil, errors.New("empty keys")
	}

	q, err := DeclareQueue(t.Ch, "")
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		err = QueueBind(t.Ch, t.Exchange, q.Name, key)
		if err != nil {
			return nil, err
		}
	}

	return Consume(t.Ch, q.Name, autoAck)
}

func (t *Topic)Reply(key string, publishing amqp.Publishing) error {
	return Publish(t.Ch, "", key, publishing)
}