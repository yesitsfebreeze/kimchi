package main

import (
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type IndentCfg struct {
	style string // "tabs" or "spaces"
	width int
}

type IndentSettings struct {
	visual IndentCfg
	save   IndentCfg
}

type Cfg struct {
	project_scan_limit  int
	startfile           string
	workspace           string
	leader_key          string
	autosave_on_close   bool
	status_bar_enabled  bool
	status_bar_position string
	indent              struct {
		visual IndentCfg
		save   IndentCfg
	}
	bindings map[string]string // binding -> command,
	strokes  map[string]string // stroke -> command,
}

func find_cfg(name string) (string, error) {
	name = strings.TrimSuffix(name, ".lua")

	exe_path, err := os.Executable()
	if err != nil {
		return "", err
	}

	exe_dir := filepath.Dir(exe_path)

	// 1. Check local ./cfg/init.lua
	local, err := filepath.Abs(filepath.Join(exe_dir, "..", "cfg", name+".lua"))
	if err != nil {
		return "", err
	}

	if file_exists(local) {
		return local, nil
	}

	// 2. Check $XDG_cfg_HOME/kimchi/init.lua
	xdg := os.Getenv("XDG_cfg_HOME")
	if xdg != "" {
		xdgPath := filepath.Join(xdg, "kimchi", name+".lua")
		if file_exists(xdgPath) {
			return xdgPath, nil
		}
	}

	// 3. Fallback to $HOME/.cfg/kimchi/init.lua
	home, err := os.UserHomeDir()
	if err == nil {
		fallback := filepath.Join(home, ".config", "kimchi", name+".lua")
		if file_exists(fallback) {
			return fallback, nil
		}
	}

	return "", os.ErrNotExist
}

func load_cfg() error {
	path, err := find_cfg("init")

	if err != nil {
		return err
	}

	L := lua.NewState()
	defer L.Close()

	set_cfg_api(L)
	parse_cfg(L)

	// Run the cfg file
	if err := L.DoFile(path); err != nil {
		return err
	}

	return nil
}

func set_cfg_api(L *lua.LState) {
	L.SetGlobal("load", L.NewFunction(load))
	L.SetGlobal("echo", L.NewFunction(echo))
	L.SetGlobal("bind", L.NewFunction(bind))
	L.SetGlobal("stroke", L.NewFunction(stroke))
}

func parse_cfg(L *lua.LState) {
	get_string := func(name string, target *string) {
		if val := L.GetGlobal(name); val.Type() == lua.LTString {
			*target = val.String()
		}
	}
	get_int := func(name string, target *int) {
		if val := L.GetGlobal(name); val.Type() == lua.LTNumber {
			*target = int(val.(lua.LNumber))
		}
	}
	get_bool := func(name string, target *bool) {
		if val := L.GetGlobal(name); val.Type() == lua.LTBool {
			*target = bool(val.(lua.LBool))
		}
	}

	get_int("project_scan_limit", &cfg.project_scan_limit)
	get_string("leader_combo", &cfg.leader_key)
	get_bool("autosave_on_close", &cfg.autosave_on_close)
	get_int("save_indent_width", &cfg.indent.save.width)
	get_string("save_indent_style", &cfg.indent.save.style)
	get_int("visual_indent_width", &cfg.indent.visual.width)
	get_string("visual_indent_style", &cfg.indent.visual.style)
}

// Lua example: load("file")
func load(L *lua.LState) int {
	log(L.ToString(1))
	return 0
}

// Lua example: echo("my mesasge")
func echo(L *lua.LState) int {
	log(L.ToString(1))
	return 0
}

// Lua example: bind("C-s, "save_buffer")
func bind(L *lua.LState) int {
	seq := L.ToString(1)
	cmd := L.ToString(2)
	cfg.bindings[seq] = cmd

	log("Binding set:", seq, "->", cmd)
	return 0
}

// Lua example: bind_stroke("ff", "fuzzy_find")
func stroke(L *lua.LState) int {
	seq := L.ToString(1)
	cmd := L.ToString(2)
	cfg.bindings[seq] = cmd

	log("Binding set:", seq, "->", cmd)
	return 0
}
