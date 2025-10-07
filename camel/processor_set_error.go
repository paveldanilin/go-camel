package camel

// setErrorProcessor sets a camel.Exchange error
type setErrorProcessor struct {
	name string
	err  error
}

func newSetErrorProcessor(name string, err error) *setErrorProcessor {
	return &setErrorProcessor{
		name: name,
		err:  err,
	}
}

func (p *setErrorProcessor) getName() string {
	return p.name
}

func (p *setErrorProcessor) Process(exchange *Exchange) {
	exchange.SetError(p.err)
}
