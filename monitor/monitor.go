package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Start(cmd *cobra.Command, args []string) error {
	totalPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		logrus.Errorf("get cpu info: %v", err)
		return err
	}
	logrus.Infof("total_percent: %v", totalPercent)

	return nil
}
