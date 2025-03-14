package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

// podResourceCmd represents the resource command
var podResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Show the resource usage of edge pods you apply",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取设备列表
		edgenodes, err := user.UDriver.ListEdgePodResource(token)
		if err != nil {
			fmt.Fprintln(os.Stderr, "List user edge pods error:", err)
			return nil
		}

		// 设备列表为空
		if len(*edgenodes) < 1 {
			fmt.Printf("User %v doesn't apply any edge pods.\n", userlogininfo.UserName)
			return nil
		}

		// 输出设备列表
		maxNumberLen := 6
		maxClientIDLen := 8
		maxNameLen := 4
		maxCpuCores := 10
		maxMemBytes := 13
		for _, edgenode := range *edgenodes {
			if len(edgenode.ClientID) > maxClientIDLen {
				maxClientIDLen = len(edgenode.ClientID)
			}
			if len(edgenode.Name) > maxNameLen {
				maxNameLen = len(edgenode.Name)
			}
		}

		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s\n", maxNumberLen, "NUMBER", maxClientIDLen, "CLIENTID", maxNameLen, "NAME", maxCpuCores, "CPU(cores)", maxMemBytes, "MEMORY(bytes)")

		for number, edgenode := range *edgenodes {

			cpuCores := int64(0)
			memBytes := int64(0)
			for _, container := range edgenode.Containers {
				cpuCores += container.CpuUse
				memBytes += container.MemUse / 1000 / 1024 / 1024
			}
			cpuCoresStr := fmt.Sprintf("%dm", cpuCores)
			memBytesStr := fmt.Sprintf("%dMi", memBytes)

			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s\n", maxNumberLen, number+1, maxClientIDLen, edgenode.ClientID, maxNameLen, edgenode.Name, maxCpuCores, cpuCoresStr, maxMemBytes, memBytesStr)
		}

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(podResourceCmd)
}
