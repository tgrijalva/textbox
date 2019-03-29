// +build !windows

package terminal

import (
	"os"

	"golang.org/x/sys/unix"
)

func TerminalSize() (width, height int, err error) {
	size, err := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(size.Col), int(size.Row), nil
}
