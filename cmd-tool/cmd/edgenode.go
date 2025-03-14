package cmd

import "github.com/spf13/cobra"

// edgenodeCmd represents the edgenode command
var edgenodeCmd = &cobra.Command{
	Use:   "edgenode",
	Short: "edgenode related commands",
}

func init() {
	rootCmd.AddCommand(edgenodeCmd)
}
