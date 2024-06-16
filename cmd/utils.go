package main

import (
	"slices"

	"github.com/nsf/termbox-go"
)

func is_delimiter(ch rune) bool {
	return ch == ' ' ||
		ch == '\t' ||
		ch == '\'' ||
		ch == '\\' ||
		ch == '"' ||
		ch == '`' ||
		ch == '^' ||
		ch == '.' ||
		ch == ',' ||
		ch == ';' ||
		ch == ':' ||
		ch == '!' ||
		ch == '?' ||
		ch == '+' ||
		ch == '-' ||
		ch == '_' ||
		ch == '(' ||
		ch == ')' ||
		ch == '{' ||
		ch == '}' ||
		ch == '@' ||
		ch == '#' ||
		ch == '~' ||
		ch == '&' ||
		ch == '*' ||
		ch == '%' ||
		ch == '$' ||
		ch == '<' ||
		ch == '>' ||
		ch == '/'
}

func is_navigation_evt(ev termbox.Event) bool {
	res := false
	if ev.Ch == 0 {
		switch ev.Key {
		case
			termbox.KeyArrowDown,
			termbox.KeyArrowLeft,
			termbox.KeyArrowRight,
			termbox.KeyArrowUp:
			res = true

		case
			termbox.KeyHome,
			termbox.KeyEnd,
			termbox.KeyPgup,
			termbox.KeyPgdn:
			res = true
		}

	} else {
		switch ev.Ch {
		case 'k', 'l', 'h', 'j', 'f', 'b':
			res = current_mode == NORMAL
		}
	}
	return res
}

func is_io_evt(ev termbox.Event) bool {
	res := false
	if ev.Ch == 0 {
		switch ev.Key {
		case
			termbox.KeyCtrlS, termbox.KeyCtrlR:
			res = true
		}

	} else {
		switch ev.Ch {
		case 'r', 'w':
			res = current_mode == NORMAL
		}
	}

	return res
}

func is_configuration_evt(ev termbox.Event) bool {
	res := false
	if ev.Ch == 0 {
		switch ev.Key {
		case
			termbox.KeyCtrlN:
			res = true
		}

	} else {
		switch ev.Ch {
		// nothing for now
		}
	}

	return res
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
func register_curr_state() {
	blen := len(textBuffer)
	text_buffer_snapshot := [][]rune{}

	for i := 0; i < blen; i++ {
		llen := len(textBuffer[i])
		text_buffer_snapshot = append(text_buffer_snapshot, make([]rune, llen))
		copy(text_buffer_snapshot[i], textBuffer[i])
	}

	undo_stack = append(undo_stack, text_buffer_snapshot)

	// clear the redo stack
	redo_stack = nil
}

func undo() {
	// get topmost state from stack
	uslen := len(undo_stack)
	if uslen < 2 {
		return
	}
	undo_state := undo_stack[uslen-1]
	undo_stack = slices.Delete(undo_stack, uslen-1, uslen)

	// update buffer state
	if uslen-2 >= 0 {
		textBuffer = undo_stack[uslen-2]
	}

	// push poped state to redo stack
	redo_stack = append(redo_stack, undo_state)
}

func redo() {
	// get topmost state from stack
	rslen := len(redo_stack)
	if rslen == 0 {
		return
	}
	redo_state := redo_stack[rslen-1]
	redo_stack = slices.Delete(redo_stack, rslen-1, rslen)

	// update buffer state
	textBuffer = redo_state

	// push poped state to undo stack
	undo_stack = append(undo_stack, redo_state)

}
