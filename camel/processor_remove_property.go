package camel

type removePropertyProcessor struct {
	name string

	propertyName string
}

func newRemovePropertyProcessor(name, propertyName string) *removePropertyProcessor {
	return &removePropertyProcessor{
		name:         name,
		propertyName: propertyName,
	}
}

func (p *removePropertyProcessor) getName() string {
	return p.name
}

func (p *removePropertyProcessor) Process(exchange *Exchange) {
	exchange.RemoveProperty(p.propertyName)
}
