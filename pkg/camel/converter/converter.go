package converter

import (
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"strconv"
	"time"
)

type Func[From any, To any] func(from From, params map[string]any) (To, error)

func (f Func[From, To]) Convert(from From, params map[string]any) (To, error) {
	return f(from, params)
}

// String -> int, float, bool, bool

func StringToInt() api.Converter[string, int] {
	return Func[string, int](func(from string, _ map[string]any) (int, error) {
		return strconv.Atoi(from)
	})
}

func StringToInt64() api.Converter[string, int64] {
	return Func[string, int64](func(from string, _ map[string]any) (int64, error) {
		return strconv.ParseInt(from, 10, 64)
	})
}

func StringToFloat64() api.Converter[string, float64] {
	return Func[string, float64](func(from string, _ map[string]any) (float64, error) {
		return strconv.ParseFloat(from, 64)
	})
}

func StringToFloat() api.Converter[string, float32] {
	return Func[string, float32](func(from string, _ map[string]any) (float32, error) {
		v, err := strconv.ParseFloat(from, 32)
		if err != nil {
			return 0, err
		}
		return float32(v), nil
	})
}

func StringToBool() api.Converter[string, bool] {
	return Func[string, bool](func(from string, _ map[string]any) (bool, error) {
		return strconv.ParseBool(from)
	})
}

func StringToDateTime() api.Converter[string, time.Time] {
	return Func[string, time.Time](func(from string, params map[string]any) (time.Time, error) {
		layout, ok := params["layout"].(string)
		if !ok {
			layout = "2006-01-02 15:04:05" // Default
		}

		t, err := time.Parse(layout, from)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	})
}
