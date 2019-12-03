package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCommand = &cobra.Command{
	Use:  "version",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("server.v1.sessionKey"))
		fmt.Println(viper.GetString("bingo.v1.theme"))
		fmt.Println(viper.GetString("system.v2.debug"))
		fmt.Printf("%s+on.%s.at.%s", version, datestamp, timestamp)
	},
}

func init() {
	rootCmd.AddCommand(versionCommand)
}
