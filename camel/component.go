package camel

type Consumer interface {
	Start() error
	Stop() error
}

type Producer interface {
	Processor
}

type Endpoint interface {
	Uri() string
	Consumer(processor Processor) (Consumer, error)
	Producer() (Producer, error)
}

type Messenger interface {
	SendMessage(message *Message) error
}

type Component interface {
	Id() string
	Endpoint(uri string) (Endpoint, error)
}

type ContextAware interface {
	SetContext(context *Context)
}
