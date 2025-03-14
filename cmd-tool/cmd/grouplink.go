package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/user"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// linkDeviceBindGroup represents the link command
var linkDeviceBindGroup = &cobra.Command{
	Use:   "link",
	Short: "Link a device bind group",
	Args:  cobra.RangeArgs(1, 5),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 参数
		devices := []request.DevInfoForGroup{}
		for index, arg := range args {
			tmp := strings.Split(arg, ":")
			if len(tmp) != 2 {
				fmt.Fprintf(os.Stderr, "Devices[%v] '%v' format error\n", index, arg)
				return nil
			}

			devices = append(devices, request.DevInfoForGroup{
				ClientID: tmp[0],
				DeviceID: tmp[1],
			})
		}

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 关联
		groupid, err := user.UDriver.LinkDeviceBindGroup(token, groupType, devices)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Link device bind group error:", err)
			return nil
		}

		// 成功
		fmt.Printf("Successfully link device bind group '%v'\n", groupid)

		return nil
	},
}

func init() {
	deviceBindGroupCmd.AddCommand(linkDeviceBindGroup)

	linkDeviceBindGroup.Flags().StringVarP(&groupType, "type", "t", "", "specify the device bind group type")
	linkDeviceBindGroup.MarkFlagRequired("type")
}
