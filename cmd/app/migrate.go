package app

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
	Long:  `Run database migrations up or down.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run migrations up",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration(func(m *migrate.Migrate) error {
			return m.Up()
		})
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Run migrations down",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration(func(m *migrate.Migrate) error {
			return m.Down()
		})
	},
}

func runMigration(direction func(*migrate.Migrate) error) {
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migrate instance")
	}

	if err := direction(m); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	log.Info().Msg("Migrations completed successfully")
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	rootCmd.AddCommand(migrateCmd)
}
