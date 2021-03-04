package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/chrischdi/xrandr-gotoggle/internal/config"
	"github.com/chrischdi/xrandr-gotoggle/pkg/xrandr"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func applyCmd() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply command",
		Long:  "Tries to apply the saved configuration if available.",
		RunE: func(cmd *cobra.Command, args []string) error {
			klog.V(4).Info("getting screens via xrandr")
			screens, err := xrandr.GetScreens()
			if err != nil {
				return err
			}

			klog.V(4).Info("getting screen checksum")
			if len(screens) < screenID+1 {
				return fmt.Errorf("screen %d was not found", screenID)
			}
			id := screens[screenID].GetCurrentID()
			klog.V(3).Infof("got checksum: %s", id)

			argsFromConfig, err := config.GetArgsForID(configPath, id)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			if len(argsFromConfig) == 0 {
				klog.V(0).Infof("unable to apply configuration: no configuration found for id %q", id)
				return nil
			}
			klog.V(3).Infof("got configuration from config: %s", strings.Join(argsFromConfig, " "))

			klog.V(4).Info("calculating current xrandr configuration")
			xrandrArgs := []string{}
			for _, screen := range screens {
				xrandrArgs = append(xrandrArgs, screen.GetXrandrArgs()...)
			}
			klog.V(3).Infof("got xrandr configuration: %s", strings.Join(xrandrArgs, " "))

			if reflect.DeepEqual(argsFromConfig, xrandrArgs) {
				klog.Info("configuration is already applied")
				return nil
			}

			if dryRun {
				klog.Infof("[dry-run] executing command: %s", config.CmdLine(argsFromConfig))
				return nil
			}
			return xrandr.Run(argsFromConfig)
		},
	}

	applyCmd.Flags().SortFlags = false
	return applyCmd
}
