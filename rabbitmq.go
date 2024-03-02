package mw

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type rabbitMQEvent struct {
	ctx     context.Context
	id      string
	content []byte
}

func (r *rabbitMQEvent) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *rabbitMQEvent) Context() context.Context { return r.ctx }

func (r *rabbitMQEvent) Bind(v interface{}) error {
	return json.Unmarshal(nil, v)
}

func (r *rabbitMQEvent) ID() string { return r.id }

type rabbitMQRouter struct {
	timeout    time.Duration
	autoCommit bool
	queue      string
	name       string
	workerPool int
	handler    HandlerFunction
}

func (r *rabbitMQRouter) Commit(e Event) error {
	// TODO:
	return nil
}

func (r *rabbitMQRouter) Producer(buffer chan<- Event) {
	// TODO:
	for {
		e := rabbitMQEvent{
			ctx:     context.Background(),
			id:      uuid.NewString(),
			content: json.RawMessage(`{}`),
		}
		buffer <- &e
		time.Sleep(time.Second)
	}
}

func (r *rabbitMQRouter) Handler() HandlerFunction { return r.handler }

func (r *rabbitMQRouter) Pool() int { return r.workerPool }

type RabbitMQBuilder struct {
	workerName, queueName string
	timeout               time.Duration
	autocommit            bool
	handler               HandlerFunction
	workerPool            int
}

func (r *RabbitMQBuilder) Timeout(t time.Duration) *RabbitMQBuilder {
	r.timeout = t
	return r
}

func (r *RabbitMQBuilder) AutoCommit(b bool) *RabbitMQBuilder {
	r.autocommit = b
	return r
}

func (r *RabbitMQBuilder) WorkerPool(n int) *RabbitMQBuilder {
	r.workerPool = n
	return r
}

func (r *RabbitMQBuilder) Build() *rabbitMQRouter {
	return &rabbitMQRouter{
		queue:      r.queueName,
		timeout:    r.timeout,
		autoCommit: r.autocommit,
		name:       r.workerName,
		handler:    r.handler,
		workerPool: r.workerPool,
	}
}

func NewRabbitMQBuilder(name, queue string, handler HandlerFunction) *RabbitMQBuilder {
	return &RabbitMQBuilder{
		workerName: name,
		queueName:  queue,
		timeout:    time.Second * 30,
		autocommit: false,
		handler:    handler,
		workerPool: 10,
	}
}
