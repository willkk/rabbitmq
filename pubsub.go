package rabbitmq

import "github.com/streadway/amqp"

type PubSub struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
	Exchange string
	Addr string
}

func NewPubSub(addr, exchange string) (*PubSub, error) {
	conn, err := Dial(addr)
	if err != nil {
		return nil, err
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		return nil, err
	}

	err = DeclareExchange(ch, exchange, "fanout")
	if err != nil {
		return nil, err
	}

	return &PubSub{conn, ch, exchange, addr}, nil
}

func (pb *PubSub)Publish(publishing amqp.Publishing) error {
	return Publish(pb.Ch, pb.Exchange, "", publishing)
}

func (pb *PubSub)Subscribe(autoAck bool)(<-chan amqp.Delivery, error) {
	q, err := DeclareQueue(pb.Ch, "")
	if err != nil {
		return nil, err
	}

	err = QueueBind(pb.Ch, pb.Exchange, q.Name, "")
	if err != nil {
		return nil, err
	}

	return Consume(pb.Ch, q.Name, autoAck)
}

func (pb *PubSub)Reply(key string, publishing amqp.Publishing) error {
	return Publish(pb.Ch, "", key, publishing)
}

