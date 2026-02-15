package converter

type Converter[From any, To any] interface {
	Convert(from From, params map[string]any) (To, error)
}

type Func[From any, To any] func(from From, params map[string]any) (To, error)

func (f Func[From, To]) Convert(from From, params map[string]any) (To, error) {
	return f(from, params)
}
