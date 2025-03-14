package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

var container string

// getPodLogCmd represents the logs command
var getPodLogCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get pod container log",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取日志
		wsHandler, err := user.UDriver.GetEdgePodLog(token, clientID, args[0], container)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Get pod container log error: %v\n", err)
			return nil
		}

		// 输出日志
		for {
			resp := response.Response{}
			if err := wsHandler.ReadJSON(&resp); err != nil {
				break
			}
			if resp.Code != 0 {
				fmt.Fprintf(os.Stderr, "%v", resp.Msg)
				break
			}
			fmt.Println(resp.Msg)
		}

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(getPodLogCmd)

	getPodLogCmd.Flags().StringVarP(&clientID, "clientid", "c", "", "specify the edge client id")
	getPodLogCmd.MarkFlagRequired("clientid")
	getPodLogCmd.Flags().StringVarP(&container, "container", "", "", "specify the edge pod container")
}
