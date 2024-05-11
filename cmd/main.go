package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// binary info
var version string
var binary_name string

// constants
var QUIT bool = false
var ROWS, COLS int
var SAVE_TO_FILE_MAX_ERROR_COUNT int = 3
var UNKOWN_SOURCE_FILENAME string = "nofile"
var LINE_NUMBER_COL_WIDTH int
var SHOW_LINE_NUMBERS bool = true

// in buffer cursor info
var currentCol, currentRow int
var offsetX, offsetY int

// buffer
var textBuffer = [][]rune{}

// status bar flags
var modified bool = false
var err_message string
var info_message string

// read/write to file
var source_file string

func print_message(col, row int, fg, bg termbox.Attribute, msg string) {
	for _, ch := range msg {
		termbox.SetCell(col, row, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
}

func insert_character(ch rune) {
	if len(textBuffer) == 0 {
		textBuffer = append(textBuffer, []rune{ch})
		return
	}
	textBuffer[currentRow] = slices.Insert(textBuffer[currentRow], currentCol, ch)
	currentCol++
}
func insert_newline() {
	if len(textBuffer) == 0 {
		textBuffer = append(textBuffer, []rune{})
		return
	}

	beforeNewLineSegment := make([]rune, currentCol)
	afterNewLineSegment := make([]rune, len(textBuffer[currentRow])-currentCol)
	copy(beforeNewLineSegment, textBuffer[currentRow][:currentCol])
	copy(afterNewLineSegment, textBuffer[currentRow][currentCol:])

	beforeNewLineSegment = append(beforeNewLineSegment, '\n')
	textBuffer[currentRow] = beforeNewLineSegment

	if currentRow+1 < len(textBuffer) {
		textBuffer = slices.Insert(textBuffer, currentRow+1, afterNewLineSegment)
	} else {
		textBuffer = append(textBuffer, afterNewLineSegment)
	}

	currentCol = 0 // these will be offseted by scroll_buffer()
	currentRow++
}
func delete_character() {

	if len(textBuffer) == 0 {
		return
	}

	l := len(textBuffer[currentRow])

	if l > 1 && currentCol > 0 {
		// middle of non empty line

		textBuffer[currentRow] = slices.Delete(textBuffer[currentRow], currentCol-1, currentCol)
		currentCol--

	} else if currentCol == 0 && l > 1 && currentRow > 0 {
		// start of non empty line
		afterCursorLineSegment := make([]rune, l)
		copy(afterCursorLineSegment, textBuffer[currentRow])

		textBuffer = slices.Delete(textBuffer, currentRow, currentRow+1)
		currentRow--

		l = len(textBuffer[currentRow])

		textBuffer[currentRow] = append(textBuffer[currentRow], afterCursorLineSegment...)

		currentCol = l

	} else if l <= 1 {
		// end of single character line

		textBuffer = slices.Delete(textBuffer, currentRow, currentRow+1)

		if currentRow > 0 {
			currentCol = len(textBuffer[currentRow-1])
			currentRow--
		} else {
			currentCol = 0
			currentRow = 0
		}
	}
}

func handle_key_events(ev termbox.Event) {

	if ev.Ch == 0 {
		// non printable characters

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

		case termbox.KeyHome:
			currentCol = 0
		case termbox.KeyEnd:
			currentCol = len(textBuffer[currentRow])

		case termbox.KeyTab:
			for i := 0; i < 4; i++ {
				insert_character(rune(' '))
			}
			modified = true

		case termbox.KeySpace:
			insert_character(rune(' '))
			modified = true

		case termbox.KeyEnter:
			insert_newline()
			modified = true

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			delete_character()
			modified = true

		case termbox.KeyCtrlS:
			write_file(source_file)
			modified = false

		case termbox.KeyCtrlN:
			SHOW_LINE_NUMBERS = !SHOW_LINE_NUMBERS

		}

	} else {
		// printable characters
		insert_character(ev.Ch)
		modified = true
	}
	if len(textBuffer) > 0 && currentCol > len(textBuffer[currentRow]) {
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
	if currentCol < offsetX+LINE_NUMBER_COL_WIDTH {
		offsetX = currentCol - LINE_NUMBER_COL_WIDTH
	}
	if currentCol >= COLS+offsetX-LINE_NUMBER_COL_WIDTH {
		offsetX = currentCol - COLS + 1 + LINE_NUMBER_COL_WIDTH
	}
}

func display_text_buffer() {
	var row, col int
	var txtBufRow, txtBufCol int
	linesCount := len(textBuffer)

	if SHOW_LINE_NUMBERS {
		LINE_NUMBER_COL_WIDTH = 2 + len(fmt.Sprintf("%d", linesCount))
	} else {
		LINE_NUMBER_COL_WIDTH = 0
	}

	for row = 0; row < ROWS; row++ {
		txtBufRow = row + offsetY // scroll by offsetY lines

		if SHOW_LINE_NUMBERS {
			line_number_as_str := fmt.Sprintf("%d", row+1)

			for i, ch := range line_number_as_str {
				termbox.SetCell(i, row, ch, termbox.ColorYellow, termbox.ColorDefault)
			}
		}

		for col = LINE_NUMBER_COL_WIDTH; col < COLS; col++ {
			txtBufCol = col - LINE_NUMBER_COL_WIDTH + offsetX // scroll by offsetX columns

			// display the text buffer content
			if txtBufRow < linesCount && txtBufCol < len(textBuffer[txtBufRow]) {
				if textBuffer[txtBufRow][txtBufCol] != '\t' {
					termbox.SetChar(col, row, textBuffer[txtBufRow][txtBufCol])
				} else {
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorDefault)
				}
			} else if txtBufRow >= linesCount { // for unreached lines print ~ as in vim
				termbox.SetCell(LINE_NUMBER_COL_WIDTH, row, rune('~'), termbox.ColorBlue, termbox.ColorDefault)
			}
		}

		// new line at the end of each row
		termbox.SetChar(col, row, rune('\n'))
	}
}
func read_file(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		err_message = "Could not read file."
		textBuffer = append(textBuffer, []rune{})
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

func write_file(filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		err_message = "Could not open file."

		if SAVE_TO_FILE_MAX_ERROR_COUNT > 0 {
			SAVE_TO_FILE_MAX_ERROR_COUNT--
			write_file("out.txt") // fallback
		}
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	rows := len(textBuffer)
	for row := 0; row < rows; row++ {

		cols := len(textBuffer[row])

		for col := 0; col < cols; col++ {
			w.WriteRune(textBuffer[row][col])
			if col == cols-1 && textBuffer[row][col] != '\n' {
				w.WriteRune('\n')
			}
		}
	}
	w.Flush()
}

func display_status_bar() {

	modified_mark := ""
	if modified {
		modified_mark = "*"
	}

	current_file := source_file
	if len(source_file) > 8 {
		current_file = source_file[:8] + "..."

	}

	left_side_content := ""
	llen := 0
	llen = len(left_side_content)

	err_msg_len := len(err_message)
	info_msg_len := len(info_message)

	if err_msg_len == 0 && info_msg_len == 0 {
		left_side_content = fmt.Sprintf(" File: %s%s\t line:%d,col:%d", modified_mark, current_file, currentRow, currentCol)
		llen = len(left_side_content)
		for i := 0; i < llen; i++ {
			termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorBlack, termbox.ColorWhite)
		}
	} else if err_msg_len > 0 {
		left_side_content = fmt.Sprintf(" Error: %s\t", err_message)
		llen = len(left_side_content)
		for i := 0; i < llen; i++ {
			termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorWhite, termbox.ColorRed)
		}
		err_message = ""
	} else if info_msg_len > 0 {
		left_side_content = fmt.Sprintf(" Info: %s\t", info_message)
		llen = len(left_side_content)
		for i := 0; i < llen; i++ {
			termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorWhite, termbox.ColorBlue)
		}
		info_message = ""
	}

	right_side_content := binary_name + " " + version
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
		source_file = UNKOWN_SOURCE_FILENAME
		info_message = "Input file messing."
		textBuffer = append(textBuffer, []rune{})
	}

	for !QUIT {

		COLS, ROWS = termbox.Size() // re-evaluate each time to synch with size change
		ROWS--                      // for the status bar

		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		scroll_text_buffer()

		display_status_bar()

		display_text_buffer()

		termbox.SetCursor(currentCol+LINE_NUMBER_COL_WIDTH-offsetX, currentRow-offsetY)

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
	if os.Getenv("ENV") == "dev" {
		binary_name = "eddy"
		version = "v0.1.0"
	}
	run_editor()
}
