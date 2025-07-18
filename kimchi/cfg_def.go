package main

func create_default_config() Cfg {
	return Cfg{
		project_scan_limit:  10,
		startfile:           "",
		workspace:           "",
		leader_key:          "C-Space",
		autosave_on_close:   true,
		status_bar_enabled:  true,
		status_bar_position: "top", // or bottom
		indent: IndentSettings{
			visual: IndentCfg{
				style: "spaces",
				width: 4,
			},
			save: IndentCfg{
				style: "spaces",
				width: 4,
			},
		},
		bindings: map[string]string{
			"C-s": "save_current_buffer",
			"C-q": "quit",
		},
		strokes: map[string]string{
			"ff": "fuzzy_find",
		},
	}
}
