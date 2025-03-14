package cmd

import (
	"github.com/spf13/cobra"
)

// iotnodeCmd represents the iotnode command
var iotnodeCmd = &cobra.Command{
	Use:   "iotnode",
	Short: "iotnode related commands",
}

func init() {
	rootCmd.AddCommand(iotnodeCmd)
}
