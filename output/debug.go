package output

import (
	"handheldui/vars"
	"log"
)

func Printf(format string, a ...any) (n int, err error) {
	if vars.Config.Logs {
		log.Printf(format, a...)
		return len(format), nil
	}
	return 0, nil
}

func Errorf(format string, a ...any) (err error) {
	if vars.Config.Logs {
		log.Printf("ERROR: "+format, a...)
		return nil
	}
	return nil
}

func Sprintf(format string, a ...any) string {
	if vars.Config.Logs {
		log.Printf(format, a...)
		return format
	}
	return ""
}
