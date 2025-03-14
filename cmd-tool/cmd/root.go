package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"linklab/device-control-v2/base-library/user"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "linklab-cli",
	Short: "linklab command line tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.linklab-cli.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".linklab-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".linklab-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {

		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		// 用户登录初始化
		userlogininfo.UserName = viper.GetString("NAMESPACE")
		userlogininfo.PassWord = viper.Sub("login").GetString("password")

		// 接口初始化
		urlinfo := user.UInfo{}
		if err := viper.Sub("interface").Unmarshal(&urlinfo); err != nil {
			fmt.Fprintln(os.Stderr, "Init config unmarshal error:", err)
		}

		user.UDriver, err = user.New(&urlinfo)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Init config file error:", err)
		}

	}
}
