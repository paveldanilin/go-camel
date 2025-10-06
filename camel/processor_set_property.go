package camel

// setPropertyProcessor sets Exchange's property
type setPropertyProcessor struct {
	id string

	name            string
	valueExpression expression
}

func newSetPropertyProcessor(id, name string, valueExpression expression) *setPropertyProcessor {
	return &setPropertyProcessor{
		id:              id,
		name:            name,
		valueExpression: valueExpression,
	}
}

func (p *setPropertyProcessor) getId() string {
	return p.id
}

func (p *setPropertyProcessor) Process(exchange *Exchange) {
	value, err := p.valueExpression.eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.SetProperty(p.name, value)
}
