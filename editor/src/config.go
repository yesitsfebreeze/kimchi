package main

type Config struct {
	FPS              int            `help:"cfg('fps', number)"`
	ProjectScanLimit int            `help:"cfg('project_scan_limit', number)"`
	XTerm            bool           `help:"cfg('xterm', true|false)"`
	AutosaveOnClose  bool           `help:"cfg('autosave_on_close', true|false)"`
	MouseEnabled     bool           `help:"cfg('mouse_enabled', true|false)"`
	CopyPaste        bool           `help:"cfg('copy_paste', true|false)"`
	SurroundingLines int            `help:"cfg('surrounding_lines', number)"`
	PromptPosition   PromptPosition `help:"cfg('prompt position', 'top'|'center'|'bottom')"`
	UnfocusedDarken  float32        `help:"cfg('unfocused_darken', number)"`
	DefaultLanguages []string       `help:"cfg('default_languages', [string])"`
	CursorTrail      CursorTrailConfig
	Statusbar        StatusbarConfig
	Indent           IndentConfig
	WhiteSpace       WhiteSpaceConfig
	Plugins          map[string]string
}

type CursorTrailConfig struct {
	Enabled bool `help:"cfg('cursor_trail.enabled', true|false)"`
	Time    int  `help:"cfg('cursor_trail.time', number(ms))"`
	Length  int  `help:"cfg('cursor_trail.length', number)"`
}

type StatusbarLayoutConfig struct {
}

type StatusbarConfig struct {
	Enabled      bool              `help:"cfg('statusbar.enabled', true|false)"`
	UseSeparator bool              `help:"cfg('statusbar.use_separator', true|false)"`
	Separator    string            `help:"cfg('statusbar.separator', string)"`
	Position     StatusbarPosition `help:"cfg('statusbar.position', 'top'|'bottom')"`
	Left         string            `help:"cfg('statusbar.left', string)"`
	Center       string            `help:"cfg('statusbar.center', string)"`
	Right        string            `help:"cfg('statusbar.right', string)"`
}

type IndentOptions struct {
	Style IndentStyle
	Width int
}

type IndentConfig struct {
	Visual IndentOptions `help:"cfg('indent.visual', { style: 'spaces'|'tabs', width: number })"`
	Save   IndentOptions `help:"cfg('indent.save', { style: 'spaces'|'tabs', width: number })"`
}

type WhiteSpaceConfig struct {
	Space rune
	Tab   rune
	Eol   rune
}

// #region Enums

type StatusbarPosition int

const (
	StatusbarTop StatusbarPosition = iota
	StatusbarBottom
)

var StatusbarPositionEnum = NewEnumMap(map[string]StatusbarPosition{
	"top":    StatusbarTop,
	"bottom": StatusbarBottom,
})

type IndentStyle int

const (
	IndentSpaces IndentStyle = iota
	IndentTabs
)

var IndentStyleEnum = NewEnumMap(map[string]IndentStyle{
	"spaces": IndentSpaces,
	"tabs":   IndentTabs,
})

type PromptPosition int

const (
	PromptTop PromptPosition = iota
	PromptCenter
	PromptBottom
)

var PromptPositionEnum = NewEnumMap(map[string]PromptPosition{
	"top":    PromptTop,
	"center": PromptCenter,
	"bottom": PromptBottom,
})

// #endregion Enums
