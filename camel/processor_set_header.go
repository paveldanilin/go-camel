package camel

// setHeaderProcessor sets Message's header
type setHeaderProcessor struct {
	id string

	name  string
	value Expr
}

func newSetHeaderProcessor(id, name string, value Expr) *setHeaderProcessor {
	return &setHeaderProcessor{
		id:    id,
		name:  name,
		value: value,
	}
}

func (p *setHeaderProcessor) getId() string {
	return p.id
}

func (p *setHeaderProcessor) Process(exchange *Exchange) {
	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().SetHeader(p.name, value)
}
