package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

// unlinkDeviceBindGroup represents the unlink command
var unlinkDeviceBindGroup = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink a device bind group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 参数
		groupid := args[0]

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 取消关联
		err = user.UDriver.UnLinkDeviceBindGroup(token, groupid)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unlink device bind group error:", err)
			return nil
		}

		// 成功
		fmt.Printf("Successfully unlink device bind group '%v'\n", groupid)

		return nil
	},
}

func init() {
	deviceBindGroupCmd.AddCommand(unlinkDeviceBindGroup)
}
