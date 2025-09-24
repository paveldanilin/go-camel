package camel

// setPropertyProcessor sets Exchange's property
type setPropertyProcessor struct {
	id string

	name  string
	value Expr
}

func newSetPropertyProcessor(id, name string, value Expr) *setPropertyProcessor {
	return &setPropertyProcessor{
		id:    id,
		name:  name,
		value: value,
	}
}

func (p *setPropertyProcessor) getId() string {
	return p.id
}

func (p *setPropertyProcessor) Process(exchange *Exchange) {
	value, err := p.value.Eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.SetProperty(p.name, value)
}
