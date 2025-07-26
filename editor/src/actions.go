package main

type Action func(args ...any)

var ActionList = map[string](Action){
	// ----------------------
	// MVP
	// ----------------------
	"None": nil,

	"Prompt":     ACT_Prompt,
	"GoToLine":   ACT_GoToLine,
	"SelectNext": ACT_SelectNext,
	"ToggleLog":  ACT_ToggleLog,
	"LineBreak":  ACT_LineBreak,

	// move
	"CursorUp":                ACT_CursorUp,
	"CursorDown":              ACT_CursorDown,
	"CursorLeft":              ACT_CursorLeft,
	"CursorRight":             ACT_CursorRight,
	"CursorPageUp":            nil,
	"CursorPageDown":          nil,
	"CursorStart":             ACT_CursorStart,
	"CursorEnd":               ACT_CursorEnd,
	"CursorToViewTop":         nil,
	"CursorToViewCenter":      nil,
	"CursorToViewBottom":      nil,
	"CursorWordRight":         nil,
	"CursorWordLeft":          nil,
	"CursorSubWordRight":      nil,
	"CursorSubWordLeft":       nil,
	"CursorParagraphPrevious": nil,
	"CursorParagraphNext":     nil,
	"CursorCenterToView":      nil,

	// view
	"ViewStart":     nil,
	"ViewEnd":       nil,
	"ViewPageUp":    nil,
	"ViewPageDown":  nil,
	"ViewLinesUp":   nil,
	"ViewLinesDown": nil,

	// Multicursor
	"SpawnMultiCursor":       nil,
	"SpawnMultiCursorUp":     nil,
	"SpawnMultiCursorDown":   nil,
	"SpawnMultiCursorSelect": nil,
	"RemoveMultiCursor":      nil,
	"RemoveAllMultiCursors":  nil,
	"SkipMultiCursor":        nil,
	"SkipMultiCursorBack":    nil,

	// select
	"SelectAll":                 nil,
	"SelectUp":                  nil,
	"SelectDown":                nil,
	"SelectLeft":                nil,
	"SelectToStart":             nil,
	"SelectToEnd":               nil,
	"SelectPageUp":              nil,
	"SelectPageDown":            nil,
	"SelectRight":               nil,
	"SelectWordRight":           nil,
	"SelectWordLeft":            nil,
	"SelectSubWordRight":        nil,
	"SelectSubWordLeft":         nil,
	"SelectLine":                nil,
	"SelectToStartOfLine":       nil,
	"SelectToStartOfText":       nil,
	"SelectToStartOfTextToggle": nil,
	"SelectToEndOfLine":         nil,
	"SelectToParagraphPrevious": nil,
	"SelectToParagraphNext":     nil,

	// modify
	"InsertNewline":      nil,
	"InsertTab":          nil,
	"Backspace":          ACT_DelLeft,
	"Delete":             ACT_DelRight,
	"DeleteWordRight":    nil,
	"DeleteWordLeft":     nil,
	"DeleteSubWordRight": nil,
	"DeleteSubWordLeft":  nil,
	"IndentSelection":    nil,
	"OutdentSelection":   nil,
	"OutdentLine":        nil,
	"IndentLine":         nil,

	// I/O
	"Open":          nil,
	"Save":          nil,
	"SaveAll":       nil,
	"SaveAs":        nil,
	"Undo":          nil,
	"Redo":          nil,
	"Copy":          nil,
	"CopyLine":      nil,
	"Cut":           nil,
	"CutLine":       nil,
	"Paste":         nil,
	"PastePrimary":  nil,
	"Duplicate":     nil,
	"DuplicateLine": nil,
	"DeleteLine":    nil,
	"Quit":          ACT_Quit,
	"QuitAll":       nil,
	"ForceQuit":     nil,

	// ----------------------
	// Nice To Have
	// ----------------------

	"Find":                  nil,
	"FindLiteral":           nil,
	"FindNext":              nil,
	"FindPrevious":          nil,
	"DiffNext":              nil,
	"DiffPrevious":          nil,
	"Autocomplete":          nil,
	"CycleAutocompleteBack": nil,

	"HalfPageUp":            nil,
	"HalfPageDown":          nil,
	"StartOfText":           nil,
	"StartOfTextToggle":     nil,
	"StartOfLine":           nil,
	"EndOfLine":             nil,
	"ToggleHelp":            nil,
	"ToggleKeyMenu":         nil,
	"ToggleDiffGutter":      nil,
	"ToggleRuler":           nil,
	"ToggleHighlightSearch": nil,
	"UnhighlightSearch":     nil,
	"ResetSearch":           nil,
	"ClearStatus":           nil,
	"ShellMode":             nil,
	"CommandMode":           nil,
	"ToggleOverwriteMode":   nil,
	"Escape":                nil,

	// "AddTab":              nil,
	// "PreviousTab":         nil,
	// "NextTab":             nil,
	// "FirstTab":            nil,
	// "LastTab":             nil,
	// "NextSplit":           nil,
	// "PreviousSplit":       nil,
	// "FirstSplit":          nil,
	// "LastSplit":           nil,
	// "Unsplit":             nil,
	// "VSplit":              nil,
	// "HSplit":              nil,
	// "ToggleMacro":         nil,
	// "PlayMacro":           nil,
	// "Suspend (Unix only)": nil,
	// "ScrollUp":            nil,
	// "ScrollDown":          nil,
	// "JumpToMatchingBrace": nil,
	// "JumpLine":            nil,
	// "Deselect":            nil,
	// "ClearInfo":           nil,

}

func GetActionProgress() {
	total := 0
	implemented := 0

	for _, action := range ActionList {
		total++
		if action != nil {
			implemented++
		}
	}

	completion := float64(implemented) / float64(total) * 100

	LogF("Action Progress: %d/%d (%.1f%%)", implemented, total, completion)
}

// strokes
func ACT_ToggleLog(args ...any) { state.DisplayLogBuffer = !state.DisplayLogBuffer }

// shortcuts

func ACT_Quit(args ...any)        { state.Quit = true }
func ACT_Prompt(args ...any)      { OpenPrompt() }
func ACT_GoToLine(args ...any)    { Log("GotoLine") }
func ACT_CursorUp(args ...any)    { CursorMove(0, -1) }
func ACT_CursorDown(args ...any)  { CursorMove(0, 1) }
func ACT_CursorLeft(args ...any)  { CursorMove(-1, 0) }
func ACT_CursorRight(args ...any) { CursorMove(1, 0) }
func ACT_SelectNext(args ...any)  { Log("SelectNext") }
func ACT_DelLeft(args ...any)     { EditDelLeft() }
func ACT_DelRight(args ...any)    { EditDelRight() }
func ACT_LineBreak(args ...any)   { EditLineBreak() }
func ACT_CursorStart(args ...any) { CursorMoveStart() }
func ACT_CursorEnd(args ...any)   { CursorMoveEnd() }
