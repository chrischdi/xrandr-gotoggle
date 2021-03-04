package main

import (
	"fmt"

	"github.com/chrischdi/xrandr-gotoggle/internal/config"
	"github.com/chrischdi/xrandr-gotoggle/pkg/xrandr"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func configCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "config command",
		Long:  "The subcommand containing configuration related commands.",
	}

	configCmd.AddCommand(
		printCurrentConfigCmd(),
		setCurrentConfigCmd(),
		viewConfigCmd(),
	)
	return configCmd
}

func printCurrentConfigCmd() *cobra.Command {
	printCurrentConfigCmd := &cobra.Command{
		Use:   "print-current-config",
		Short: "print current configuration",
		Long:  "Prints the current configuration which could be added to the configuration file.",
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

			klog.V(4).Info("calculating current xrandr configuration")
			xrandrArgs := []string{}
			for _, screen := range screens {
				xrandrArgs = append(xrandrArgs, screen.GetXrandrArgs()...)
			}
			c := config.Config{}
			c[id] = xrandrArgs
			klog.Infof("Printing xrandr-gotoggle config entry:")
			klog.Infof(config.YamlLine(id, xrandrArgs))
			klog.Infoln()
			klog.Infof("Printing xrandr command:")
			klog.Infof(config.CmdLine(xrandrArgs))
			return nil
		},
	}

	printCurrentConfigCmd.Flags().SortFlags = false
	return printCurrentConfigCmd
}
func setCurrentConfigCmd() *cobra.Command {
	setCurrentConfigCmd := &cobra.Command{
		Use:   "set-current-config",
		Short: "set current configuration at the configuration file",
		Long:  "Adds or overwrites the current configuration to the configuration file.",
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

			klog.V(4).Info("calculating current xrandr configuration")
			xrandrArgs := []string{}
			for _, screen := range screens {
				xrandrArgs = append(xrandrArgs, screen.GetXrandrArgs()...)
			}

			cfg, err := config.ReadConfig(configPath)
			if err != nil {
				return err
			}

			cfg[id] = xrandrArgs

			return config.WriteConfig(configPath, cfg)
		},
	}

	setCurrentConfigCmd.Flags().SortFlags = false
	return setCurrentConfigCmd
}

func viewConfigCmd() *cobra.Command {
	viewConfigCmd := &cobra.Command{
		Use:   "view",
		Short: "view config",
		Long:  "Prints the current configuration from the configuration file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.ReadConfig(configPath)
			if err != nil {
				return err
			}

			for id, args := range cfg {
				klog.Info(config.YamlLine(id, args))
			}
			return nil
		},
	}

	viewConfigCmd.Flags().SortFlags = false
	return viewConfigCmd
}
