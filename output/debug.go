package output

import (
	"fmt"
	"handheldui/vars"
)

func Printf(format string, a ...any) (n int, err error) {
	if vars.Debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
