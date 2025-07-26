package main

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func ConfigLoadLuaApi(L *lua.LState) {
	api := map[string]func(*lua.LState) int{
		"load":         LuaLoad,
		"echo":         LuaEcho,
		"cfg":          LuaCfg,
		"clear_cfg":    LuaClearCfg,
		"bind":         LuaBind,
		"clear_bind":   LuaClearBind,
		"stroke":       LuaStroke,
		"clear_stroke": LuaClearStroke,
		"style":        LuaStyle,
		"clear_style":  LuaStyle,
	}

	for name, fn := range api {
		L.SetGlobal(name, L.NewFunction(fn))
	}
}

// load("file") // relative to current file
func LuaLoad(L *lua.LState) int {
	filename := EnforceExtension(L.CheckString(1), ".lua")

	basedir := filepath.Dir(ConfigCurrentFile)
	path, err := ResolveFile(filename, basedir)
	if err != nil {
		L.RaiseError("%s", err.Error())
	}

	ConfigRunLuaFile(path)
	return 0
}

// echo("my mesasge")
func LuaEcho(L *lua.LState) int {
	Log(L.ToString(1))
	return 0
}

// cfg('key', 'value')
func LuaCfg(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckAny(2)

	err := ConfigSet(key, val)
	if err != nil {
		L.RaiseError("should be: %s", ConfigMetaData[key])
	}

	return 1
}

// clear_cfg('key')
func LuaClearCfg(L *lua.LState) int {
	key := L.CheckString(1)
	val := state.DefaultConfig.Get(key)

	luaVal, ok := val.(lua.LValue)
	if !ok {
		L.RaiseError("value for key '%s' is not a Lua value", key)
		return 0
	}

	err := ConfigSet(key, luaVal)
	if err != nil {
		L.RaiseError("should be: %s", ConfigMetaData[key])
	}

	return 1
}

// style("key", "#foreground", "#background", <style>) // relative to current file
func LuaStyle(L *lua.LState) int {
	key := L.CheckString(1)
	foreground := L.CheckString(2)
	background := L.OptString(3, "")
	style := L.OptString(4, "")

	state.Theme.SetStyle(key, foreground, background, style)

	// basedir := filepath.Dir(ConfigCurrentFile)
	// path, err := ResolveFile(filename, basedir)
	// if err != nil {
	// 	L.RaiseError("%s", err.Error())
	// }

	// ConfigRunLuaFile(path)
	return 0
}

// TODO: safety check
// bind("name", "key-combo")
func LuaBind(L *lua.LState) int {
	action := L.CheckString(1)
	bind := L.CheckString(2)
	state.Binds.Shortcuts[action] = bind
	return 0
}

// TODO: safety check
// unbind("ctrl-s")
func LuaClearBind(L *lua.LState) int {
	seq := L.CheckString(1)
	key := KeyByValue(state.Binds.Shortcuts, seq)
	delete(state.Binds.Shortcuts, key)
	return 0
}

// TODO: safety check
// stroke("fuzzy_find", "ff")
func LuaStroke(L *lua.LState) int {
	action := L.CheckString(1)
	seq := L.CheckString(2)
	state.Binds.Strokes[action] = seq
	return 0
}

// TODO: safety check
// destroke("ff")
func LuaClearStroke(L *lua.LState) int {
	seq := L.CheckString(1)
	key := KeyByValue(state.Binds.Strokes, seq)
	delete(state.Binds.Strokes, key)
	return 0
}
