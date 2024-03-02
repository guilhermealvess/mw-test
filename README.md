# mw
Multi Worker Golang

# Example
```go
    manager := mw.New()
    worker := mw.NewRabbitMQBuilder("worker", "worker_example__queue", func(e Event) error { return nil }).
        Timeout(time.Secound*10).
        WorkerPool(10*1.5).
        Build()

    manager.Register(worker)
    manager.Run()
```