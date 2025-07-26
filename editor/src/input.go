package main

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

func GetModPrefix(mods tcell.ModMask) string {
	prefix := ""
	if mods&tcell.ModShift != 0 {
		prefix += "shift-"
	}
	if mods&tcell.ModCtrl != 0 {
		prefix += "ctrl-"
	}
	if mods&tcell.ModAlt != 0 {
		prefix += "alt-"
	}
	if mods&tcell.ModMeta != 0 {
		prefix += "meta-"
	}

	return prefix
}

func TranslateInput(ev *tcell.EventKey) string {
	prefix := GetModPrefix(ev.Modifiers())

	// remove mod
	re := regexp.MustCompile(`[^+]+$`)
	bind := re.FindString(ev.Name())

	if ev.Key() == tcell.KeyBackspace || ev.Rune() == '\b' || ev.Rune() == 127 {
		bind = "backspace"
	} else {
		// strip rune[] away
		re = regexp.MustCompile(`\[(.*?)\]`)
		match := re.FindStringSubmatch(bind)
		if len(match) > 1 {
			bind = match[1]
			r := []rune(bind)
			if len(r) == 1 && unicode.IsUpper(r[0]) {
				prefix = strings.ReplaceAll(prefix, "shift-", "")
				prefix += "shift-"
			}
			bind = strings.ToLower(bind)
		}
	}

	return strings.ToLower(prefix + bind)
}

func TryActionExecute(binds map[string]string, bind string) bool {
	action_name := KeyByValue(binds, bind)

	if action_name == "" {
		return false
	}

	action_name = SnakeToCamel(action_name)
	action, ok := ActionList[action_name]
	if !ok {
		LogErrF("Action '%s' not found", action_name)
		return false
	}
	if action != nil {
		if DEBUG {
			Log("executing:", action_name)
		}
		action()
	} else {
		if DEBUG {
			Log("Executing (Not Implemented):", action_name)
		}
	}
	return true
}

func IsValidtRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func HandlePrompt(bind string, r rune) bool {
	if !state.Prompt.Visible {
		return false
	}

	// TODO: also use ctrl-q(Quit)
	if bind == "esc" {
		ClosePrompt()
		return true
	}

	if bind == "enter" {
		SubmitPrompt()
		ClearPrompt()
		ClosePrompt()
		return true
	}

	return false

	// if IsValidtRune(r) {
	// 	SlateInput(r)
	// }
	// return true
}

func HandleAction(bind string) bool {
	if TryActionExecute(state.CoreBinds.Shortcuts, bind) {
		return true
	}
	if TryActionExecute(state.Binds.Shortcuts, bind) {
		return true
	}

	return false
}

func HandleInput(ev *tcell.EventKey) {
	bind := SanitizeBind(TranslateInput(ev))
	r := ev.Rune()

	if DEBUG {
		Log("pressed ", bind)
	}

	if HandlePrompt(bind, r) {
		return
	}

	if HandleAction(bind) {
		return
	}

	// Only insert if rune is printable and no modifiers except possibly Shift
	if r != 0 && (ev.Modifiers() == 0 || ev.Modifiers() == tcell.ModShift) {
		EditInsertRune(r)
	}
}

func HandleMouseInput(ev *tcell.EventMouse) {
	buttons := ev.Buttons()
	if buttons&tcell.WheelUp != 0 {
		WithActiveArea(func(a *Area) { a.Scroll(0, -1) })
	}
	if buttons&tcell.WheelDown != 0 {
		WithActiveArea(func(a *Area) { a.Scroll(0, 1) })
	}
}
