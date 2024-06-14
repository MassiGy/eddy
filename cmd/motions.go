package main

import "slices"

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
