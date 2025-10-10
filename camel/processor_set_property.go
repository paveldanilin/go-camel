package camel

// setPropertyProcessor sets Exchange's property
type setPropertyProcessor struct {
	name string

	propertyName    string
	valueExpression expression
}

func newSetPropertyProcessor(name, propertyName string, valueExpression expression) *setPropertyProcessor {
	return &setPropertyProcessor{
		name:            name,
		propertyName:    propertyName,
		valueExpression: valueExpression,
	}
}

func (p *setPropertyProcessor) getName() string {
	return p.name
}

func (p *setPropertyProcessor) Process(exchange *Exchange) {
	value, err := p.valueExpression.eval(exchange)
	if err != nil {
		exchange.SetError(err)
		return
	}

	exchange.SetProperty(p.propertyName, value)
}
