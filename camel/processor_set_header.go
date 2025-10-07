package camel

// setHeaderProcessor sets Message's header
type setHeaderProcessor struct {
	name string

	headerName      string
	valueExpression expression
}

func newSetHeaderProcessor(name, headerName string, valueExpression expression) *setHeaderProcessor {
	return &setHeaderProcessor{
		name:            name,
		headerName:      headerName,
		valueExpression: valueExpression,
	}
}

func (p *setHeaderProcessor) getName() string {
	return p.name
}

func (p *setHeaderProcessor) Process(exchange *Exchange) {
	value, err := p.valueExpression.eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.Message().SetHeader(p.headerName, value)
}
