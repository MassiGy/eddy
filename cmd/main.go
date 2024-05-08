package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/MassiGy/eddy/types"
	"github.com/gdamore/tcell/v2"
)

// will be populated by the build command
var version string
var binary_name string

var s tcell.Screen
var defStyle, boxStyle tcell.Style

var err error
var logFile *os.File

// for the line numbers column
var offset int

var limit int
var ox, oy int
var linesCharCount map[int]int

var emptyChar rune

var jumpList types.JumpList

func main() {
	if os.Getenv("CURR_DEV_ENV") == "dev" {
		version = "v0.1.0"
		binary_name = "eddy"
	}

	setup()
	defer quit()
	defer logFile.Close()

	// Event loop
	for {
		// Update screen
		s.Show()

		// Show cursor
		s.ShowCursor(ox, oy)
		s.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
		// showNumberLine()

		printed := false
		if !printed {

			logMsg("\n")
			for i, p := range jumpList {

				if p.Val == '\n' {
					logMsg("found a newline in jumpList (printer)\n")
					logMsg("\n")
				} else {
					logMsg(string(p.Val))
				}

				if i > 40 {
					logMsg("\n")
				}
			}
		}

		updatelinesCharCount()
		ev := s.PollEvent()

		printed = true

		switch ev := ev.(type) {

		case *tcell.EventResize:
			s.Sync()

		case *tcell.EventKey:

			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlQ {
				return

			} else if ev.Key() == tcell.KeyEnter {
				registerNewKey(ev)

			} else if ev.Key() == tcell.KeyUp {
				handleKeyUp()

			} else if ev.Key() == tcell.KeyDown {
				handleKeyDown()

			} else if ev.Key() == tcell.KeyLeft {
				handleKeyLeft()

			} else if ev.Key() == tcell.KeyRight {
				handleKeyRight()

			} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
				handleBackSpace()
			} else {
				registerNewKey(ev)
				// handlePrintableChars(ev)
			}

			s.Clear()
			showNumberLine()
			setContent()
			s.ShowCursor(ox+1, oy)
		}
	}
}

func setup() {
	// Setup our variables and initial state
	offset = 4
	limit = 100
	ox, oy = offset+1, 0

	linesCharCount = make(map[int]int)

	// for debugging
	if os.Getenv("DEBUG") == "true" {
		logFile, err = os.OpenFile("./log.out", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic("failed to open log file")
		}
	}

	// Initialize screen
	s, err = tcell.NewScreen()
	if err != nil {
		panic(err.Error())
	}

	boxStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkGray)
	defStyle = tcell.StyleDefault.Foreground(tcell.ColorNavajoWhite).Background(tcell.ColorDarkGray)

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(boxStyle)
	s.Clear()

}
func quit() {

	maybePanic := recover()

	s.Fini() // free the resources then panic back
	if maybePanic != nil {
		panic(maybePanic)
	}
}
func logMsg(msg string) {
	if os.Getenv("DEBUG") == "true" {
		logFile.WriteString(msg)
	}
}

func showNumberLine() {

	for k, _ := range linesCharCount {
		runes := []rune(fmt.Sprintf("%d", k+1))
		i := 0
		for i = 0; i < len(runes); i++ {
			s.SetContent(i, k, runes[i], nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkBlue))
		}
		for j := i; j < offset; j++ {
			s.SetContent(j, k, rune(' '), nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkBlue))
		}
	}
}

func updatelinesCharCount() {
	for k := range linesCharCount {
		linesCharCount[k] = 0
	}

	for _, p := range jumpList {
		linesCharCount[p.Y] += 1
	}
}

func registerNewKey(ev *tcell.EventKey) {

	l := len(jumpList)
	targetIndex := ox - offset - 1

	counter := 0
	for _, v := range linesCharCount {
		if counter < oy {
			targetIndex += v
		}
	}

	logMsg(fmt.Sprintf("target index %#v\n", targetIndex))

	if targetIndex >= l {
		ox++
		jumpList = append(jumpList, types.JumpPoint{
			X:   ox,
			Y:   oy,
			Val: ev.Rune(),
		})

		if ev.Key() == tcell.KeyEnter {
			ox = offset + 1
			oy++
		}
	}

	if targetIndex < l {
		ox++
		jumpList = slices.Insert(jumpList, targetIndex+1, types.JumpPoint{
			X:   ox,
			Y:   oy,
			Val: ev.Rune(),
		})

		if ev.Key() == tcell.KeyEnter {
			jumpList[targetIndex].NewLine = true
			ox = offset + 1
			oy++
		}

	}
}

func handleBackSpace() {

	l := len(jumpList)

	if l == 0 {
		return
	}
	if l == 1 {
		jumpList = types.JumpList{}
		ox = offset + 1
		oy = 0
	}
	if l >= 2 {
		// pop the last elem
		jumpList = slices.Delete(jumpList, l-1, l)

		ox = jumpList[l-2].X
		oy = jumpList[l-2].Y
	}
}

func handleKeyUp() {
	if oy > 0 {
		oy--
	}
}
func handleKeyDown() {
	if linesCharCount[oy+1] > 0 {
		oy++

		if ox > linesCharCount[oy] {
			ox = linesCharCount[oy] + offset + 1
		}
	}
}

func handleKeyLeft() {
	if ox > offset+1 {
		ox--
		s.ShowCursor(ox, oy)

	} else if oy > 0 && linesCharCount[oy-1] > 1 {

		oy--
		ox = linesCharCount[oy] + offset + 1

	}
}

func handleKeyRight() {
	if ox-offset-1 < linesCharCount[oy] {
		ox++
	} else {
		handleKeyDown()
	}
}

func setContent() {
	newLinesCount := 0

	for _, p := range jumpList {
		if p.Val == '\t' {
			s.SetContent(p.X+1, p.Y+newLinesCount, ' ', nil, defStyle)
			s.SetContent(p.X+2, p.Y+newLinesCount, ' ', nil, defStyle)
			s.SetContent(p.X+3, p.Y+newLinesCount, ' ', nil, defStyle)
			s.SetContent(p.X+4, p.Y+newLinesCount, ' ', nil, defStyle)
			ox += 4
		} else if p.NewLine {
			newLinesCount++
			// oy++
		} else {
			s.SetContent(p.X, p.Y+newLinesCount, p.Val, nil, defStyle)
		}
	}
}
