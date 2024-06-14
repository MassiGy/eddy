package main

import "github.com/nsf/termbox-go"

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
