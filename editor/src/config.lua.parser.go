package main

import (
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

type LuaCustomParser func(val lua.LValue) (reflect.Value, error)

var luaTypeParsers = map[reflect.Type]LuaCustomParser{}

func LuaWrapEnumParser[T ~int](enum EnumMap[T]) func(lua.LValue) (reflect.Value, error) {
	return func(val lua.LValue) (reflect.Value, error) {
		if val.Type() != lua.LTString {
			return reflect.Value{}, fmt.Errorf("expected string")
		}
		parsed, err := enum.Parse(val.String())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(parsed), nil
	}
}

func ConfigRegisterLuaTypes() {
	luaTypeParsers[reflect.TypeOf(IndentSpaces)] = LuaWrapEnumParser(IndentStyleEnum)
	luaTypeParsers[reflect.TypeOf(StatusbarTop)] = LuaWrapEnumParser(StatusbarPositionEnum)
}

func ApplyConfigLuaValue(target reflect.Value, luaVal lua.LValue) error {
	parser, ok := luaTypeParsers[target.Type()]
	if !ok {
		return fmt.Errorf("no Lua parser registered for type: %v", target.Type())
	}

	converted, err := parser(luaVal)
	if err != nil {
		return err
	}

	target.Set(converted)
	return nil
}
