package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xpzouying/sysmonitor/monitor"
)

var (
	rootCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the system monitor",
		RunE:  monitor.Start,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
