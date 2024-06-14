package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/nsf/termbox-go"
)

// binary info
var version string
var binary_name string

type mode int

const (
	NORMAL mode = iota
	INSERT
	VISUAL
	PROMPT
)

// constants
var QUIT bool = false
var WRAP bool = true //TODO: add this to a config file
var WRAP_AFTER int = 80
var ROWS, COLS int
var SAVE_TO_FILE_MAX_ERROR_COUNT int = 3
var UNKOWN_SOURCE_FILENAME string = "nofile"
var LINE_NUMBER_COL_WIDTH int
var SHOW_LINE_NUMBERS bool = true
var LINES_COUNT int

// in buffer cursor info
var currentCol, currentRow int
var offsetX, offsetY int

// buffer
var textBuffer = [][]rune{}

// status bar flags
var modified bool = false
var err_message string
var info_message string
var keylog_message string = "	"
var current_mode mode

// read/write to file
var source_file string

func is_delimiter(ch rune) bool {
	return ch == ' ' ||
		ch == '\t' ||
		ch == ',' ||
		ch == ';' ||
		ch == ':' ||
		ch == '!' ||
		ch == '?' ||
		ch == '"' ||
		ch == '\'' ||
		ch == '.' ||
		ch == '/'
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
	l := len(textBuffer)
	if l == 0 {
		textBuffer = append(textBuffer, []rune{})
		return
	}

	beforeNewLineSegment := make([]rune, currentCol)
	afterNewLineSegment := make([]rune, len(textBuffer[currentRow])-currentCol)
	copy(beforeNewLineSegment, textBuffer[currentRow][:currentCol])
	copy(afterNewLineSegment, textBuffer[currentRow][currentCol:])

	textBuffer[currentRow] = beforeNewLineSegment

	if currentRow+1 < l {
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

		if currentCol >= l {
			// end of non empty line
			textBuffer[currentRow] = slices.Delete(textBuffer[currentRow], l-1, l)
		} else {
			// middle of non empty line
			textBuffer[currentRow] = slices.Delete(textBuffer[currentRow], currentCol-1, currentCol)
		}

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

	} else if currentCol > 0 && l <= 1 { /*TODO: is this really a valid case, compare to the one below*/
		// end of single character line or empty line

		textBuffer = slices.Delete(textBuffer, currentRow, currentRow+1)

		if currentRow > 0 {
			currentCol = len(textBuffer[currentRow-1])
			currentRow--
		} else {
			currentCol = 0
			currentRow = 0
		}
	} else if currentCol == 0 && l <= 1 {
		// start of single character line or empty line

		var ch rune
		if l != 0 {
			// capture the current line char
			ch = textBuffer[currentRow][currentCol]
		}

		textBuffer = slices.Delete(textBuffer, currentRow, currentRow+1)

		if currentRow > 0 {
			if l != 0 {
				// append captured char to prev line
				textBuffer[currentRow-1] = append(textBuffer[currentRow-1], ch)
			}
			currentCol = len(textBuffer[currentRow-1])
			currentRow--
		} else {
			currentCol = 0
			currentRow = 0
		}
	}
}
func delete_word(direction int) {
	if len(textBuffer) == 0 {
		return
	}

	l := len(textBuffer[currentRow])

	if l <= 1 {
		delete_character()
		return
	}

	if currentCol <= 1 && direction == -1 {
		delete_character()
		return

	} else if currentCol >= l-1 && direction == 1 {
		if currentRow+1 < len(textBuffer) {
			currentCol = 0
			currentRow++
			delete_character()
			return
		}
	}

	word_len := 1 // count the current char
	for i := currentCol + direction; i >= 0 && i < l; i += direction {
		if word_len > 0 && is_delimiter(textBuffer[currentRow][i]) {
			word_len++ // count the space or the \n
			break
		}
		word_len++
	}
	if word_len > 0 && direction == -1 {
		textBuffer[currentRow] = slices.Delete(textBuffer[currentRow], max(currentCol-word_len, 0), currentCol)
		currentCol = max(0, currentCol-word_len)
		l = max(0, l-word_len)
	} else if word_len > 0 && direction == 1 {
		textBuffer[currentRow] = slices.Delete(textBuffer[currentRow], currentCol, min(currentCol+word_len, len(textBuffer[currentRow])))
		l = max(0, l-word_len)
	}

	if l <= 0 {
		textBuffer = slices.Delete(textBuffer, currentRow, currentRow+1)
		if currentRow > 0 {
			currentCol = len(textBuffer[currentRow-1])
			currentRow--
		}
	}
}

func jump_word(direction int) {
	if len(textBuffer) == 0 {
		return
	}

	l := len(textBuffer[currentRow])

	if l <= 1 {
		if direction == 1 {
			currentCol = 0
			currentRow++
			return
		} else { // == -1
			if currentRow > 0 {
				currentCol = max(0, len(textBuffer[currentRow-1]))
				currentRow--
			}
			return
		}
	}

	if currentCol <= 1 && direction == -1 {
		if currentRow > 0 {
			currentCol = max(0, len(textBuffer[currentRow-1]))
			currentRow--
		}
		return
	} else if currentCol >= l-1 && direction == 1 {
		if currentRow < len(textBuffer)-1 {
			currentCol = 0
			currentRow++
		}
		return
	}

	word_len := 1
	for i := currentCol + direction; i >= 0 && i < l; i += direction {
		if word_len > 0 && is_delimiter(textBuffer[currentRow][i]) {
			break
		}
		word_len++
	}
	if word_len > 0 && direction == -1 {
		currentCol = max(0, currentCol-word_len)
	} else if word_len > 0 && direction == 1 {
		currentCol += word_len
	}
}

func handle_key_events(ev termbox.Event) {

	if ev.Ch == 0 {
		// non printable characters
		switch ev.Key {

		/* MANAGING MODES */
		case termbox.KeyCtrlQ:
			keylog_message = "Ctrl+q"
			QUIT = true
			return

		case termbox.KeyEsc:
			keylog_message = "Esc"
			current_mode = NORMAL

		/*  NAVIGATION  */
		case termbox.KeyArrowDown:
			keylog_message = "Down"
			if currentRow < len(textBuffer)-1 {
				currentRow++
			}

		case termbox.KeyArrowUp:
			keylog_message = "Up"
			if currentRow > 0 {
				currentRow--
			}

		case termbox.KeyArrowLeft:
			keylog_message = "Left"
			if currentCol > 0 {
				currentCol--
			} else if currentRow > 0 {
				currentCol = max(0, len(textBuffer[currentRow-1])-1)
				currentRow--
			}

		case termbox.KeyArrowRight:
			keylog_message = "Right"
			if len(textBuffer) == 0 {
				break
			} else if currentCol < len(textBuffer[currentRow])-1 {
				currentCol++
			} else if currentRow+1 < len(textBuffer) {
				currentCol = 0
				currentRow++
			}

		case termbox.KeyHome:
			keylog_message = "Start"
			currentCol = 0

		case termbox.KeyEnd:
			keylog_message = "End"
			currentCol = len(textBuffer[currentRow])

		case termbox.KeyPgup:
			keylog_message = "PgUp"
			currentRow = 0

		case termbox.KeyPgdn:
			keylog_message = "PgDown"
			currentRow = len(textBuffer) - 1

		/* DELIMETERS */
		case termbox.KeyTab:
			keylog_message = "Tab"
			if current_mode == INSERT {
				for i := 0; i < 4; i++ {
					insert_character(rune(' '))
				}
				modified = true
			}

		case termbox.KeySpace:
			keylog_message = "Space"
			if current_mode == INSERT {
				insert_character(rune(' '))
				modified = true
			}

		case termbox.KeyEnter:
			keylog_message = "Enter"
			if current_mode == INSERT {
				insert_newline()
				modified = true
			}

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			keylog_message = "BackSpace"
			if current_mode == INSERT {
				delete_character()
				modified = true
			}

		/* I/O on the file */
		case termbox.KeyCtrlS:
			keylog_message = "Ctrl+s"
			write_file(source_file)
			modified = false

		case termbox.KeyCtrlR:
			keylog_message = "Ctrl+r"
			textBuffer = nil
			read_file(source_file)
			modified = false
			current_mode = NORMAL

		/* EDITOR EXTRA CONTROL */
		case termbox.KeyCtrlN:
			keylog_message = "Ctrl+n"
			SHOW_LINE_NUMBERS = !SHOW_LINE_NUMBERS

		}

	} else {
		// printable characters
		if current_mode == INSERT {
			keylog_message = "	"
			insert_character(ev.Ch)
			modified = true

		} else if current_mode == NORMAL {
			keylog_message = string(ev.Ch)

			switch ev.Ch {

			case 'q':
				QUIT = true
				return

			case 'e', 'i':
				current_mode = INSERT

			case 'v':
				current_mode = VISUAL

			case 'p', '?', ':':
				current_mode = PROMPT

			case 'r': // reload
				textBuffer = nil
				read_file(source_file)
				current_mode = NORMAL
				modified = false

			case 'w':
				write_file(source_file)
				modified = false

			case 'k':
				if currentRow > 0 {
					currentRow--
				}

			case 'j':
				if currentRow < len(textBuffer)-1 {
					currentRow++
				}

			case 'h':
				if currentCol > 0 {
					currentCol--
				} else if currentRow > 0 {
					currentCol = max(0, len(textBuffer[currentRow-1])-1)
					currentRow--
				}

			case 'l':
				if len(textBuffer) == 0 {
					break
				} else if currentCol < len(textBuffer[currentRow])-1 {
					currentCol++
				} else if currentRow+1 < len(textBuffer) {
					currentCol = 0
					currentRow++
				}

			case 'f':
				jump_word(1)

			case 'b':
				jump_word(-1)

			case 'D':
				delete_word(1)
				modified = true

			case 'd':
				delete_word(-1)
				modified = true

			case 'C':
				delete_word(1)
				current_mode = INSERT
				modified = true

			case 'c':
				delete_word(-1)
				current_mode = INSERT
				modified = true
			}

		}
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
	if currentCol < offsetX {
		offsetX = currentCol
	}
	if currentCol >= COLS+offsetX {
		offsetX = currentCol - COLS + 1
	}
}

func wrap() {

	wrappedTextBuffer := [][]rune{}
	for _, row := range textBuffer {
		rowLen := len(row)

		if rowLen >= WRAP_AFTER {

			remainder := make([]rune, rowLen)
			copy(remainder, row)
			lremainder := len(remainder)
			for lremainder > WRAP_AFTER {

				// find the prev delimiter
				offset := 0
				for i := 0; i < WRAP_AFTER; i++ {
					if is_delimiter(remainder[WRAP_AFTER-i]) {
						offset = -(i - 1) // -1 to let the delimiter on prevline
						break
					}
				}
				chuck := make([]rune, WRAP_AFTER+offset)
				copy(chuck, remainder[:(WRAP_AFTER+offset)])
				wrappedTextBuffer = append(wrappedTextBuffer, chuck)

				remainder = slices.Delete(remainder, 0, WRAP_AFTER+offset)

				lremainder = len(remainder) // update the length
			}
			if len(remainder) > 0 {
				wrappedTextBuffer = append(wrappedTextBuffer, remainder)
			}

		} else {
			wrappedTextBuffer = append(wrappedTextBuffer, row)
		}
	}
	textBuffer = wrappedTextBuffer
}
func display_text_buffer() {
	var row, col int
	var txtBufRow, txtBufCol int
	linesCount := len(textBuffer)

	for row = 0; row < ROWS; row++ {
		txtBufRow = row + offsetY // scroll by offsetY lines

		if SHOW_LINE_NUMBERS {
			line_number_as_str := fmt.Sprintf("%d", txtBufRow+1)

			for i, ch := range line_number_as_str {
				termbox.SetCell(i, row, ch, termbox.ColorYellow, termbox.ColorDefault)
			}
		}

		for col = LINE_NUMBER_COL_WIDTH; col < COLS+LINE_NUMBER_COL_WIDTH; col++ {
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
	}
}

func display_status_bar() {

	modified_mark := ""
	if modified {
		modified_mark = "*"
	}

	line_numbers_mark := ""
	if SHOW_LINE_NUMBERS {
		line_numbers_mark = "[ln]"
	}

	wrap_mark := ""
	if WRAP {
		wrap_mark = fmt.Sprintf("[WRP=%d]", WRAP_AFTER)
	}

	current_file := source_file
	if len(source_file) > 10 {
		current_file = source_file[:10] + "~"
	}

	left_side_content := ""
	switch current_mode {
	case INSERT:
		left_side_content += "INSERT"
	case NORMAL:
		left_side_content += "NORMAL"
	case VISUAL:
		left_side_content += "VISUAL"
	case PROMPT:
		left_side_content += "PROMPT"
	}
	left_side_content += ""

	llen := 0
	llen = len(left_side_content)

	err_msg_len := len(err_message)
	info_msg_len := len(info_message)

	if err_msg_len == 0 && info_msg_len == 0 {
		left_side_content += fmt.Sprintf(" %s%s \t%d,%d \t%s %s", modified_mark, current_file, currentRow, currentCol, line_numbers_mark, wrap_mark)
		llen = len(left_side_content)
		for i := 0; i < llen; i++ {
			termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorBlack, termbox.ColorWhite)
		}
	} else if err_msg_len > 0 {
		left_side_content += fmt.Sprintf(" Error: %s\t", err_message)
		llen = len(left_side_content)
		for i := 0; i < llen; i++ {
			termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorWhite, termbox.ColorRed)
		}
		err_message = ""
	} else if info_msg_len > 0 {
		left_side_content += fmt.Sprintf(" Info: %s\t", info_message)
		llen = len(left_side_content)
		for i := 0; i < llen; i++ {
			termbox.SetCell(i, ROWS, rune(left_side_content[i]), termbox.ColorWhite, termbox.ColorBlue)
		}
		info_message = ""
	}

	// right_side_content := binary_name + " " + version
	right_side_content := fmt.Sprintf("%s\t", keylog_message)
	rlen := len(right_side_content)

	padding := COLS - llen - rlen
	if SHOW_LINE_NUMBERS {
		padding += LINE_NUMBER_COL_WIDTH
	}

	for i := 0; i < padding; i++ {
		termbox.SetCell(i+llen, ROWS, rune(' '), termbox.ColorWhite, termbox.ColorWhite)
	}

	for i := 0; i < rlen; i++ {
		termbox.SetCell(i+llen+padding, ROWS, rune(right_side_content[i]), termbox.ColorBlack, termbox.ColorWhite)
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
		line = scanner.Text() // this does not read the \n
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
		}
		// since we do not read the \n in read_file
		// and we do not insert \n to the textbuff
		w.WriteRune('\n')
	}
	w.Flush()
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
		about_file_path := fmt.Sprintf("%s/.config/%s-%s/about", os.Getenv("HOME"), binary_name, version)

		// check for about file existance
		if _, err := os.Stat(about_file_path); err == nil {
			source_file = about_file_path
			read_file(source_file)
		} else {
			source_file = UNKOWN_SOURCE_FILENAME
			info_message = "Input file messing."
			textBuffer = append(textBuffer, []rune{})
		}
	}
	currentCol = LINE_NUMBER_COL_WIDTH
	currentRow = 0
	current_mode = NORMAL

	for !QUIT {
		if SHOW_LINE_NUMBERS {
			LINE_NUMBER_COL_WIDTH = 2 + len(fmt.Sprintf("%d", len(textBuffer)))
		} else {
			LINE_NUMBER_COL_WIDTH = 0
		}
		if WRAP && current_mode == NORMAL {
			wrap() // wrap & update text buffer
			offsetX = 0
			currentCol = min(currentCol, len(textBuffer[currentRow]))
		}

		COLS, ROWS = termbox.Size()   // re-evaluate each time to synch with size change
		ROWS--                        // for the status bar
		COLS -= LINE_NUMBER_COL_WIDTH // for the linenumber column

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
