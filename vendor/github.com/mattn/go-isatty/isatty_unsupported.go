// +build dragonfly nacl plan9

package isatty

// IsTerminal return true if the file descriptor is terminal.
func IsTerminal(fd uintptr) bool {
	return false
}
