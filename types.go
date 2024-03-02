package mw

import "context"

type Event interface {
	SetContext(context.Context)
	Context() context.Context
	Bind(interface{}) error
	ID() string
}

type HandlerFunction func(Event) error

type Router interface {
	Commit(Event) error
	Producer(chan<- Event)
	Handler() HandlerFunction
	Pool() int
}
