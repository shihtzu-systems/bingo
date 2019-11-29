package cmd

import (
	"github.com/shihtzu-systems/bingo/pkg/fsystem"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var version = "0.0.0"
var datestamp = "20101010"
var timestamp = "013042"

var configPath string

var rootCmd = &cobra.Command{
	Use: "bingo",
}

func Execute() error {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is .bingo.yaml)")

	return rootCmd.Execute()
}

func initConfig() {
	dout, err := fsystem.ReadFsFile("app.datestamp")
	if err != nil {
		log.Fatal(err)
	}
	datestamp = string(dout)

	tout, err := fsystem.ReadFsFile("app.timestamp")
	if err != nil {
		log.Fatal(err)
	}
	timestamp = string(tout)

	vout, err := fsystem.ReadFsFile("app.version")
	if err != nil {
		log.Fatal(err)
	}
	version = string(vout)

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".bingo")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("BINGO")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Warn("error occurred while reading config: ", err)
	}
}
