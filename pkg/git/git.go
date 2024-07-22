package git

import (
	"fmt"
	"os/exec"
)

func GitPull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = "website"

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, output)
	}

	return nil
}
