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

func (m *manager) consumer(buffer <-chan Event) {
	for _, router := range m.routers {
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
}

func (m *manager) producer(buffer chan<- Event) {
	for _, router := range m.routers {
		go router.Producer(buffer)
	}
}

func (m *manager) Run() {
	buffer := make(chan Event, m.poolSize)
	m.consumer(buffer)
	m.producer(buffer)

	for {
	}
}
