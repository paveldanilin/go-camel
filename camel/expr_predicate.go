package camel

import (
	"fmt"
	"strconv"
)

func newPredicateFromExpr(expr Expr) PredicateFunc {
	return func(exchange *Exchange) (bool, error) {
		v, err := expr.Eval(exchange)
		if err != nil {
			return false, err
		}
		return toBool(v)
	}
}

func toBool(v any) (bool, error) {
	switch x := v.(type) {
	case nil:
		return false, nil
	case bool:
		return x, nil
	case int:
		return x != 0, nil
	case int8:
		return x != 0, nil
	case int16:
		return x != 0, nil
	case int32:
		return x != 0, nil
	case int64:
		return x != 0, nil
	case uint:
		return x != 0, nil
	case uint8:
		return x != 0, nil
	case uint16:
		return x != 0, nil
	case uint32:
		return x != 0, nil
	case uint64:
		return x != 0, nil
	case float32:
		return x != 0, nil
	case float64:
		return x != 0, nil
	case string:
		if b, err := strconv.ParseBool(x); err == nil {
			return b, nil
		}
		if x == "1" {
			return true, nil
		}
		if x == "0" {
			return false, nil
		}
		return false, fmt.Errorf("cannot convert string %q to bool", x)
	default:
		return false, fmt.Errorf("unsupported type %T", v)
	}
}
