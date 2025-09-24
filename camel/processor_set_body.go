package camel

// setBodyProcessor sets a camel.Message body
type setBodyProcessor struct {
	id    string
	value Expr
}

func newSetBodyProcessor(id string, value Expr) *setBodyProcessor {
	return &setBodyProcessor{
		id:    id,
		value: value,
	}
}

func (p *setBodyProcessor) getId() string {
	return p.id
}

func (p *setBodyProcessor) Process(exchange *Exchange) {
	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().Body = value
}
