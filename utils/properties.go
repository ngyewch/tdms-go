package utils

import (
	"time"
)

func GetInt(props map[string]any, name string) (int, bool, error) {
	v, exists := props[name]
	if !exists {
		return 0, false, nil
	}
	n, err := AsInt(v)
	if err != nil {
		return 0, true, err
	}
	return n, true, nil
}

func GetUint(props map[string]any, name string) (uint, bool, error) {
	v, exists := props[name]
	if !exists {
		return 0, false, nil
	}
	n, err := AsUint(v)
	if err != nil {
		return 0, true, err
	}
	return n, true, nil
}

func GetFloat64(props map[string]any, name string) (float64, bool, error) {
	v, exists := props[name]
	if !exists {
		return 0, false, nil
	}
	n, err := AsFloat64(v)
	if err != nil {
		return 0, true, err
	}
	return n, true, nil
}

func GetString(props map[string]any, name string) (string, bool, error) {
	v, exists := props[name]
	if !exists {
		return "", false, nil
	}
	s, err := AsString(v)
	if err != nil {
		return "", true, err
	}
	return s, true, nil
}

func GetTime(props map[string]any, name string) (time.Time, bool, error) {
	v, exists := props[name]
	if !exists {
		return time.Time{}, false, nil
	}
	t, err := AsTime(v)
	if err != nil {
		return time.Time{}, true, err
	}
	return t, true, nil
}
