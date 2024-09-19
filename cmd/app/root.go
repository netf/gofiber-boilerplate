package app

import (
	"fmt"
	"os"

	"github.com/netf/gofiber-boilerplate/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Golang Fiber Boilerplate API",
	Long:  `A boilerplate for building RESTful APIs using Golang and the Fiber web framework.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
}
