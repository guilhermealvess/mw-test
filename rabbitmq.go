package mw

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQEvent struct {
	ctx     context.Context
	id      string
	content []byte
	msg     *amqp.Delivery
}

func (r *rabbitMQEvent) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *rabbitMQEvent) Context() context.Context { return r.ctx }

func (r *rabbitMQEvent) Bind(v interface{}) error {
	return json.Unmarshal(r.content, v)
}

func (r *rabbitMQEvent) ID() string { return r.id }

type rabbitMQRouter struct {
	timeout    time.Duration
	autoCommit bool
	queue      string
	name       string
	workerPool int
	handler    HandlerFunction
	connection *amqp.Connection
}

func (r *rabbitMQRouter) Commit(e Event) error {
	if err := e.(*rabbitMQEvent).msg.Ack(false); err != nil {
		return err
	}
	return nil
}

func (r *rabbitMQRouter) Producer(buffer chan<- Event) {
	ch, err := r.connection.Channel()
	if err != nil {
		return
	}

	msgs, err := ch.Consume(
		r.queue,
		r.name,
		r.autoCommit,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return
	}

	defer ch.Close()

	for msg := range msgs {
		buffer <- &rabbitMQEvent{
			ctx:     context.Background(),
			id:      uuid.NewString(),
			content: msg.Body,
			msg:     &msg,
		}
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
	conn                  *amqp.Connection
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
		connection: r.conn,
	}
}

func NewRabbitMQBuilder(name, queue string, handler HandlerFunction, conn *amqp.Connection) *RabbitMQBuilder {
	return &RabbitMQBuilder{
		workerName: name,
		queueName:  queue,
		timeout:    time.Second * 30,
		autocommit: false,
		handler:    handler,
		workerPool: 10,
		conn:       conn,
	}
}

func NewConnectRabbitMQ(uri string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	return conn, nil
}
