package main

import (
	"os"

	"k8s.io/klog"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
