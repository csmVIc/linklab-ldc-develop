package cmd

import "github.com/spf13/cobra"

// edgepodCmd represents the edgepod command
var edgepodCmd = &cobra.Command{
	Use:   "edgepod",
	Short: "edgepod related commands",
}

func init() {
	rootCmd.AddCommand(edgepodCmd)
}
