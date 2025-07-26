package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	homedir "github.com/mitchellh/go-homedir"
	lua "github.com/yuin/gopher-lua"
)

const CONFIG_FILE_NAME = "kitsune.lua"

var ConfigLuaState *lua.LState
var ConfigCurrentFile string

func (c *Config) Dump() {
	spew.Dump(c)
}

func (c *Config) Get(key string) any {
	parts := strings.Split(key, ".")
	var current reflect.Value = reflect.ValueOf(c).Elem()

	for _, part := range parts {
		if current.Kind() == reflect.Struct {
			field := current.FieldByNameFunc(func(n string) bool {
				return strings.EqualFold(n, part)
			})

			if !field.IsValid() {
				return nil
			}

			current = field

		} else if current.Kind() == reflect.Map {
			mapKey := reflect.ValueOf(part)
			val := current.MapIndex(mapKey)
			if !val.IsValid() {
				return nil
			}
			current = val

		} else {
			return nil // unsupported nesting
		}
	}

	if !current.IsValid() || !current.CanInterface() {
		return nil
	}

	return current.Interface()
}

func ConfigThrowErr(file string, line string, msg string) {
	fmt.Fprintln(os.Stderr, "\033[31mConfig Error:\033[0m")

	lineNum, err := strconv.Atoi(line)
	if err != nil {
		lineNum = 0
	}
	// Read source line
	srcLine := ""
	data, err := os.ReadFile(file)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		if lineNum-1 < len(lines) {
			srcLine = lines[lineNum-1]
		}
	}

	// Print location
	fmt.Fprintf(os.Stderr, "%s:%d\n", file, lineNum)

	// Print source line
	if srcLine != "" {
		fmt.Fprintf(os.Stderr, "  %2d | %s\n", lineNum, srcLine)
		fmt.Fprintf(os.Stderr, "     | \033[31m^ %s\033[0m\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "  \033[31m%s\033[0m\n", msg)
	}

	os.Exit(1)
}

func ConfigRunLuaFile(filename string) {
	if !FileExists(filename) {
		ThrowErrF(ErrFConfigFileDoesNotExist, filename)
	}

	ConfigCurrentFile = filename
	if err := ConfigLuaState.DoFile(filename); err != nil {
		msg := err.Error()
		if idx := strings.IndexRune(msg, '\n'); idx != -1 {
			msg = msg[:idx]
		}
		parts := strings.Split(msg, ":")

		ConfigThrowErr(
			strings.Trim(parts[0], " "),
			strings.Trim(parts[1], " "),
			strings.Trim(parts[2], " "),
		)
	}
}

func ConfigInit() {
	ConfigMeta()
	ConfigRegisterLuaTypes()

	ConfigLuaState = lua.NewState()
	ConfigLoadLuaApi(ConfigLuaState)
	defer ConfigLuaState.Close()

	state.Config = DefaultConfig()

	config_file, err := ConfigResolve(CONFIG_FILE_NAME, false)

	if err == nil {
		Log("User-Config:", config_file)
		state.ConfigFile = config_file
		ConfigRunLuaFile(state.ConfigFile)
	}

	if state.Args.Path != "" {
		ConfigParseInputPath(state.Args.Path)
	} else {
		if cwd, err := os.Getwd(); err == nil {
			ConfigParseInputPath(cwd)
		}
	}

	if state.ProjectFile != "" {
		Log("Project-Config:", state.ProjectFile)
		ConfigRunLuaFile(state.ProjectFile)
	}

	SanitizeBinds()

	ConfigCurrentFile = ""
}

func ConfigSet(path string, val lua.LValue) error {
	parts := strings.Split(path, ".")
	for i := range parts {
		parts[i] = SnakeToCamel(parts[i])
	}

	path = strings.Join(parts, ".")

	if !AllowedConfigKeys[path] {
		ThrowErrF(ErrFUnkownConfigKey, path)
	}

	v := reflect.ValueOf(&state.Config).Elem()

	for i := range parts {
		part := parts[i]

		if i == len(parts)-2 {
			field := v.FieldByNameFunc(func(n string) bool {
				return strings.EqualFold(n, part)
			})
			if field.IsValid() && field.Kind() == reflect.Map {
				key := reflect.ValueOf(parts[i+1])
				valv, err := ConfigConvertValue(val, field.Type().Elem())
				if err != nil {
					return err
				}
				field.SetMapIndex(key, valv)
				return nil
			}
		}

		field := v.FieldByNameFunc(func(n string) bool {
			return strings.EqualFold(n, part)
		})

		if !field.IsValid() {
			return ErrF(ErrFInvalidConfigField, parts[i])
		}

		if i == len(parts)-1 {
			// Set leaf value
			valv, err := ConfigConvertValue(val, field.Type())
			if err != nil {
				return err
			}
			field.Set(valv)
			return nil
		}

		// Descend into struct
		if field.Kind() == reflect.Struct {
			v = field
		} else if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			v = field.Elem()
		} else {
			return fmt.Errorf("unexpected type at %s", part)
		}
	}

	return nil
}

func ConfigConvertValue(val lua.LValue, t reflect.Type) (reflect.Value, error) {
	// First: try custom type parser
	if parser, ok := luaTypeParsers[t]; ok {
		return parser(val)
	}

	// Then fallback to basic kinds
	switch t.Kind() {
	case reflect.String:
		if val.Type() != lua.LTString {
			return reflect.Value{}, ErrF(ErrFUnkownLuaType, "expected "+val.Type().String())
		}
		return reflect.ValueOf(val.String()).Convert(t), nil

	case reflect.Int:
		if n, ok := val.(lua.LNumber); ok {
			return reflect.ValueOf(int(n)).Convert(t), nil
		}
		return reflect.Value{}, ErrF(ErrFUnkownLuaType, "expected "+val.Type().String())

	case reflect.Bool:
		if b, ok := val.(lua.LBool); ok {
			return reflect.ValueOf(bool(b)).Convert(t), nil
		}
		return reflect.Value{}, ErrF(ErrFUnkownLuaType, "expected "+val.Type().String())
	}

	return reflect.Value{}, ErrF(ErrFUnkownLuaType, "expected "+val.Type().String())
}

func ConfigResolve(name string, is_full_path bool) (string, error) {
	filename := EnforceExtension(name, ".lua")
	file := ""

	if is_full_path {

		if FileExists(filename) {
			return filename, nil
		}

		return "", os.ErrNotExist
	} else {
		home_dir, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		file = filepath.Join(home_dir, ".config/kitsune", filename)
		if FileExists(file) {
			return file, nil
		}

		exe_path, err := os.Executable()
		if err != nil {
			return "", err
		}

		file = filepath.Join(filepath.Dir(exe_path), CONFIG_FILE_NAME)
		if FileExists(file) {
			return file, nil
		}
	}

	return "", os.ErrNotExist
}

func ConfigParseInputPath(startpath string) {

	sp, err := filepath.Abs(startpath)
	if err != nil {
		LogErr("Inputfile is corrupt:", startpath)
	}

	if !FileExists(sp) {
		LogErr("Inputfile doesn't exist:", sp)
		return
	}
	pathtype, path, err := ClassifyPath(sp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	if pathtype == PathKind(PathFile) {
		state.InputFile = path
		state.Project = filepath.Dir(path)
	}

	if pathtype == PathKind(PathDir) {
		state.Project = path
		state.InputFile = ""
	}

	projectFile := FindClosestProjectFile()

	if projectFile != "" {
		state.ProjectFile = projectFile
		state.Project = filepath.Dir(projectFile)
	}
}

func FindClosestProjectFile() string {
	start := state.Project
	limit := state.Config.ProjectScanLimit

	path := start
	for i := 0; i <= limit; i++ {
		cfg_path := filepath.Join(path, CONFIG_FILE_NAME)
		if FileExists(cfg_path) {
			return cfg_path
		}

		parent := filepath.Dir(path)
		if parent == path {
			// reached filesystem root
			break
		}
		path = parent
	}

	// not found
	return ""

}
