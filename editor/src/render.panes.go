package main

// just so you know the main idea is to only have\support two panes
// there is really not much i can think of where you'd need more
// and if you can use tmux or w/e to do that

// lets draw a separator between the two panes
// hmm, i think i could draw two status lines, makes more sense
// one per split
// have to think about this a bit, i guess the areas should have a statusbar
// visible or not

type PaneLayout int

const (
	PaneLayoutHorizontal PaneLayout = iota
	PaneLayoutVertical
)

type Pane struct {
	ID        int
	Visible   bool
	Area      *Area
	StatusBar *StatusBar
	Split     float32
}

type Panes struct {
	One    *Pane
	Two    *Pane
	Layout PaneLayout
}

func (p *Pane) NormalizeSplit() {
	p.Split = Clamp(p.Split, 0, 1)
	if !p.Visible {
		p.Area.Hidden = true
		p.Split = 0
	} else {
		p.Area.Hidden = false
		if p.Split <= 0 {
			p.Split = 1
		}
	}
}

func (p *Panes) NormalizeSplit() {
	if p.One == nil || p.Two == nil {
		LogErr("Panes not initialized properly")
		return
	}

	p.One.NormalizeSplit()
	p.Two.NormalizeSplit()

	max := -1
	if p.Layout == PaneLayoutHorizontal {
		max = state.Screen.Size.X
	}

	if p.Layout == PaneLayoutVertical {
		max = state.Screen.Size.Y
	}

	current := float32(max) * (p.One.Split + p.Two.Split)
	if current > float32(max) {
		p.One.Split = (p.One.Split * float32(max)) / current
		p.Two.Split = (p.Two.Split * float32(max)) / current
	}

	if p.Layout == PaneLayoutHorizontal {
		height := state.Screen.Size.Y
		p.One.Area.SetSize(Vec2{X: int(float32(max) * p.One.Split), Y: height})
		p.Two.Area.SetSize(Vec2{X: int(float32(max) * p.Two.Split), Y: height})
		p.Two.Area.SetPosition(Vec2{X: int(float32(max) * p.One.Split), Y: 0})
	}

	if p.Layout == PaneLayoutVertical {
		width := state.Screen.Size.X
		p.One.Area.SetSize(Vec2{X: width, Y: int(float32(max) * p.One.Split)})
		p.Two.Area.SetSize(Vec2{X: width, Y: int(float32(max) * p.Two.Split)})
		p.Two.Area.SetPosition(Vec2{X: 0, Y: int(float32(max) * p.One.Split)})
	}

}

func InitPanes() {
	state.Panes = &Panes{
		One: &Pane{
			ID:      0,
			Visible: true,
			Area:    CreateArea("Pane1").SetZIndex(BASE_ZINDEX),
			// StatusBar: nil,
			Split: 1,
		},
		Two: &Pane{
			ID:      1,
			Visible: false,
			Area:    CreateArea("Pane2").SetZIndex(BASE_ZINDEX),
			// StatusBar: nil,
			Split: 1,
		},
	}

	// CreateStatusBar(state.Panes.One)
	// CreateStatusBar(state.Panes.Two)
}

func UpdatePanes() {
	state.Panes.Update()
}

func (p *Panes) Update() {
	p.NormalizeSplit()
	// p.One.StatusBar.Update()
	// p.Two.StatusBar.Update()
}
