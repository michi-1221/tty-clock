package digital

// Block5 is the default 5-row full-block font (tty-clock derived, widened to
// 2 cells per stroke). Each digit glyph is 6 cells wide; the colon is 2 cells
// wide. Every row of a single glyph shares that glyph's width. There is no
// space glyph: the colon-off blank is generated at the colon's width (★5).
var Block5 = map[rune][]string{
	'0': {"██████", "██  ██", "██  ██", "██  ██", "██████"},
	'1': {"    ██", "    ██", "    ██", "    ██", "    ██"},
	'2': {"██████", "    ██", "██████", "██    ", "██████"},
	'3': {"██████", "    ██", "██████", "    ██", "██████"},
	'4': {"██  ██", "██  ██", "██████", "    ██", "    ██"},
	'5': {"██████", "██    ", "██████", "    ██", "██████"},
	'6': {"██████", "██    ", "██████", "██  ██", "██████"},
	'7': {"██████", "    ██", "    ██", "    ██", "    ██"},
	'8': {"██████", "██  ██", "██████", "██  ██", "██████"},
	'9': {"██████", "██  ██", "██████", "    ██", "██████"},
	':': {"  ", "██", "  ", "██", "  "},
}
