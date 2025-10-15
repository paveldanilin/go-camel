package camel

type marshalProcessor struct {
	name       string
	dataFormat DataFormat
}

func newMarshalProcessor(name string, dataFormat DataFormat) *marshalProcessor {
	return &marshalProcessor{
		name:       name,
		dataFormat: dataFormat,
	}
}

func (p *marshalProcessor) getName() string {
	return p.name
}

func (p *marshalProcessor) Process(exchange *Exchange) {
	body, err := p.dataFormat.Marshal(exchange.Message().Body)
	if err != nil {
		exchange.SetError(err)
		return
	}
	exchange.Message().Body = body
}
