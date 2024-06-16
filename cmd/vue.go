package main

import (
	"fmt"
	"slices"

	"github.com/nsf/termbox-go"
)

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
	if len(source_file) > 15 {
		current_file = source_file[:15] + "~"
	}

	left_side_content := ""
	switch current_mode {
	case INSERT:
		left_side_content += "INSERT"
	case NORMAL:
		left_side_content += "NORMAL"
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
