package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/chrischdi/xrandr-gotoggle/pkg/xrandr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	dryRun     bool
	configPath = path.Join(os.Getenv("HOME"), ".config/xrandr-gotoggle.yaml")
	screenID   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "xrandr-gotoggle",
	Short:         "A CLI tool to handle xrandr configurations",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "If true, only print the command to be executed")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", configPath, "Path to the configuration yaml file.")
	rootCmd.PersistentFlags().IntVar(&screenID, "screen-id", 0, "ID of the xrandr screen to handle.")

	rootCmd.AddCommand(
		applyCmd(),
		configCmd(),
		printChecksumCmd(),
	)

	klog.InitFlags(nil)
	flag.Set("skip_headers", "true")
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func printChecksumCmd() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "print-checksum",
		Short: "print-checksum command",
		Long:  "Calculates and prints the screen checksum.",
		RunE: func(cmd *cobra.Command, args []string) error {
			klog.V(4).Info("getting edid checksum")
			screens, err := xrandr.GetScreens()
			if err != nil {
				return err
			}

			if len(screens) < screenID+1 {
				return fmt.Errorf("screen %d was not found", screenID)
			}

			id := screens[screenID].GetCurrentID()
			klog.Infof("got screen checksum: %s", id)
			return nil
		},
	}

	applyCmd.Flags().SortFlags = false
	return applyCmd
}
