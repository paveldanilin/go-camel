package camel

// setBodyProcessor sets a camel.Message body
type setBodyProcessor struct {
	id              string
	valueExpression expression
}

func newSetBodyProcessor(id string, valueExpression expression) *setBodyProcessor {
	return &setBodyProcessor{
		id:              id,
		valueExpression: valueExpression,
	}
}

func (p *setBodyProcessor) getId() string {
	return p.id
}

func (p *setBodyProcessor) Process(exchange *Exchange) {
	value, err := p.valueExpression.eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().Body = value
}
