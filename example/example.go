package example

import (
	"time"

	"github.com/guilhermealvess/mw"
)

var (
	handlerName = "example"
	queueName   = "example__queue"
	timeout     = time.Second * 30
)

func handler(e mw.Event) error {
	var m map[string]interface{}
	if err := e.Bind(&m); err != nil {
		return err
	}

	return nil
}

func main() {
	conn, err := mw.NewConnectRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	workerOne := mw.NewRabbitMQRouter(
		handlerName,
		queueName,
		conn,
		handler,
		mw.WithTimeout(timeout),
		mw.WithAutoCommit(),
		mw.WithMultiplierWorkerPool(2),
	)

	workerTwo := mw.NewRabbitMQRouter(
		handlerName,
		queueName,
		conn,
		handler,
	)

	manager := mw.New()
	manager.Register(workerOne, workerTwo)
	manager.Run()
}
