package main

import (
	"log"
	"os"

	"github.com/MassiGy/eddy/types"
	"github.com/gdamore/tcell/v2"
)

// will be populated by the build command
var version string
var binary_name string
var logFile *os.File

func main() {
	if os.Getenv("CURR_DEV_ENV") == "dev" {
		version = "v0.1.0"
		binary_name = "eddy"
	}
	defStyle := tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(boxStyle)
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Event loop
	ox, oy := 0, 0
	var r rune

	var jumpList types.JumpList
	jumpList = append(jumpList, types.JumpPoint{
		X:   ox,
		Y:   oy,
		Val: r,
	})

	limit := 100
	for {
		// Update screen
		s.Show()
		s.ShowCursor(ox, oy)
		s.SetCursorStyle(tcell.CursorStyleBlinkingBlock)

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {

		case *tcell.EventResize:
			s.Sync()

		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlQ {
				return

			} else if ev.Key() == tcell.KeyEnter {
				oy++
				ox = 0
				jumpList = append(jumpList, types.JumpPoint{
					X:   ox,
					Y:   oy,
					Val: r,
				})
				s.ShowCursor(ox, oy)

			} else if ev.Key() == tcell.KeyCtrlC {
				s.Clear()

			} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyCtrlZ {

				// get the prev point
				prevJumpPoint := jumpList[len(jumpList)-1]
				jumpList = jumpList[:len(jumpList)-1]
				ox = prevJumpPoint.X
				oy = prevJumpPoint.Y

				s.SetContent(ox, oy, r, nil, boxStyle)
				s.ShowCursor(ox, oy)

			} else {

				for _, k := range types.PrintableKeys {
					if ev.Rune() == rune(k) {
						s.SetContent(ox, oy, ev.Rune(), nil, defStyle)
						jumpList = append(jumpList, types.JumpPoint{
							X:   ox,
							Y:   oy,
							Val: ev.Rune(),
						})

						ox++
						if ox > limit {
							ox = 0
							oy++
						}

						s.ShowCursor(ox, oy)
					}
				}

			}
		}
	}
}
