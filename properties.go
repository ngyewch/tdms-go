package tdms

import (
	"fmt"
	"time"
)

func getInt(props map[string]any, name string) (int, bool, error) {
	v, exists := props[name]
	if !exists {
		return 0, false, nil
	}
	n, err := asInt(v)
	if err != nil {
		return 0, true, err
	}
	return n, true, nil
}

func getFloat64(props map[string]any, name string) (float64, bool, error) {
	v, exists := props[name]
	if !exists {
		return 0, false, nil
	}
	n, err := asFloat64(v)
	if err != nil {
		return 0, true, err
	}
	return n, true, nil
}

func getString(props map[string]any, name string) (string, bool, error) {
	v, exists := props[name]
	if !exists {
		return "", false, nil
	}
	s, err := asString(v)
	if err != nil {
		return "", true, err
	}
	return s, true, nil
}

func getTime(props map[string]any, name string) (time.Time, bool, error) {
	v, exists := props[name]
	if !exists {
		return time.Time{}, false, nil
	}
	t, err := asTime(v)
	if err != nil {
		return time.Time{}, true, err
	}
	return t, true, nil
}

// ---

func asInt(v any) (int, error) {
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

func asFloat64(v any) (float64, error) {
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

func asString(v any) (string, error) {
	switch v1 := v.(type) {
	case string:
		return v1, nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", v)
	}
}

func asTime(v any) (time.Time, error) {
	switch v1 := v.(type) {
	case time.Time:
		return v1, nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %v to time.Time", v)
	}
}
