package wrappers

import (
	"fmt"
	"handheldui/output"
	"os"
	"os/exec"
)

// UnzipFile calls the system to unzip the file and provides progress information
func UnzipFile(src, dest string) error {
	// Check if the unzip command is available on the system
	_, err := exec.LookPath("unzip")
	if err != nil {
		return fmt.Errorf("unzip command not found on system: %v", err)
	}

	// Prepare the unzip command to extract the file
	output.Sprintf("%s %s %s %s", "unzip", "-o", src, "-d", dest)
	cmd := exec.Command("unzip", "-o", src, "-d", dest)

	// Redirect stdout and stderr to monitor the progress
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the unzip command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing unzip: %v", err)
	}

	return nil
}
