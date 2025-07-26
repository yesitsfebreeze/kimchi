package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func GetFormattedTime() string {
	now := time.Now()
	return now.Format("15:04")
}

func SnakeToCamel(input string) string {
	var result = input
	if strings.ContainsAny(input, " ") {
		parts := strings.Split(input, " ")
		for i := range parts {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
			}
		}

		result = strings.Join(parts, "")
	}

	if strings.ContainsAny(input, "_") {
		parts := strings.Split(input, "_")
		for i := range parts {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
			}
		}

		result = strings.Join(parts, "")
	}

	return strings.ToUpper(result[:1]) + result[1:]
}

// TODO: maybe we reverse the maps that use this initially/onchange
// then this can be a constant lookup time
func KeyByValue(data map[string]string, search string) string {
	for k, v := range data {
		if v == search {
			return k
		}
	}
	return ""
}

type PathKind int

const (
	PathNone PathKind = iota
	PathFile
	PathDir
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func ClassifyPath(path string) (PathKind, string, error) {
	cs, err := filepath.Abs(path)
	if err != nil {
		return PathNone, "", err
	}

	if !FileExists(cs) {
		return PathNone, cs, nil
	}

	switch {
	case IsDir(cs):
		return PathDir, cs, nil
	case IsFile(cs):
		return PathFile, cs, nil
	default:
		return PathNone, cs, fmt.Errorf("unsupported path type: %s", cs)
	}
}

func FlattenStructKeys(v interface{}, prefix string) []string {
	var keys []string
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		name := field.Name
		fullKey := name
		if prefix != "" {
			fullKey = prefix + "." + name
		}

		if fieldVal.Kind() == reflect.Struct {
			keys = append(keys, FlattenStructKeys(fieldVal.Interface(), fullKey)...)
		} else {
			keys = append(keys, fullKey)
		}
	}

	return keys
}

func EnforceExtension(name string, ext string) string {
	if !strings.HasSuffix(name, ext) {
		return name + ext
	}
	return name
}

func ResolveFile(path string, basedir string) (string, error) {
	var result string = ""

	if FileExists(path) {
		return path, nil
	} else {

		if !filepath.IsAbs(path) {
			resolved := filepath.Join(basedir, path)
			if FileExists(resolved) {
				result = resolved
			}
		}

		if result == "" {
			abs, err := filepath.Abs(path)
			if err == nil && FileExists(abs) {
				result = abs
			}
		}
	}

	if result == "" {
		return "", fmt.Errorf(ErrFConfigFileDoesNotExist, path)
	}

	return result, nil
}

type NumberPrimitives interface {
	~int | ~int64 | ~float64 | ~float32 | ~uint | ~uint64
}

func Tern[T NumberPrimitives](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func Clamp[T NumberPrimitives](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
