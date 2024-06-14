package main

import (
	"os"
)

var (
	// binary info(assigned at build time)
	version     string
	binary_name string
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
	LINES_COUNT           int
	SHOW_LINE_NUMBERS     bool = true

	// in buffer cursor info
	currentCol, currentRow int
	offsetX, offsetY       int

	// buffer
	textBuffer = [][]rune{}

	// status bar flags
	modified       bool = false
	err_message    string
	info_message   string
	keylog_message string = "	"
	current_mode   mode

	// read/write to file
	source_file                  string
	SAVE_TO_FILE_MAX_ERROR_COUNT int = 3
)

func main() {
	if os.Getenv("ENV") == "dev" {
		binary_name = "eddy"
		version = "v0.1.0"
	}
	run_editor()
}
