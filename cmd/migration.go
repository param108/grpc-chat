package cmd

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	// Required to call the init function. Dependency of postgres migrate
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"os"
)

const ()

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Perform DB migrations",
	RunE:  migrateCmdF,
}

var RollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Perform DB rollbacks",
	RunE:  rollbackCmdF,
}

func migrateCmdF(command *cobra.Command, args []string) error {
	//logger.SetupLogger(*config.LogSettings.LogLevel)
	return runDatabaseMigrations()
}

func rollbackCmdF(command *cobra.Command, args []string) error {
	//logger.SetupLogger(*config.LogSettings.LogLevel)
	return runDatabaseRollback()
}

func createMigrate() (*migrate.Migrate, error) {

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	path := os.Getenv("MIGRATION_PATH")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username,
		password, host, port,
		dbname)

	fmt.Println(url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		fmt.Printf("failed to load the database: %s", err)
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Printf("ping to the database host failed: %s", err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		fmt.Printf("failed to prepare migration: %s", err)
		return nil, err
	}

	return m, nil
}

func runDatabaseMigrations() error {
	m, err := createMigrate()
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Printf("no changes needed")
			return nil
		}
		fmt.Printf("migration failed: %s", err)
		return err
	}

	fmt.Printf("migration successful")
	return nil
}

func runDatabaseRollback() error {
	m, err := createMigrate()
	if err != nil {
		return err
	}
	err = m.Steps(-1)
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Printf("no changes needed %s", err.Error())
			return nil
		}
		fmt.Printf("rollback failed: %s", err)
		return err
	}

	fmt.Printf("rollback successful")
	return nil
}
