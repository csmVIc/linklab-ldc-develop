package cmd

import (
	"fmt"
	"os"

	"linklab/device-control-v2/base-library/user"

	"github.com/spf13/cobra"
)

// listIoTCmd represents the list command
var listIoTCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the list of devices you use",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return err
		}

		// 获取设备列表
		devices, err := user.UDriver.ListIoTDevice(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "List user devices error:", err)
			return err
		}

		// 设备列表为空
		if len(*devices) < 1 {
			fmt.Printf("User %v doesn't use any devices.\n", userlogininfo.UserName)
			return nil
		}

		// 输出设备列表
		maxNumberLen := 6
		maxClientIDLen := 8
		maxDeviceIDLen := 8
		for _, device := range *devices {
			if len(device.ClientID) > maxClientIDLen {
				maxClientIDLen = len(device.ClientID)
			}
			if len(device.DeviceID) > maxDeviceIDLen {
				maxDeviceIDLen = len(device.DeviceID)
			}
		}
		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s\n", maxNumberLen, "NUMBER", maxClientIDLen, "CLIENTID", maxDeviceIDLen, "DEVICEID")
		for number, device := range *devices {
			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s\n", maxNumberLen, number+1, maxClientIDLen, device.ClientID, maxDeviceIDLen, device.DeviceID)
		}

		return nil
	},
}

func init() {
	iotnodeCmd.AddCommand(listIoTCmd)
}
