package expr

type Expression struct {
	Language   string // "simple", "constant",...
	Expression any    // 'string' for simple, 'any' for constant
}

func Simple(expression string) Expression {
	return Expression{
		Language:   "simple",
		Expression: expression,
	}
}

func Constant(value any) Expression {
	return Expression{
		Language:   "constant",
		Expression: value,
	}
}
