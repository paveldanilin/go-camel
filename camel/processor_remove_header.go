package camel

type removeHeaderProcessor struct {
	name string

	headerName string
}

func newRemoveHeaderProcessor(name, headerName string) *removeHeaderProcessor {
	return &removeHeaderProcessor{
		name:       name,
		headerName: headerName,
	}
}

func (p *removeHeaderProcessor) getName() string {
	return p.name
}

func (p *removeHeaderProcessor) Process(exchange *Exchange) {
	exchange.Message().RemoveHeader(p.headerName)
}
