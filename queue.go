package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
	"time"
	"errors"
)

type MQQueue struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
	Q *amqp.Queue
	Del <-chan amqp.Delivery
}

func (mq MQQueue)Init(addr, queue string) error {
	var err error
	mq.Conn, err = Dial(addr)
	if err != nil {
		return err
	}

	mq.Q, err = DeclareQueue(mq.Conn, queue)
	if err != nil {
		return err
	}



	return nil
}

func (mq MQQueue)Close() {
	if mq.Conn != nil {
		mq.Conn.Close()
	}
}

func (mq MQQueue)RPC(publishing amqp.Publishing)(
	*amqp.Delivery, error) {

	err := mq.Consume()
	if err != nil {
		log.Fatalf("Consume failed.")
		return nil, err
	}

	err = mq.Send(publishing)
	if err != nil {
		log.Fatalf("Send failed.")
		return nil, err
	}

	select {
	case <-time.After(time.Second*30):
		return nil, errors.New("time out")
	case resp := <-mq.Del:
		if resp.CorrelationId == publishing.CorrelationId {
			return &resp, nil
		}
	}

	return nil, errors.New("no CorrelationId matched")
}

func Receive(queue string) (<-chan amqp.Delivery, error){
	var mq MQQueue
	err := mq.Dial(mqAddr)
	if err != nil {
		return nil, err
	}

	err = mq.Channel()
	if err != nil {
		return nil, err
	}

	err = mq.Queue(queue)
	if err != nil {
		return nil, err
	}

	err = mq.Consume()
	if err != nil {
		return nil, err
	}

	return mq.Del, nil
}