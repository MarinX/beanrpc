// beanrpc
package beanrpc

import (
	"encoding/json"
	"github.com/kr/beanstalk"
	"time"
)

type BeanWorker struct {
	address string
	conn    *beanstalk.Conn
	rpc     map[string]HandlerFunc
	closed  bool
	running bool
}

type HandlerFunc func(c *Context)

type Request struct {
	Method string
	Params interface{}
}

type Context struct {
	buff []byte
	id   uint64
}

func New(address string) *BeanWorker {
	return &BeanWorker{
		address: address,
		rpc:     make(map[string]HandlerFunc),
	}
}

func (t *Context) Buff() []byte {
	return t.buff
}

func (t *Context) Bind(rcv interface{}) error {
	var r Request
	r.Params = rcv
	return json.Unmarshal(t.buff, &r)
}

func (t *Context) Id() uint64 {
	return t.id
}

func (t *BeanWorker) Open(tube string) error {
	conn, err := beanstalk.Dial("tcp", t.address)
	if err != nil {
		return err
	}
	conn.Tube = beanstalk.Tube{conn, tube}
	conn.TubeSet = *beanstalk.NewTubeSet(conn, tube)
	t.conn = conn
	return nil
}

func (t *BeanWorker) On(name string, h HandlerFunc) {
	t.rpc[name] = h
}

func (t *BeanWorker) Put(method string, params interface{}, priority uint32) error {
	var r Request
	r.Method = method
	r.Params = params
	buff, err := json.Marshal(&r)
	if err != nil {
		return err
	}

	t.conn.Put(buff, priority, 0*time.Second, 180*time.Second)
	return nil
}

func (t *BeanWorker) Run() {
	t.running = true
	for {
		if t.closed {
			break
		}

		id, buff, err := t.conn.Reserve(5 * time.Second)
		if err != nil {
			continue
		}

		var r Request
		err = json.Unmarshal(buff, &r)
		if err == nil {
			if t.rpc[r.Method] != nil {
				t.rpc[r.Method](&Context{
					buff: buff,
					id:   id,
				})
			}
		}
		t.conn.Delete(id)
	}
	t.running = false
}

func (t *BeanWorker) Close() error {
	t.closed = true
	for {
		if !t.running {
			break
		}
		time.Sleep(1 * time.Second)

	}
	return t.conn.Close()
}

