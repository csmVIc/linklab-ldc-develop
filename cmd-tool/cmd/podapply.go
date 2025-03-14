package cmd

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/user"
	"os"

	"github.com/spf13/cobra"
)

var podYaml string
var useEdgeRegistry bool
var createIngress bool

// applyPodCmd represents the apply command
var applyPodCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the pod configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 读取文件
		yamlbinary, err := ioutil.ReadFile(podYaml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read pod yaml %v error: %v\n", podYaml, err)
			return nil
		}

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 上传文件
		yamlhash, err := user.UDriver.UploadEdgePodYaml(token, yamlbinary)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Upload pod yaml error:", err)
			return nil
		}

		// 部署服务
		clientid, ingressMap, err := user.UDriver.ApplyEdgePod(token, yamlhash, useEdgeRegistry, createIngress)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Apply edge pod error:", err)
			return nil
		}

		// 成功部署
		fmt.Printf("Successfully deploy pod to edge client '%v'\n", clientid)

		// 显示Ingress
		if ingressMap != nil && len(ingressMap) > 0 {

			fmt.Println("")

			maxPortLen := 4
			maxIngressLen := 7
			for key, val := range ingressMap {
				if len(key) > maxPortLen {
					maxPortLen = len(key)
				}
				if len(val) > maxIngressLen {
					maxIngressLen = len(val)
				}
			}

			fmt.Printf("%-[1]*s  %-[3]*s\n", maxPortLen, "PORT", maxIngressLen, "INGRESS")

			for key, val := range ingressMap {
				fmt.Printf("%-[1]*s  %-[3]*s\n", maxPortLen, key, maxIngressLen, val)
			}
		}

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(applyPodCmd)

	applyPodCmd.Flags().StringVarP(&podYaml, "filename", "f", "", "specify the pod yaml configuration")
	applyPodCmd.MarkFlagRequired("filename")

	applyPodCmd.Flags().BoolVarP(&useEdgeRegistry, "use-edgeregistry", "", false, "use edge image registry")
	applyPodCmd.Flags().BoolVarP(&createIngress, "create-ingress", "", false, "create container service ingress")
}
