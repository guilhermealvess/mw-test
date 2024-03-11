# mw
Multi Worker Golang

# Example
```go
    manager := mw.New()
    worker := mw.NewRabbitMQRouter(
		"handlerName",
		"queueName",
		conn,
		handler,
		mw.WithTimeout(timeout),
		mw.WithAutoCommit(),
		mw.WithMultiplierWorkerPool(2),
	)

    manager.Register(worker)
    manager.Run()
```