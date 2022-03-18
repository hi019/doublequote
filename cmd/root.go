package main

import (
	"fmt"
	"os"

	dq "doublequote"
	"doublequote/pkg/domain"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var cfg domain.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "doublequote",
	Short: "A brief description of your application",
}

// Execute adds all child commands to the root command and sets flags appropriately
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./doublequote.toml)")
}

// initConfig reads in config file and env variables if set
func initConfig() {
	viper.SetConfigName("default")
	viper.SetConfigType("toml")

	defaultCfg, err := dq.Assets.Open("assets/default.toml")
	cobra.CheckErr(err)
	err = viper.ReadConfig(defaultCfg)
	cobra.CheckErr(err)

	// Use the config path passed via argument, if it exists. If it doesn't, then look for
	// a config file in the current directory
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		path, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(path)
		viper.SetConfigName("doublequote")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // Read in environment variables that match

	// If a config file is found, read it in
	if err := viper.MergeInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		cobra.CheckErr(err)
	}
}
