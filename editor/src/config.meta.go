package main

import (
	"reflect"
	"strings"
)

var AllowedConfigKeys = map[string]bool{}
var ConfigMetaData map[string]string

func ConfigExtractMetaData(key string) string {
	v := reflect.ValueOf(state.Config)
	t := v.Type()
	return ConfigExtractMetaRecursive(v, t, key, "")
}

func ConfigExtractMetaRecursive(v reflect.Value, t reflect.Type, targetKey string, prefix string) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		key := strings.ToLower(field.Name)

		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if fullKey == targetKey {
			return field.Tag.Get("help")
		}

		if field.Type.Kind() == reflect.Struct {
			if meta := ConfigExtractMetaRecursive(fieldVal, field.Type, targetKey, fullKey); meta != "" {
				return meta
			}
		}
	}
	return ""
}

func ConfigMeta() {
	ConfigMetaData = make(map[string]string)
	for _, key := range FlattenStructKeys(state.Config, "") {
		AllowedConfigKeys[key] = true
		lwr := strings.ToLower(key)
		ConfigMetaData[lwr] = ConfigExtractMetaData(lwr) // << string
	}
}
