package cmd

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/spf13/cobra"
)

// podExecCmd represents the exec command
var podExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command in a container",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 参数
		pod := args[0]
		commands := args[1:]

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 建立连接
		wsHandler, err := user.UDriver.EdgePodExec(token, clientID, pod, container, commands)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Execute command error: %v\n", err)
			return nil
		}

		// 配置输出
		stdin := streams.NewIn(os.Stdin)
		if err := stdin.SetRawTerminal(); err != nil {
			fmt.Fprintf(os.Stderr, "Set raw terminal error: %v\n", err)
			return nil
		}
		defer stdin.RestoreTerminal()

		// 退出通道
		inputExitChan := make(chan error, 1)
		outputExitChan := make(chan error, 1)

		// 输入
		go func() {
			buf := make([]byte, 1024)
			for {
				len, err := stdin.Read(buf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Stdin read error: %v\n", err)
					inputExitChan <- err
					return
				}
				if len > 0 {
					if err := wsHandler.WriteJSON(&msg.EdgeClientPodExec{
						Type: msg.NormalPodExec,
						Msg:  string(buf[0:len]),
					}); err != nil {
						fmt.Fprintf(os.Stderr, "Websocket write error: %v\n", err)
						inputExitChan <- err
						return
					}
				}
			}
		}()

		// 输出
		go func() {
			for {
				resp := response.Response{}
				if err := wsHandler.ReadJSON(&resp); err != nil {
					// fmt.Fprintf(os.Stderr, "Websocket read error: %v\n", err)
					outputExitChan <- err
					return
				}
				if resp.Code != 0 {
					fmt.Fprintf(os.Stderr, "%v", resp.Msg)
					outputExitChan <- nil
					return
				}
				fmt.Print(resp.Msg)
			}
		}()

		// 监控
		select {
		case <-inputExitChan:
			break
		case <-outputExitChan:
			break
		}

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(podExecCmd)

	podExecCmd.Flags().StringVarP(&clientID, "clientid", "c", "", "specify the edge client id")
	podExecCmd.MarkFlagRequired("clientid")
	podExecCmd.Flags().StringVarP(&container, "container", "", "", "specify the edge pod container")
}
