package camel

// setErrorProcessor sets a camel.Exchange error
type setErrorProcessor struct {
	id  string
	err error
}

func newSetErrorProcessor(id string, err error) *setErrorProcessor {
	return &setErrorProcessor{
		id:  id,
		err: err,
	}
}

func (p *setErrorProcessor) getId() string {
	return p.id
}

func (p *setErrorProcessor) Process(exchange *Exchange) {
	exchange.SetError(p.err)
}
