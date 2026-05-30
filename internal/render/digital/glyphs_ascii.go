package digital

// ASCII5 is the fallback font for terminals that can't render block glyphs
// (TERM=dumb, non-UTF-8 locale) or when font:"ascii" is requested. Digits are
// 3 cells wide; the colon is 1 cell wide.
var ASCII5 = map[rune][]string{
	'0': {"###", "# #", "# #", "# #", "###"},
	'1': {"  #", "  #", "  #", "  #", "  #"},
	'2': {"###", "  #", "###", "#  ", "###"},
	'3': {"###", "  #", "###", "  #", "###"},
	'4': {"# #", "# #", "###", "  #", "  #"},
	'5': {"###", "#  ", "###", "  #", "###"},
	'6': {"###", "#  ", "###", "# #", "###"},
	'7': {"###", "  #", "  #", "  #", "  #"},
	'8': {"###", "# #", "###", "# #", "###"},
	'9': {"###", "# #", "###", "  #", "###"},
	':': {" ", ":", " ", ":", " "},
}
