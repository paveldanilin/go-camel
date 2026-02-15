package converter

import "strconv"

const (
	// FloatFormat is the format fmt is one of
	//  - 'b' (-ddddp±ddd, a binary exponent),
	//  - 'e' (-d.dddde±dd, a decimal exponent),
	//  - 'E' (-d.ddddE±dd, a decimal exponent),
	//  - 'f' (-ddd.dddd, no exponent),
	//  - 'g' ('e' for large exponents, 'f' otherwise),
	//  - 'G' ('E' for large exponents, 'f' otherwise),
	//  - 'x' (-0xd.ddddp±ddd, a hexadecimal fraction and binary exponent), or
	//  - 'X' (-0Xd.ddddP±ddd, a hexadecimal fraction and binary exponent).
	FloatFormat = "format"

	// FloatPrecision is the precision prec controls the number of digits (excluding the exponent) printed by the
	// 'e', 'E', 'f', 'g', 'G', 'x', and 'X' formats. For 'e', 'E', 'f', 'x', and 'X',
	// it is the number of digits after the decimal point.
	// For 'g' and 'G' it is the maximum number of significant digits (trailing zeros are removed).
	// The special precision -1 uses the smallest number of digits necessary such that ParseFloat will return f exactly.
	// The exponent is written as a decimal integer; for all formats other than 'b', it will be at least two digits.
	FloatPrecision = "precision"
	FloatBitSize   = "bitsize"
)

func Float64ToString() Converter[float64, string] {
	return Func[float64, string](func(from float64, params map[string]any) (string, error) {
		format := byte('g')
		precision := -1
		bitSize := 64

		if v, exists := params[FloatFormat]; exists {
			format = v.(byte)
		}

		if v, exists := params[FloatPrecision]; exists {
			precision = v.(int)
		}

		if v, exists := params[FloatBitSize]; exists {
			bitSize = v.(int)
		}

		return strconv.FormatFloat(from, format, precision, bitSize), nil
	})
}
