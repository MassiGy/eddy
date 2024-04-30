package types

import "github.com/gdamore/tcell/v2"

var PrintableKeys []tcell.Key = []tcell.Key{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
	'y', 'z', '@', '#', '{', '}', '[', ']', '/', '\\', '~', '-',
	'.', '_', '|', '\'', '"', '&', '²', '(', ')', '=', '^', '+',
	' ',
}

type JumpPoint struct {
	X   int
	Y   int
	Val rune
}

type JumpList []JumpPoint
