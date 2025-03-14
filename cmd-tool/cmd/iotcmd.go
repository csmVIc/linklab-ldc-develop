package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

var clientID, deviceID string

// iotCmd represents the cmd command
var iotCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Send commands to the device serial port",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return err
		}

		for i := 0; i < len(args); i++ {
			// 发送命令
			err = user.UDriver.SendIoTCmd(token, clientID, deviceID, args[i]+"\r\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Send command \"%v\" error: %v\n", args[i], err)
				return err
			}
			fmt.Printf("Send command \"%v\" success\n", args[i])
		}

		return nil
	},
}

func init() {
	iotnodeCmd.AddCommand(iotCmd)

	iotCmd.Flags().StringVarP(&clientID, "clientid", "c", "", "specify the device management client id")
	iotCmd.MarkFlagRequired("clientid")

	iotCmd.Flags().StringVarP(&deviceID, "deviceid", "d", "", "specify the device id")
	iotCmd.MarkFlagRequired("deviceid")
}
