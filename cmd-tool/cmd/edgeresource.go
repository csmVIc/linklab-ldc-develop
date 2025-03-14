package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

// edgeResourceCmd represents the resource command
var edgeResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Show the resource usage of edge nodes",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 获取设备列表
		edgenodes, err := user.UDriver.ListEdgeResource(token)
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
		maxCpuCores := 10
		maxCpuPercent := 4
		maxMemBytes := 13
		maxMemPercent := 7
		// maxGpu := 3
		for _, edgenode := range *edgenodes {
			if len(edgenode.ClientID) > maxClientIDLen {
				maxClientIDLen = len(edgenode.ClientID)
			}
			if len(edgenode.Name) > maxNameLen {
				maxNameLen = len(edgenode.Name)
			}
		}

		// fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s  %-[15]*s\n", maxNumberLen, "NUMBER", maxClientIDLen, "CLIENTID", maxNameLen, "NAME", maxCpuCores, "CPU(cores)", maxCpuPercent, "CPU%", maxMemBytes, "MEMORY(bytes)", maxMemPercent, "MEMORY%", maxGpu, "GPU")
		fmt.Printf("%-[1]*s  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s\n", maxNumberLen, "NUMBER", maxClientIDLen, "CLIENTID", maxNameLen, "NAME", maxCpuCores, "CPU(cores)", maxCpuPercent, "CPU%", maxMemBytes, "MEMORY(bytes)", maxMemPercent, "MEMORY%")

		for number, edgenode := range *edgenodes {

			cpuPercent := float64(edgenode.CpuUse) / float64(edgenode.CpuAll) * 100
			memBytes := edgenode.MemUse / 1000 / 1024 / 1024
			memPercent := float64(edgenode.MemUse) / float64(edgenode.MemAll) * 100

			cpuCoresStr := fmt.Sprintf("%dm", edgenode.CpuUse)
			cpuPercentStr := fmt.Sprintf("%.f%%", cpuPercent)
			memBytesStr := fmt.Sprintf("%dMi", memBytes)
			memPercentStr := fmt.Sprintf("%.f%%", memPercent)

			// fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s  %-[15]*d\n", maxNumberLen, number+1, maxClientIDLen, edgenode.ClientID, maxNameLen, edgenode.Name, maxCpuCores, cpuCoresStr, maxCpuPercent, cpuPercentStr, maxMemBytes, memBytesStr, maxMemPercent, memPercentStr, maxGpu, int(edgenode.NvidiaGpuAll))
			fmt.Printf("%-[1]*d  %-[3]*s  %-[5]*s  %-[7]*s  %-[9]*s  %-[11]*s  %-[13]*s\n", maxNumberLen, number+1, maxClientIDLen, edgenode.ClientID, maxNameLen, edgenode.Name, maxCpuCores, cpuCoresStr, maxCpuPercent, cpuPercentStr, maxMemBytes, memBytesStr, maxMemPercent, memPercentStr)
		}

		return nil
	},
}

func init() {
	edgenodeCmd.AddCommand(edgeResourceCmd)
}
