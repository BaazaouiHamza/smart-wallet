package db

import (
	"embed"
	"fmt"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
)

//go:embed migration
var migrationFS embed.FS

func MigrateUp(connectionString string) error {
	driver, err := postgres.New(connectionString)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	if _, err := migration.Migrate(driver, &migration.EmbedMigrationSource{
		EmbedFS: migrationFS,
		Dir:     "migration",
	}, migration.Up, 0); err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	return nil
}
