//go:build unix

package caps

import "golang.org/x/sys/unix"

// cellAspectFromFd returns the cell height/width ratio from the terminal's
// reported pixel size, or 0 when the terminal does not report pixels.
func cellAspectFromFd(fd uintptr) float64 {
	ws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	if err != nil || ws == nil || ws.Xpixel == 0 || ws.Ypixel == 0 || ws.Col == 0 || ws.Row == 0 {
		return 0
	}
	cellW := float64(ws.Xpixel) / float64(ws.Col)
	cellH := float64(ws.Ypixel) / float64(ws.Row)
	if cellW == 0 {
		return 0
	}
	return cellH / cellW
}
