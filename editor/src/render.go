package main

import (
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
)

func StartRender() {
	Update()
	RenderAreas()

	if state.Config.FPS <= 0 {
		state.Config.FPS = 16
	} else if state.Config.FPS > 60 {
		state.Config.FPS = 60
	}

	go func() {

		tick := time.Second / time.Duration(state.Config.FPS)
		ticker := time.NewTicker(tick)
		defer func() {
			ticker.Stop()
			Shutdown()
		}()

		lastFrame := time.Now()

		for !state.Quit {
			<-ticker.C
			now := time.Now()
			state.Delta = now.Sub(lastFrame).Milliseconds()
			lastFrame = now
			Update()
			RenderAreas()
		}
	}()

	for !state.Quit {
		ev := state.Screen.Poll()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			HandleInput(ev)
		case *tcell.EventMouse:
			HandleMouseInput(ev)
		case *tcell.EventResize:
			state.Screen.UpdateSize()
		}

	}
}

func RenderAreas() {
	var zIndexes []int
	for z := range state.ZIndexedAreas {
		zIndexes = append(zIndexes, z)
	}
	sort.Ints(zIndexes)

	for _, z := range zIndexes {
		areas := state.ZIndexedAreas[z]
		for _, area := range areas {
			area.Draw()
		}
	}

	state.Screen.FlushAll()
}
