package mw

type manager struct {
	routers  []Router
	poolSize int
}

func New() *manager {
	return &manager{
		routers: []Router{},
	}
}

func (m *manager) Register(r Router) {
	m.poolSize += r.Pool()
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

func (m *manager) producer(router Router, buffer chan<- Event) {
	go router.Producer(buffer)
}

func (m *manager) Run() {
	for _, router := range m.routers {
		buffer := make(chan Event, m.poolSize)
		go m.consumer(router, buffer)
		go m.producer(router, buffer)
	}

	for {
	}
}
