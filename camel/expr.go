package camel

type Expr interface {
	Eval(message *Message) (any, error)
}
