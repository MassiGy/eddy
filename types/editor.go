package types

import "github.com/gdamore/tcell/v2"

var PrintableKeys []tcell.Key = []tcell.Key{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'y', 'z', 'Y', 'Z', '@', '#', '{', '}', '[', ']', '/', '\\',
	'.', '_', '|', '\'', '"', '&', '²', '(', ')', '=', '^', '+',
	' ', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-',
	'°', '%', '*', '$', '!', ':', ',', '?', '<', '>', '€', '£',
	'~', ';',
}

type JumpPoint struct {
	X       int
	Y       int
	Val     rune
	NewLine bool
}

type JumpList []JumpPoint
