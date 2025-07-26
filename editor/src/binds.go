package main

import (
	"regexp"
	"strings"
)

type Shortcuts map[string]string
type Strokes map[string]string

type Binds struct {
	Shortcuts Shortcuts
	Strokes   Strokes
}

func SanitizeBind(bind string) string {
	bind = strings.ReplaceAll(bind, "+", "-")
	bind = regexp.MustCompile(`-+`).ReplaceAllString(bind, "-")
	return strings.ToLower(bind)
}

func SanitizeBinds() {
	cleaned := make(Shortcuts)

	for k, v := range state.Binds.Shortcuts {
		sanitized := SanitizeBind(v)

		if IsCoreBind(sanitized) {
			LogErrF("%s is a reserved bind!", v)
			continue
		} else {

			if existing, ok := cleaned[sanitized]; ok && existing != v {
				LogErrF("Warning: bind collision after cleaning: %q and %q both map to %q\n", k, existing, sanitized)
			}

			cleaned[k] = sanitized
		}
	}

	state.Binds.Shortcuts = cleaned
}

func IsCoreBind(bind string) bool {
	_, exists := state.CoreBinds.Shortcuts[bind]
	return exists
}
