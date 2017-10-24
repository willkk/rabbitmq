package rabbitmq

import (
	"testing"
	"github.com/streadway/amqp"
)

func TestSend(t *testing.T) {
	conn, err := Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		t.Fatalf("dail failed.")
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		t.Fatalf("open channel failed.")
	}
	body := "hello, test queue"
	err = Publish(ch, "test_queue", amqp.Publishing{
		Type: "plain/text",
		Body: []byte(body),
	})

	if err != nil {
		t.Errorf("send failed. err=%s", body, err)
	}
	t.Logf("send:%s", body)
}

func TestRpc(t *testing.T) {
	MQInit("amqp://guest:guest@localhost:5672")
	conn, err := Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		t.Fatalf("dail failed.")
	}
	ch, err := OpenChannel(conn)
	if err != nil {
		t.Fatalf("open channel failed.")
	}

	body := "hello, test queue"
	del, err := Rpc(ch,"test_queue", amqp.Publishing{
		Type: "plain/text",
		Body: []byte(body),
	})

	if err != nil {
		t.Errorf("rpc failed. err=%s", err)
	}

	t.Logf("rpc:%s, recv:%s", body, del.Body)
}

func TestReceive(t *testing.T) {
	MQInit("amqp://guest:guest@localhost:5672")
	ch, err := Receive("test_queue")

	if err != nil {
		t.Errorf("rpc failed. err=%s", err)
	}

	for d := range ch {
		t.Logf("receive:%s", d.Body)
	}
}

