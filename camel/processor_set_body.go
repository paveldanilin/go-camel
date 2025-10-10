package camel

// setBodyProcessor sets a camel.Message body
type setBodyProcessor struct {
	name            string
	valueExpression expression
}

func newSetBodyProcessor(name string, valueExpression expression) *setBodyProcessor {
	return &setBodyProcessor{
		name:            name,
		valueExpression: valueExpression,
	}
}

func (p *setBodyProcessor) getName() string {
	return p.name
}

func (p *setBodyProcessor) Process(exchange *Exchange) {
	value, err := p.valueExpression.eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().Body = value
}
