package numeric

import (
	"fmt"
	"strconv"
)

func Float64Value(value any) (float64, error) { //nolint:cyclop
	switch tValue := value.(type) {
	case int:
		return float64(tValue), nil
	case int8:
		return float64(tValue), nil
	case int16:
		return float64(tValue), nil
	case int32:
		return float64(tValue), nil
	case int64:
		return float64(tValue), nil
	case uint:
		return float64(tValue), nil
	case uint8:
		return float64(tValue), nil
	case uint16:
		return float64(tValue), nil
	case uint32:
		return float64(tValue), nil
	case uint64:
		return float64(tValue), nil
	case float32:
		return float64(tValue), nil
	case float64:
		return tValue, nil
	case string:
		floatValue, err := strconv.ParseFloat(tValue, 64)
		if err != nil {
			return 0, fmt.Errorf("parse float64: %w", err)
		}

		return floatValue, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}
