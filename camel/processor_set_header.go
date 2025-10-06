package camel

// setHeaderProcessor sets Message's header
type setHeaderProcessor struct {
	id string

	name            string
	valueExpression expression
}

func newSetHeaderProcessor(id, name string, valueExpression expression) *setHeaderProcessor {
	return &setHeaderProcessor{
		id:              id,
		name:            name,
		valueExpression: valueExpression,
	}
}

func (p *setHeaderProcessor) getId() string {
	return p.id
}

func (p *setHeaderProcessor) Process(exchange *Exchange) {
	value, err := p.valueExpression.eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().SetHeader(p.name, value)
}
