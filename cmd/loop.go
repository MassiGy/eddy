package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

func run_editor() {

	undo_stack = main_undo_stack
	redo_stack = main_redo_stack

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
		textBuffer = get_homescreen_text()
		source_file = UNKOWN_SOURCE_FILENAME
		info_message = "Viewing Eddy's Homescreen."
	}

	main_buffer = textBuffer // for the first file read

	currentCol = 0
	currentRow = 0
	current_mode = NORMAL

	register_curr_state() // first snapshot

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
		if current_mode == PROMPT {
			prompt_mode_buffer = textBuffer
			prompt_mode_undo_stack = undo_stack
			prompt_mode_redo_stack = redo_stack
			prompt_mode_curr_col = currentCol
			prompt_mode_curr_row = currentRow
		} else {
			main_buffer = textBuffer
			main_undo_stack = undo_stack
			main_redo_stack = redo_stack
			main_curr_col = currentCol
			main_curr_row = currentRow
		}
	}
}
