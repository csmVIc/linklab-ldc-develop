package main

import (
	"linklab/device-control-v2/cmd-tool/cmd"
)

func init() {
	// log.SetLevel(log.DebugLevel)
	// log.SetReportCaller(true)
	// log.SetFormatter(&log.TextFormatter{
	// 	DisableColors: true,
	// })
}

func main() {
	cmd.Execute()
}
