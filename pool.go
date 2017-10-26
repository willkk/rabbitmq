package rabbitmq

import (
	"github.com/streadway/amqp"
	"errors"
	"sync"
)

type Conn struct {
	Conn *amqp.Connection
	InUse bool
}

type Pool struct {
	MaxCount int
	MinCount int
	CurrCount int
	CurrInUse int
	Addr string
	Conns map[*amqp.Connection]Conn

	lock *sync.Mutex
}

func NewPool(addr string, max, min int) (*Pool, error) {
	conns := make(map[*amqp.Connection]Conn)
	for i := 0; i < min; i++ {
		conn, err := Dial(addr)
		if err != nil {
			return nil, err
		}
		conns[conn] = Conn{conn, false}
	}

	p := &Pool{
		max,
		min,
		min,
		0,
		addr,
		conns,
		new(sync.Mutex),
	}

	return p, nil
}

func (p *Pool)Conn() (*amqp.Connection, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.CurrCount == p.MaxCount && p.CurrCount == p.CurrInUse {
		return nil, errors.New("exceeds max connection")
	}

	// open new one
	if p.CurrInUse == p.CurrCount {
		conn, err := Dial(p.Addr)
		if err != nil {
			return nil, err
		}
		p.Conns[conn] = Conn{conn, true}
		p.CurrCount++
		p.CurrInUse++
		return conn, nil
	}

	// get one from pool
	for _, cn := range p.Conns {
		if cn.InUse == false {
			cn.InUse = true
			p.CurrInUse++
			return cn.Conn, nil
		}
	}

	return nil, errors.New("no conn available")
}

func (p *Pool)Release(conn *amqp.Connection) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for key, cn := range p.Conns {
		if key == conn {
			if p.CurrCount == p.MinCount {
				cn.InUse = false
				p.CurrInUse--
			} else {
				key.Close()
				delete(p.Conns, key)
				p.CurrCount--
				p.CurrInUse--
			}

			break
		}
	}
}

func (p *Pool)Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for key, _ := range p.Conns {
		key.Close()
		delete(p.Conns, key)
	}
}