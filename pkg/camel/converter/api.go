package converter

type Converter[From any, To any] interface {
	Convert(from From, params map[string]any) (To, error)
}
