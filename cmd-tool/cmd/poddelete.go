package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

// deletePodCmd represents the delete command
var deletePodCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the edge pod",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 删除Pod
		for i := 0; i < len(args); i++ {
			if err := user.UDriver.DeleteEdgePod(token, clientID, args[i]); err != nil {
				fmt.Fprintf(os.Stderr, "Delete pod \"%v\" error: %v\n", args[i], err)
				return nil
			}
			fmt.Printf("Delete pod \"%v\" success\n", args[i])
		}

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(deletePodCmd)

	deletePodCmd.Flags().StringVarP(&clientID, "clientid", "c", "", "specify the edge client id")
	deletePodCmd.MarkFlagRequired("clientid")
}
