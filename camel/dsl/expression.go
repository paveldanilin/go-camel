package dsl

type Expression struct {
	Language   string // "simple", "constant",...
	Expression string // Language != "constant"
	Value      any
}

func Simple(expression string) Expression {
	return Expression{
		Language:   "simple",
		Expression: expression,
	}
}

func Constant(value any) Expression {
	return Expression{
		Language: "constant",
		Value:    value,
	}
}
