package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

// listEdgeCmd represents the list command
var listEdgeCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the list of edge nodes",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取设备列表
		edgenodes, err := user.UDriver.ListEdgeNode(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "List edge nodes error:", err)
			return nil
		}

		// 设备列表为空
		if len(*edgenodes) < 1 {
			fmt.Printf("Edge nodes list empty.\n")
			return nil
		}

		// 输出设备列表
		maxNumberLen := 6
		maxClientIDLen := 8
		maxNameLen := 4
		maxReadyLen := 5
		maxArchLen := 4
		maxOSLen := 2
		maxOSImageLen := 7
		maxIpLen := 2
		for _, edgenode := range *edgenodes {
			if len(edgenode.ClientID) > maxClientIDLen {
				maxClientIDLen = len(edgenode.ClientID)
			}
			if len(edgenode.Name) > maxNameLen {
				maxNameLen = len(edgenode.Name)
			}
			if len(edgenode.Ready) > maxReadyLen {
				maxReadyLen = len(edgenode.Ready)
			}
			if len(edgenode.Architecture) > maxArchLen {
				maxArchLen = len(edgenode.Architecture)
			}
			if len(edgenode.OS) > maxOSLen {
				maxOSLen = len(edgenode.OS)
			}
			if len(edgenode.OSImage) > maxOSImageLen {
				maxOSImageLen = len(edgenode.OSImage)
			}
			if len(edgenode.IpAddress) > maxIpLen {
				maxIpLen = len(edgenode.IpAddress)
			}
		}

		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s  %-[15]*s\n", maxNumberLen, "NUMBER", maxClientIDLen, "CLIENTID", maxNameLen, "NAME", maxReadyLen, "READY", maxArchLen, "ARCH", maxOSLen, "OS", maxOSImageLen, "OSIMAGE", maxIpLen, "IP")

		for number, edgenode := range *edgenodes {
			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s  %-[15]*s\n", maxNumberLen, number+1, maxClientIDLen, edgenode.ClientID, maxNameLen, edgenode.Name, maxReadyLen, edgenode.Ready, maxArchLen, edgenode.Architecture, maxOSLen, edgenode.OS, maxOSImageLen, edgenode.OSImage, maxIpLen, edgenode.IpAddress)
		}

		return nil
	},
}

func init() {
	edgenodeCmd.AddCommand(listEdgeCmd)
}
