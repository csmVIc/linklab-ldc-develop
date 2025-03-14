package cmd

import "github.com/spf13/cobra"

// deviceBindGroupCmd represents the group command
var deviceBindGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "group related commands",
}

func init() {
	iotnodeCmd.AddCommand(deviceBindGroupCmd)
}
