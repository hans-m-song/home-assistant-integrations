package util

import "fmt"

type ErrNotFound struct{ Key any }

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("key %v not found", e.Key)
}

type ErrWrongType struct{ Key, Expected, Actual any }

func (e ErrWrongType) Error() string {
	return fmt.Sprintf("key %v is not of type %T, got %T", e.Key, e.Expected, e.Actual)
}

func GetByKey[Value any, Key comparable](data map[Key]any, key Key) (Value, error) {
	zero := *new(Value)
	raw, ok := data[key]
	if !ok {
		return zero, ErrNotFound{key}
	}

	val, ok := raw.(Value)
	if !ok {
		return zero, ErrWrongType{key, zero, raw}
	}

	return val, nil
}
