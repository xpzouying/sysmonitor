package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/sysmonitor/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		logrus.Errorln(err)
	}
}
