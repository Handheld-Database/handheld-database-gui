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

func Errorf(format string, a ...any) (err error) {
	if vars.Debug {
		return fmt.Errorf(format, a...)
	}
	return nil
}

func Sprintf(format string, a ...any) string {
	if vars.Debug {
		return fmt.Sprintf(format, a...)
	}
	return ""
}
