//go:build !unix

package caps

// cellAspectFromFd has no portable pixel-size query off Unix; callers default
// to a 2:1 cell when this returns 0.
func cellAspectFromFd(uintptr) float64 { return 0 }
