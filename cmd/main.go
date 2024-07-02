package main

import (
	"runtime"
)

type mode int

const (
	NORMAL mode = iota
	INSERT
	PROMPT
)

const (
	// constants
	WRAP                   bool   = true //TODO: add this to a config file
	WRAP_AFTER             int    = 80
	UNKOWN_SOURCE_FILENAME string = "nofile"
)

var (
	// editor vars
	QUIT                  bool = false
	ROWS, COLS            int
	LINE_NUMBER_COL_WIDTH int
	SHOW_LINE_NUMBERS     bool = true

	// in buffer cursor info
	currentCol, currentRow int
	offsetX, offsetY       int
	main_curr_col          int
	main_curr_row          int
	prompt_mode_curr_col   int
	prompt_mode_curr_row   int

	// buffers
	textBuffer         [][]rune
	main_buffer        = [][]rune{}
	prompt_mode_buffer = [][]rune{}

	// undo redo
	undo_stack             [][][]rune
	main_undo_stack        = [][][]rune{}
	prompt_mode_undo_stack = [][][]rune{}

	redo_stack             [][][]rune
	main_redo_stack        = [][][]rune{}
	prompt_mode_redo_stack = [][][]rune{}

	last_modification_key string

	// status bar flags
	modified       bool = false
	err_message    string
	info_message   string
	keylog_message string = ""
	current_mode   mode

	// read/write to file
	source_file                  string
	SAVE_TO_FILE_MAX_ERROR_COUNT int = 3

	// OS based divergences
	TAB string = ""
)

func main() {
	if runtime.GOOS == "windows" {
		TAB = "    " // 4 spaces
		keylog_message = "    "
	} else {
		TAB = "\t" // linux & darwin
		keylog_message = "\t"
	}
	run_editor()
}
