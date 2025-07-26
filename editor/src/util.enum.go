package main

import (
	"fmt"
	"strings"
)

type EnumMap[T ~int] struct {
	strToVal map[string]T
	valToStr map[T]string
}

func NewEnumMap[T ~int](m map[string]T) EnumMap[T] {
	reverse := make(map[T]string)
	for k, v := range m {
		reverse[v] = k
	}
	return EnumMap[T]{strToVal: m, valToStr: reverse}
}

func (e EnumMap[T]) Parse(s string) (T, error) {
	val, ok := e.strToVal[strings.ToLower(s)]
	if !ok {
		var zero T
		return zero, fmt.Errorf("invalid value: %s", s)
	}
	return val, nil
}

func (e EnumMap[T]) String(val T) string {
	str, ok := e.valToStr[val]
	if !ok {
		return "unknown"
	}
	return str
}
