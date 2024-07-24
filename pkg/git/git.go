package git

import (
	"fmt"
	"os/exec"
)

const repoURL = "https://github.com/kubernetes/website.git"

func GitClone() error {
	cmd := exec.Command("git", "clone", repoURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %v, output: %s", err, output)
	}

	return nil
}

func GitPull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = "website"

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, output)
	}

	return nil
}
