package camel

type Processor interface {
	Process(message *Message) error
}
