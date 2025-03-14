package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// listPodCmd represents the list command
var listPodCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the list of edge pods you apply",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取设备列表
		edgepods, err := user.UDriver.ListEdgePod(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "List user edge pods error:", err)
			return nil
		}

		// 设备列表为空
		if len(*edgepods) < 1 {
			fmt.Printf("User %v doesn't apply any edge pods.\n", userlogininfo.UserName)
			return nil
		}

		// 输出设备列表
		maxNumberLen := 6
		maxClientIDLen := 8
		maxNodenameLen := 8
		maxNameLen := 4
		maxReadyLen := 5
		maxRestartsLen := 8
		maxCreateTimeLen := 10
		for _, edgepod := range *edgepods {
			if len(edgepod.ClientID) > maxClientIDLen {
				maxClientIDLen = len(edgepod.ClientID)
			}
			if len(edgepod.NodeName) > maxNodenameLen {
				maxNodenameLen = len(edgepod.NodeName)
			}
			if len(edgepod.Name) > maxNameLen {
				maxNameLen = len(edgepod.Name)
			}
			if len(edgepod.Ready) > maxReadyLen {
				maxReadyLen = len(edgepod.Ready)
			}
		}

		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s\n", maxNumberLen, "NUMBER", maxClientIDLen, "CLIENTID", maxNodenameLen, "NODENAME", maxNameLen, "NAME", maxReadyLen, "READY", maxRestartsLen, "RESTARTS", maxCreateTimeLen, "CREATETIME")

		for number, edgepod := range *edgepods {
			restarts := 0
			for _, container := range edgepod.Containers {
				if container.RestartCount > restarts {
					restarts = container.RestartCount
				}
			}

			createTime := time.Unix(0, edgepod.CreateTime)
			createTimeStr := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", createTime.Year(), createTime.Month(), createTime.Day(), createTime.Hour(), createTime.Minute(), createTime.Second())

			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*d  %-[13]*s\n", maxNumberLen, number+1, maxClientIDLen, edgepod.ClientID, maxNodenameLen, edgepod.NodeName, maxNameLen, edgepod.Name, maxReadyLen, edgepod.Ready, maxRestartsLen, restarts, maxCreateTimeLen, createTimeStr)
		}

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(listPodCmd)
}
