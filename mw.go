package mw

import "sync"

type manager struct {
	routers []Router
}

func New() *manager {
	return &manager{
		routers: []Router{},
	}
}

func (m *manager) Register(r Router) {
	m.routers = append(m.routers, r)
}

func (m *manager) consumer(router Router, buffer <-chan Event) {
	for i := 0; i < router.Pool(); i++ {
		go func(r Router) {
			for e := range buffer {
				if err := r.Handler()(e); err != nil {
					continue
				}

				r.Commit(e)
			}
		}(router)
	}
}

func (m *manager) producer(router Router, buffer chan<- Event, wg *sync.WaitGroup) {
	router.Producer(buffer)
	wg.Done()
}

func (m *manager) Run() {
	wg := sync.WaitGroup{}
	for _, router := range m.routers {
		wg.Add(1)
		buffer := make(chan Event, router.Pool())
		go m.consumer(router, buffer)
		go m.producer(router, buffer, &wg)
	}

	wg.Wait()
}
