package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listDefineDeviceBindGroup represents the listdefine command
var listDefineDeviceBindGroup = &cobra.Command{
	Use:   "listdefine",
	Short: "List the list of defined device bind groups",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取列表
		groups, err := user.UDriver.ListDefineDeviceBindGroup(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "List define device bind groups error:", err)
			return nil
		}

		// 列表为空
		if len(*groups) < 1 {
			fmt.Printf("Defined device bind group list empty.\n")
			return nil
		}

		// 输出列表
		maxNumberLen := 6
		maxTypeLen := 4
		maxBoardsLen := 6
		boardsJoin := []string{}
		for index, elem := range *groups {
			boardsJoin = append(boardsJoin, strings.Join(elem.Boards, ","))
			if len(elem.Type) > maxTypeLen {
				maxTypeLen = len(elem.Type)
			}
			if len(boardsJoin[index]) > maxBoardsLen {
				maxBoardsLen = len(boardsJoin[index])
			}
		}

		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s\n", maxNumberLen, "NUMBER", maxTypeLen, "TYPE", maxBoardsLen, "BOARDS")

		for index, elem := range *groups {
			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s\n", maxNumberLen, index+1, maxTypeLen, elem.Type, maxBoardsLen, boardsJoin[index])
		}

		return nil
	},
}

// listLinkDeviceBindGroup represents the listlink command
var listLinkDeviceBindGroup = &cobra.Command{
	Use:   "listlink",
	Short: "List the list of linked device bind groups",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取列表
		groups, err := user.UDriver.ListLinkDeviceBindGroup(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "List link device bind groups error:", err)
			return nil
		}

		// 列表为空
		if len(*groups) < 1 {
			fmt.Printf("Linked device bind group list empty.\n")
			return nil
		}

		// 输出列表
		maxNumberLen := 6
		maxIDLen := 2
		maxDevicesLen := 7
		devicesJoin := []string{}
		for index, elem := range *groups {

			sb := strings.Builder{}
			for jndex, delem := range elem.Devices {
				if jndex > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(delem.ClientID)
				sb.WriteString(":")
				sb.WriteString(delem.DeviceID)
			}
			devicesJoin = append(devicesJoin, sb.String())

			if len(elem.ID) > maxIDLen {
				maxIDLen = len(elem.ID)
			}
			if len(devicesJoin[index]) > maxDevicesLen {
				maxDevicesLen = len(devicesJoin[index])
			}
		}

		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s\n", maxNumberLen, "NUMBER", maxIDLen, "ID", maxDevicesLen, "DEVICES")

		for index, elem := range *groups {
			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s\n", maxNumberLen, index+1, maxIDLen, elem.ID, maxDevicesLen, devicesJoin[index])
		}

		return nil
	},
}

func init() {
	deviceBindGroupCmd.AddCommand(listDefineDeviceBindGroup)
	deviceBindGroupCmd.AddCommand(listLinkDeviceBindGroup)
}
