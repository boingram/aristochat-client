package aristochat

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
)

func init() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Unable to open logrus.log for writing, exiting!\n %v", err))
	}
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.DebugLevel)
}
