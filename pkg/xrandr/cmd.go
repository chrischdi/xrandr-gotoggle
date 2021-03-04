package xrandr

import (
	"fmt"
	"os/exec"
)

const xrandrPath = "/usr/bin/xrandr"

func Run(args []string) error {
	cmd := exec.Command(xrandrPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xrandr exited with %v: %s", err, out)
	}
	return err
}

// GetScreens returns all the screens info from xrandr output
func GetScreens() (Screens, error) {
	cmd := exec.Command(xrandrPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	screens, err := parseScreens(string(output))
	if err != nil {
		return nil, err
	}

	return screens, nil
}
