package main

//   → strokes
// > → search actions by fuzzy match
// : → execute Lua commands
// @ → go to symbol
// % → search text
// = → evaluate/inspect

type PromptMode int

const (
	StrokePrompt PromptMode = iota
	FuzzyPrompt
)

type Prompt struct {
	Visible bool
	Query   string
	Mode    PromptMode
	Area    *Area
	// Filtered []CommandEntry
}

func InitPrompt() {
	a := CreateArea("Prompt").SetZIndex(FRONT_ZINDEX)
	a.SetBorderStyle(BorderRounded)
	a.SetBorderUsage(BorderUsage{
		Top:    true,
		Bottom: true,
		Left:   true,
		Right:  true,
	})

	p := &Prompt{
		Visible: false,
		Query:   "",
		Mode:    StrokePrompt,
		Area:    a,
	}

	ctx := &Context{
		Buffer: &Buffer{},
	}
	p.Area.Context = ctx
	state.Prompt = p
}

func OpenPrompt() {
	state.Prompt.Visible = true
}

func ClosePrompt() {
	state.Prompt.Visible = false
}

func ClearPrompt() {
	state.Prompt.Query = ""
}

func PromptInput(r rune) {
	state.Prompt.Query += string(r)
}

func SubmitPrompt() {
	query := state.Prompt.Query
	if !TryActionExecute(state.Binds.Strokes, state.Prompt.Query) {
		LogErr("Nothing found for", query)
	}
}

// IMPORTANT:
// need to focus this as current buffer
// i cant type today,
// anyhow, i'm not sure why cursor/manipulation is not working
func UpdatePrompt() {
	p := state.Prompt

	p.Area.Hidden = !p.Visible
	if !p.Visible {
		return
	}

	// if query != lastquery -> update request

	p.Area.Context.Buffer.SetLine(0, []rune("> "+p.Query))

	p.Area.SetSize(Vec2{
		X: state.Screen.Size.X,
		Y: 1,
	})
	p.Area.SetPosition(Vec2{
		X: 0,
		Y: 1, // TOOD: get from status bar
	})

	// p.Area.CalculateDrawBounds()
	// for x := 0; x < p.Area.Bounds.start.X; x++ {
	// 	for y := 0; y < p.Area.Bounds.start.Y; y++ {
	// 		state.Screen.Cell(Vec2{X: x, Y: y}, ' ', state.Theme.Cursor.Main)
	// 	}
	// }

	// // w, h := w, kit.Screen.Size.Y
	// w := state.Screen.Size.X
	// h := 1
	// // DrawPanel(w/2-cmd_width/2, h/4-4, cmd_width, 3)

	// var startY int
	// var startX int = 0

	// if state.Config.PromptPosition == PromptTop {
	// 	startY = 1
	// }

	// if state.Config.PromptPosition == PromptCenter {
	// 	// use border and a frame for it
	// 	w = int(w / 3)
	// 	startY = state.Screen.Size.Y / 2
	// 	startX = w/2 - w/2
	// }

	// if state.Config.PromptPosition == PromptBottom {
	// 	startY = state.Screen.Size.Y - 2 // +1 because of statusbar
	// }

	// state.Screen.ClearRegion(
	// 	startX,
	// 	startY,
	// 	w,
	// 	h,
	// )

	// state.Screen.Cell(startX, startY, '>', state.Theme.Text)

	// for i, r := range state.Prompt.Query {
	// 	state.Screen.Cell(startX+i+1, startY, r, state.Theme.Text)
	// }

	// // // draw the first line with the query

}
