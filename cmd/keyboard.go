package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

var (
	// this is specific to this file,
	// it basically tracks how many
	// redos have been perfomed since
	// the last undo
	undo_redo_counter int = 0
)

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
			register_curr_state()

		case termbox.KeySpace:
			keylog_message = "Space"
			if current_mode == INSERT {
				insert_character(rune(' '))
				modified = true
			}
			register_curr_state()

		case termbox.KeyEnter:
			keylog_message = "Enter"
			if current_mode == INSERT {
				insert_newline()
				modified = true
			}
			register_curr_state()

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			keylog_message = "BackSpace"
			if current_mode == INSERT {
				delete_character()
				modified = true
			}
			register_curr_state()

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
			register_curr_state()

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
			if is_delimiter(ev.Ch) {
				register_curr_state()
			}

		} else if current_mode == NORMAL {
			keylog_message = string(ev.Ch)

			switch ev.Ch {

			case 'q':
				QUIT = true
				return

			case 'e', 'i':
				current_mode = INSERT

			case '?', ':':
				current_mode = PROMPT

			case 'r': // reload
				textBuffer = nil
				read_file(source_file)
				current_mode = NORMAL
				modified = false
				register_curr_state()

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
				register_curr_state()

			case 'd':
				delete_word(-1)
				modified = true
				register_curr_state()

			case 'C':
				delete_word(1)
				current_mode = INSERT
				modified = true
				register_curr_state()

			case 'c':
				delete_word(-1)
				current_mode = INSERT
				modified = true
				register_curr_state()

			case 'u':
				undo()
				modified = true
				undo_redo_counter = 0

			case 'U':
				redo()
				modified = true

				// first redo (keep stack for potential subsequent redos)
				if strings.Compare(last_modification_key, "U") != 0 && undo_redo_counter == 0 {
					undo_redo_counter++

					// not first redo & last motion was not redo
					// (reset stack since buffer might be changed)
				} else if strings.Compare(last_modification_key, "U") != 0 {
					redo_stack = nil // clear redo stack
					undo_redo_counter = 0
				}
			}

		}
	}
	if len(textBuffer) > 0 && currentCol > len(textBuffer[currentRow]) {
		currentCol = len(textBuffer[currentRow])
	}

	if !is_configuration_evt(ev) &&
		!is_io_evt(ev) &&
		!is_navigation_evt(ev) {

		last_modification_key = string(ev.Ch) // update last key
	}

}
