package xrandr

import (
	"crypto/sha256"
	"fmt"

	"k8s.io/klog"
)

func (s *Screen) GetXrandrArgs() []string {
	cmd := []string{}
	cmd = append(cmd, "--screen", fmt.Sprintf("%d", s.No))
	for _, monitor := range s.Monitors {
		cmd = append(cmd, "--output", monitor.ID)
		if monitor.Connected {
			cmd = append(cmd, "--mode", monitor.Resolution.String())
			cmd = append(cmd, "--pos", monitor.Position.String())
			if monitor.Primary {
				cmd = append(cmd, "--primary")
			}
		} else {
			cmd = append(cmd, "--off")
		}
	}
	return cmd
}

func (s *Screen) GetCurrentID() string {
	var txt string
	for _, monitor := range s.Monitors {
		klog.Infof("handling monitor %s", monitor.ID)
		if monitor.Connected {
			klog.Infof("adding monitor %s", monitor.ID)
			txt = txt + monitor.ID + "\n"
		}
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(txt)))
}
