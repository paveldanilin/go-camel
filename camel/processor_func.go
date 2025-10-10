package camel

type funcProcessor struct {
	name     string
	userFunc func(*Exchange)
}

func newFuncProcessor(name string, userFunc func(*Exchange)) *funcProcessor {
	return &funcProcessor{
		name:     name,
		userFunc: userFunc,
	}
}

func (p *funcProcessor) getName() string {
	return p.name
}

func (p *funcProcessor) Process(exchange *Exchange) {
	p.userFunc(exchange)
}
