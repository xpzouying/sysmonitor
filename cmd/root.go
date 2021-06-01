package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	if err := loadConfig(); err != nil {
		return err
	}

	return rootCmd.Execute()
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	return viper.ReadInConfig()
}
