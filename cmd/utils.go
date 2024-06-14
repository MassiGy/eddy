package main

import "slices"

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
