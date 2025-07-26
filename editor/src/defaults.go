package main

func DefaultConfig() Config {
	return Config{
		// FPS:              60,
		FPS:              0,
		ProjectScanLimit: 10,
		XTerm:            true,
		AutosaveOnClose:  true,
		MouseEnabled:     true,
		CopyPaste:        true,
		SurroundingLines: 0,
		UnfocusedDarken:  0.5,
		CursorTrail: CursorTrailConfig{
			Enabled: false,
			Time:    100, // milliseconds
			Length:  8,   // number of trail cells
		},
		PromptPosition: PromptTop,

		Statusbar: StatusbarConfig{
			Enabled:      true,
			UseSeparator: true,
			Separator:    "─", // default separator
			Position:     StatusbarBottom,
			Left:         "{git} {file} {pos} ({cursor_count})",
			Center:       "",
			Right:        "{type} {format} | {time}",
		},
		Indent: IndentConfig{
			Visual: IndentOptions{
				Style: IndentTabs,
				Width: 2,
			},
			Save: IndentOptions{
				Style: IndentSpaces,
				Width: 4,
			},
		},
		WhiteSpace: WhiteSpaceConfig{
			Space: '•',
			Tab:   '→',
			Eol:   '¬',
		},
		Plugins: map[string]string{},
	}
}

func CoreBinds() Binds {
	return Binds{
		Shortcuts: map[string]string{
			"CursorUp":    "up",
			"CursorDown":  "down",
			"CursorLeft":  "left",
			"CursorRight": "right",
			"CursorStart": "home",
			"CursorEnd":   "end",

			"DelLeft":   "backspace",
			"DelRight":  "delete",
			"LineBreak": "enter",
		},
		Strokes: map[string]string{
			"ff": "FuzzyFind",
			"ll": "ToggleLog",
		},
	}

}

func DefaultBinds() Binds {
	return Binds{
		Shortcuts: map[string]string{
			"Prompt":          "ctrl-space",
			"Quit":            "ctrl-q",
			"SaveBuffer":      "ctrl-s",
			"SaveAllBuffers":  "ctrl-shift-s",
			"CloseBuffer":     "ctrl-w",
			"CloseAllBuffers": "ctrl-shift-w",
			"GotoLine":        "ctrl-l",
			"Find":            "ctrl-f",
			"SelectNext":      "ctrl-d",
		},
		Strokes: map[string]string{
			"FuzzyFind": "ff",
			"ToggleLog": "ll",
		},
	}
}
