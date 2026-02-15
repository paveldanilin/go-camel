package converter

import "strconv"

func IntToString() Converter[int, string] {
	return Func[int, string](func(from int, _ map[string]any) (string, error) {
		return strconv.Itoa(from), nil
	})
}
