package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/edkadigital/startmeup/config"
	"github.com/edkadigital/startmeup/pkg/log"
	"github.com/edkadigital/startmeup/pkg/migrations"
	"github.com/edkadigital/startmeup/pkg/services"
)

func main() {
	// Parse command line flags
	var (
		migrateRiver   bool
		forceRiver     bool
		migrateSchemas bool
	)
	flag.BoolVar(&migrateRiver, "river", true, "Run River queue migrations")
	flag.BoolVar(&forceRiver, "force-river", false, "Force applying River migrations regardless of what's already applied")
	flag.BoolVar(&migrateSchemas, "schemas", true, "Run Ent schema migrations using Atlas")
	flag.Parse()

	// Start a new container to access the database
	c := services.NewContainer()
	defer func() {
		// Gracefully shutdown all services.
		if err := c.Shutdown(); err != nil {
			log.Default().Error("shutdown failed", "error", err)
		}
	}()

	// Run Ent schema migrations if requested
	if migrateSchemas {
		fmt.Println("Running Ent schema migrations using Atlas...")

		// Construct the connection string for Atlas based on the environment
		var dsn string
		switch c.Config.App.Environment {
		case config.EnvTest:
			dsn = c.Config.Database.TestConnection
		default:
			dsn = c.Config.Database.Connection
		}

		// Validate DSN presence
		if dsn == "" {
			fmt.Println("Error: Database connection string (DSN) is empty for the current environment.")
			os.Exit(1)
		}

		// Determine the path to the migrations directory relative to the workspace root
		migrationDir, err := filepath.Abs("ent/migrate/migrations")
		if err != nil {
			fmt.Printf("Error finding migration directory path: %v\n", err)
			os.Exit(1)
		}

		// Ensure the migration directory exists
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			fmt.Printf("Migration directory does not exist: %s\nRun 'go generate ./ent/...' first to generate migrations.\n", migrationDir)
			os.Exit(1)
		}

		// Prepare the atlas command
		// We use 'file://' prefix for the directory path
		cmd := exec.Command("atlas", "migrate", "apply",
			"--dir", fmt.Sprintf("file://%s", migrationDir),
			"--url", dsn,
		)

		// Capture output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running Atlas migrations: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Ent schema migrations completed!")
	}

	// Run River migrations if requested
	if migrateRiver {
		fmt.Println("Running River queue migrations...")

		// Create a River migrations manager
		riverManager := migrations.NewRiverManager(c.Database)

		var err error
		if forceRiver {
			// Force apply all migrations
			err = riverManager.ApplyAllMigrations(context.Background())
		} else {
			// Apply only pending migrations
			err = riverManager.ApplyPendingMigrations(context.Background())
		}

		if err != nil {
			fmt.Printf("Error running River migrations: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("River migrations completed successfully!")
	}

	fmt.Println("All migrations completed!")
}
