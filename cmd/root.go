package cmd

import (
	"github.com/shihtzu-systems/bingo/pkg/fsystem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var version = "0.0.0"
var datestamp = "20101010"
var timestamp = "013042"

var configPath string
var logger *zap.Logger

var rootCmd = &cobra.Command{
	Use: "bingo",
}

func Execute() error {

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is .bingo.yaml)")

	return rootCmd.Execute()
}

func init() {
	logger, _ = zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))
	cobra.OnInitialize(onInitialize)
}

func onInitialize() {
	dout, err := fsystem.ReadFsFile("app.datestamp")
	if err != nil {
		logger.Fatal("unable to read file: app.datestamp", zap.Error(err))
	}
	datestamp = string(dout)

	tout, err := fsystem.ReadFsFile("app.timestamp")
	if err != nil {
		logger.Fatal("unable to read file: app.timestamp", zap.Error(err))
	}
	timestamp = string(tout)

	vout, err := fsystem.ReadFsFile("app.version")
	if err != nil {
		logger.Fatal("unable to read file: app.version", zap.Error(err))
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
		logger.Debug("using config file:", zap.Field{
			Type:   zapcore.StringType,
			Key:    "config file",
			String: viper.ConfigFileUsed(),
		})
	} else {
		logger.Warn("error while reading config file:", zap.Field{
			Type:   zapcore.StringType,
			Key:    "config file",
			String: viper.ConfigFileUsed(),
		}, zap.Error(err))
	}
}

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running command", zap.Error(err))
	}
	return err
}
