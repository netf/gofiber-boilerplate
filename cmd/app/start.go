package app

import (
	"github.com/netf/gofiber-boilerplate/internal/api"
	"github.com/netf/gofiber-boilerplate/internal/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.NewDatabase(cfg.DatabaseURL)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not initialize database")
		}

		if cfg.Environment == "development" {
			log.Info().Msg("Running in development mode. Applying auto-migrations...")
			if err := db.AutoMigrate(database); err != nil {
				log.Fatal().Err(err).Msg("Could not run auto-migrations")
			}
		}

		api.Start(cfg, database)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
