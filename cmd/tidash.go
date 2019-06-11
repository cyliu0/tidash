package cmd

import (
	"os"
	"time"

	"github.com/cyliu0/tidash/dashboard"
	"github.com/cyliu0/tidash/pd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tidashCmd = &cobra.Command{
	Use:  "tidash",
	Long: "A terminal dashboard for monitoring TiKV cluster leaders",
	Run: func(cmd *cobra.Command, args []string) {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Fatalf("Failed to set log level, err: %v", err)
		}
		logrus.SetLevel(level)
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Fatalf("Failed to open log file: %v, err: %v", logFile, err)
		}
		logrus.SetOutput(file)
		pd.InitPDClient(pdApiAddr)
		dashboard.InitDash(time.Duration(updateInterval) * time.Second)
	},
}

func Execute() {
	if err := tidashCmd.Execute(); err != nil {
		logrus.Fatalf("Command err: " + err.Error())
	}
}

var pdApiAddr string
var updateInterval int
var logLevel string
var logFile string

func init() {
	tidashCmd.Flags().StringVarP(&pdApiAddr, "pd-api-addr", "p", "", "PD API address")
	tidashCmd.Flags().IntVarP(&updateInterval, "update-interval", "i", 1, "Dashboard update interval in seconds")
	tidashCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level")
	tidashCmd.Flags().StringVarP(&logFile, "log-file", "f", "tidash.log", "Log File")
}
