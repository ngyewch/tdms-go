package utils

import (
	"fmt"
	"time"
)

func AsInt(v any) (int, error) {
	switch v1 := v.(type) {
	case int:
		return v1, nil
	case uint:
		return int(v1), nil
	case int8:
		return int(v1), nil
	case uint8:
		return int(v1), nil
	case int16:
		return int(v1), nil
	case uint16:
		return int(v1), nil
	case int32:
		return int(v1), nil
	case uint32:
		return int(v1), nil
	case int64:
		return int(v1), nil
	case uint64:
		return int(v1), nil
	default:
		return 0, fmt.Errorf("cannot convert %v to int", v)
	}
}

func AsUint(v any) (uint, error) {
	switch v1 := v.(type) {
	case int:
		return uint(v1), nil
	case uint:
		return v1, nil
	case int8:
		return uint(v1), nil
	case uint8:
		return uint(v1), nil
	case int16:
		return uint(v1), nil
	case uint16:
		return uint(v1), nil
	case int32:
		return uint(v1), nil
	case uint32:
		return uint(v1), nil
	case int64:
		return uint(v1), nil
	case uint64:
		return uint(v1), nil
	default:
		return 0, fmt.Errorf("cannot convert %v to uint", v)
	}
}

func AsFloat64(v any) (float64, error) {
	switch v1 := v.(type) {
	case int:
		return float64(v1), nil
	case uint:
		return float64(v1), nil
	case int8:
		return float64(v1), nil
	case uint8:
		return float64(v1), nil
	case int16:
		return float64(v1), nil
	case uint16:
		return float64(v1), nil
	case int32:
		return float64(v1), nil
	case uint32:
		return float64(v1), nil
	case int64:
		return float64(v1), nil
	case uint64:
		return float64(v1), nil
	case float32:
		return float64(v1), nil
	case float64:
		return v1, nil
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", v)
	}
}

func AsString(v any) (string, error) {
	switch v1 := v.(type) {
	case string:
		return v1, nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", v)
	}
}

func AsTime(v any) (time.Time, error) {
	switch v1 := v.(type) {
	case time.Time:
		return v1, nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %v to time.Time", v)
	}
}
