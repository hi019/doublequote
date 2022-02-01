package main

import (
	"context"
	"os"
	"os/signal"

	dq "doublequote"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the application",
	RunE:  RunServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func RunServe(_ *cobra.Command, _ []string) error {
	cfg := dq.Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	_, cleanup, err := initializeApplication(&cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Wait for CTRL+C.
	<-ctx.Done()

	return nil
}
