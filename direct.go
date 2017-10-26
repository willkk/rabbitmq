package rabbitmq

import (
	"github.com/streadway/amqp"
	"errors"
)

type Director struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
	Exchange string
	Addr string
	pool *Pool
}

func NewDirector(addr, exchange string) (*Director, error) {
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

	err = DeclareExchange(ch, exchange, "direct")
	if err != nil {
		return nil, err
	}

	return &Director{conn, ch, exchange, addr, pool}, nil
}

func (d *Director)Publish(key string, publishing amqp.Publishing) error {
	return Publish(d.Ch, d.Exchange, key, publishing)
}

func (d *Director)Close() {
	d.pool.Close()
	d.Ch.Close()
}

func (d *Director)Consume(keys []string, autoAck bool)(<-chan amqp.Delivery, error) {
	if len(keys) == 0 {
		return nil, errors.New("empty keys")
	}

	q, err := DeclareQueue(d.Ch, "")
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		err = QueueBind(d.Ch, d.Exchange, q.Name, key)
		if err != nil {
			return nil, err
		}
	}

	return Consume(d.Ch, q.Name, autoAck)
}

func (d *Director)Reply(key string, publishing amqp.Publishing) error {
	return Publish(d.Ch, "", key, publishing)
}