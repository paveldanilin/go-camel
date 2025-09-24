package camel

type funcProcessor struct {
	id       string
	userFunc func(*Exchange)
}

func newFuncProcessor(id string, userFunc func(*Exchange)) *funcProcessor {
	return &funcProcessor{
		id:       id,
		userFunc: userFunc,
	}
}

func (p *funcProcessor) getId() string {
	return p.id
}

func (p *funcProcessor) Process(exchange *Exchange) {
	p.userFunc(exchange)
}
