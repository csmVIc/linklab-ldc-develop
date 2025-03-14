package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"linklab/device-control-v2/base-library/tool"
	"linklab/device-control-v2/base-library/user"

	"github.com/spf13/cobra"
)

var imageSourcePath, buildArchitecture string

// imageBuildCmd represents the build command
var imageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a container image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 参数
		imagename := args[0]

		// 绝对路径
		var err error
		imageSourcePath, err = filepath.Abs(imageSourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Get image source directory '%v' absolute path error: %v\n", imageSourcePath, err)
			return nil
		}

		// 检查目录
		dockerfilePath := filepath.Join(imageSourcePath, "Dockerfile")
		if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "'Dockerfile' does not exist in directory '%v' error\n", imageSourcePath)
			return nil
		}

		// 压缩文件
		zipBinary, err := tool.ZipDirectory(imageSourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Zip image source directory '%v' error: %v\n", imageSourcePath, err)
			return nil
		}

		// 获取token
		token, err := user.UDriver.UserLogin(userlogininfo.UserName, userlogininfo.PassWord)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Get user token error:", err)
			return nil
		}

		// 上传文件
		filehash, err := user.UDriver.UploadEdgeImageSource(token, zipBinary)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Upload image source error:", err)
			return nil
		}

		// 指令集架构选择
		nodeselector := map[string]string{}
		if len(buildArchitecture) > 0 {
			nodeselector["kubernetes.io/arch"] = buildArchitecture
		}

		// 部署服务
		clientid, err := user.UDriver.BuildEdgeImage(token, filehash, imagename, nodeselector)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Build edge image error:", err)
			return nil
		}

		// 成功部署
		fmt.Printf("Successfully deploy build image pod to edge client '%v'\n", clientid)

		return nil
	},
}

func init() {
	edgepodCmd.AddCommand(imageBuildCmd)

	imageBuildCmd.Flags().StringVarP(&imageSourcePath, "path", "p", "", "specify the image source directory")
	imageBuildCmd.MarkFlagRequired("path")

	imageBuildCmd.Flags().StringVarP(&buildArchitecture, "arch", "", "", "specify the image architecture")
}
