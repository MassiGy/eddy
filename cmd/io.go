package main

import (
	"bufio"
	"os"
)

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
