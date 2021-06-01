package monitor

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xpzouying/sysmonitor/alarm"
)

func Start(cmd *cobra.Command, args []string) error {
	if err := serve(); err != nil {
		logrus.Errorf("serve error: %v", err)
		return err
	}

	return nil
}

func serve() error {
	alarmMsg := make(chan string, 128)

	alarmer, err := newAlarm()
	if err != nil {
		return err
	}
	alarmer.Notify(alarmMsg)

	interval := viper.GetInt("SysMonitor.interval")
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := doMetrics(alarmMsg); err != nil {
				logrus.Errorf("do_metrics: %v", err)
			}
		}
	}

	return nil
}

func newAlarm() (alarm.Alarm, error) {
	switch alarmType := viper.GetString("Alarm.type"); alarmType {
	case "weixin":
		hubAPI := viper.GetString("Alarm.send_api")
		towx := viper.GetString("Alarm.towx")
		return alarm.NewWeixin(hubAPI, towx)
	}

	return nil, errors.New("unknown alarm type")
}

func doMetrics(alarmMsgs chan string) error {
	cpuPercent, memPercent, err := getCPUAndMemoryMetrics()
	if err != nil {
		return err
	}

	logrus.Infof("cpu total_percent: %v", cpuPercent)
	logrus.Infof("mem total_percent: %v", memPercent)

	cpuThr := viper.GetFloat64("AlarmRules.CPU")
	MemThr := viper.GetFloat64("AlarmRules.Mem")

	if cpuPercent <= cpuThr && memPercent <= MemThr {
		return nil
	}

	msg := fmt.Sprintf("cpu percent: %v\nmem percent: %v", cpuPercent, memPercent)
	alarmMsgs <- msg

	return nil
}

func getCPUAndMemoryMetrics() (cpuPercent, memPercent float64, err error) {
	interval := viper.GetInt("SysMonitor.interval")

	cpuPercents, err := cpu.Percent(time.Duration(interval)*time.Second, false)
	if err != nil {
		return 0, 0, errors.Wrap(err, "get cpu percent error")
	}
	cpuPercent = cpuPercents[0]

	memp, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, errors.Wrap(err, "get mem virtual memory error")
	}
	memPercent = memp.UsedPercent

	return
}
