package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ROWS, COLS int
var QUIT bool = false
var LINE_NUMBER_COL_WIDTH int = 5
var currentCol, currentRow int
var offsetX, offsetY int

var textBuffer = [][]rune{}
var source_file string

func print_message(col, row int, fg, bg termbox.Attribute, msg string) {
	for _, ch := range msg {
		termbox.SetCell(col, row, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
}
func handle_key_events(ev termbox.Event) {

	switch ev.Key {
	case termbox.KeyCtrlQ, termbox.KeyEsc:
		QUIT = true
		return
	case termbox.KeyArrowDown:
		if currentRow < len(textBuffer)-1 {
			currentRow++
		}

	case termbox.KeyArrowUp:
		if currentRow != 0 {
			currentRow--
		}

	case termbox.KeyArrowLeft:
		if currentCol != 0 {
			currentCol--
		} else if currentRow > 0 {
			currentCol = len(textBuffer[currentRow-1])
			currentRow--
		}

	case termbox.KeyArrowRight:
		if currentCol < len(textBuffer[currentRow]) {
			currentCol++
		} else if currentRow+1 < len(textBuffer) {
			currentCol = 0
			currentRow++
		}

	}
	scroll_text_buffer()
	if currentCol > len(textBuffer[currentRow]) {
		currentCol = len(textBuffer[currentRow])
	}

}

func scroll_text_buffer() {
	if currentRow >= ROWS+offsetY {
		offsetY = currentRow - ROWS + 1
	}
	if currentRow < offsetY {
		offsetY = currentRow
	}
	if currentCol < offsetX {
		offsetX = currentCol
	}
	if currentCol >= COLS+offsetX {
		offsetX = currentCol - COLS + 1
	}
}

func display_text_buffer() {
	var row, col int
	var txtBufRow, txtBufCol int
	linesCount := len(textBuffer)

	for row = 0; row < ROWS; row++ {
		txtBufRow = row + offsetY // scroll by offsetY lines

		for col = 0; col < COLS; col++ {
			txtBufCol = col + offsetX // scroll by offsetX columns

			// display the text buffer content
			if txtBufRow < linesCount && txtBufCol < len(textBuffer[txtBufRow]) {
				if textBuffer[txtBufRow][txtBufCol] != '\t' {
					termbox.SetChar(col, row, textBuffer[txtBufRow][txtBufCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorDefault)
				}
			} else if txtBufRow >= linesCount { // print * if line is empty
				termbox.SetCell(0, row, rune('*'), termbox.ColorBlue, termbox.ColorDefault)
			}
		}

		// new line at the end of each row
		termbox.SetChar(col, row, rune('\n'))
	}
}
func read_file(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error on openning file %s : %v\n", filename, err.Error())
		return
	}
	defer file.Close()

	source_file = filename
	lineNumber := 0

	scanner := bufio.NewScanner(file)
	var line string
	var l int

	for scanner.Scan() {
		textBuffer = append(textBuffer, []rune{})
		line = scanner.Text()
		l = len(line)

		for i := 0; i < l; i++ {
			textBuffer[lineNumber] = append(textBuffer[lineNumber], rune(line[i]))
		}
		lineNumber++
	}
	if lineNumber == 0 {
		textBuffer = append(textBuffer, []rune{})
	}
}
func display_status_bar() {

	left_side_content := fmt.Sprintf(" File: %s\t line:%d,col:%d", source_file, currentRow, currentCol)
	llen := len(left_side_content)

	for i := 0; i < llen; i++ {
		termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorBlack, termbox.ColorWhite)
	}
	right_side_content := "Eddy v0.1.0 "
	rlen := len(right_side_content)

	padding := COLS - llen - rlen

	for i := 0; i < padding; i++ {
		termbox.SetCell(i+llen, ROWS, rune(' '), termbox.ColorWhite, termbox.ColorWhite)
	}

	for i := 0; i < rlen; i++ {
		termbox.SetCell(i+llen+padding, ROWS, rune(right_side_content[i]), termbox.ColorBlack, termbox.ColorWhite)
	}
}

func run_editor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()

	if len(os.Args) > 1 {
		source_file = os.Args[1]
		read_file(source_file)
	} else {
		source_file = "out.txt"
		textBuffer = append(textBuffer, []rune{})
	}
	for !QUIT {

		COLS, ROWS = termbox.Size() // re-evaluate each time to synch with size change
		ROWS--                      // for the status bar

		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		display_status_bar()
		display_text_buffer()
		termbox.SetCursor(currentCol-offsetX, currentRow-offsetY)
		termbox.Flush()

		evt := termbox.PollEvent()
		switch evt.Type {
		case termbox.EventKey:
			handle_key_events(evt)
		case termbox.EventError:
			return // call defered routines
		}
	}
}
func main() {
	run_editor()
}
