package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

var groupType string

// createDeviceBindGroup represents the create command
var createDeviceBindGroup = &cobra.Command{
	Use:   "create",
	Short: "Create a device bind group",
	Args:  cobra.RangeArgs(1, 5),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 参数
		boardtypes := args

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 创建
		err = user.UDriver.CreateDeviceBindGroup(token, groupType, boardtypes)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Create device bind group error:", err)
			return nil
		}

		// 成功创建
		fmt.Printf("Successfully create device bind group '%v'\n", groupType)

		return nil
	},
}

func init() {
	deviceBindGroupCmd.AddCommand(createDeviceBindGroup)

	createDeviceBindGroup.Flags().StringVarP(&groupType, "type", "t", "", "specify the device bind group type")
	createDeviceBindGroup.MarkFlagRequired("type")
}
