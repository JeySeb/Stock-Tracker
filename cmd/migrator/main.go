package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	var direction = flag.String("direction", "up", "Direction to run migrations (up or down)")
	var steps = flag.Int("steps", 0, "Number of steps to run migrations (0 for all)")
	var force = flag.Int("force", -1, "Force migration to specific version (fixes dirty state)")
	var specificMigration = flag.String("migration", "", "Run specific migration by number (e.g., '004'). Use with -direction or 'reset' to run down then up")
	flag.Parse()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set and is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to the database")

	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{})
	if err != nil {
		log.Fatal("Failed to create CockroachDB driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"cockroachdb",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}

	// Handle force option to fix dirty database state
	if *force >= 0 {
		log.Printf("Forcing migration version to %d", *force)
		if err := m.Force(*force); err != nil {
			log.Fatal("Failed to force migration version:", err)
		}
		log.Println("Migration version forced successfully")
		return
	}

	// Handle specific migration execution
	if *specificMigration != "" {
		if err := runSpecificMigration(m, *specificMigration, *direction); err != nil {
			log.Fatal("Specific migration failed:", err)
		}
		return
	}

	// Handle regular migration execution
	switch *direction {
	case "up":
		if *steps > 0 {
			log.Printf("Running %d steps forward", *steps)
			err = m.Steps(*steps)
		} else {
			log.Println("Running all steps forward")
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			log.Printf("Running %d steps backward", *steps)
			err = m.Steps(-*steps)
		} else {
			log.Println("Running all steps backward")
			err = m.Down()
		}
	default:
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", *direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed:", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply - database is up to date")
	} else {
		log.Println("Migration completed successfully")
	}
}

// runSpecificMigration handles execution of a specific migration by number
func runSpecificMigration(m *migrate.Migrate, migrationNumber, direction string) error {
	// Parse migration number
	migrationNum, err := strconv.Atoi(migrationNumber)
	if err != nil {
		return fmt.Errorf("invalid migration number '%s': must be a valid integer", migrationNumber)
	}

	// Get current migration version
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in dirty state, please run with -force flag to fix")
	}

	log.Printf("Current migration version: %d", currentVersion)
	log.Printf("Target migration: %03d", migrationNum)

	switch direction {
	case "reset":
		log.Printf("🔄 Resetting migration %03d (down then up)", migrationNum)

		// First, migrate down to the migration before our target
		targetDownVersion := uint(migrationNum - 1)
		log.Printf("⬇️  Migrating down to version %d", targetDownVersion)
		if err := m.Migrate(targetDownVersion); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to migrate down to version %d: %w", targetDownVersion, err)
		}

		// Then migrate up to our target migration
		targetUpVersion := uint(migrationNum)
		log.Printf("⬆️  Migrating up to version %d", targetUpVersion)
		if err := m.Migrate(targetUpVersion); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to migrate up to version %d: %w", targetUpVersion, err)
		}

		log.Printf("✅ Migration %03d reset completed successfully", migrationNum)

	case "up":
		targetVersion := uint(migrationNum)
		log.Printf("⬆️  Running migration %03d UP", migrationNum)
		if err := m.Migrate(targetVersion); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to migrate to version %d: %w", migrationNum, err)
		}
		log.Printf("✅ Migration %03d UP completed successfully", migrationNum)

	case "down":
		if migrationNum <= 0 {
			return fmt.Errorf("cannot migrate down to version %d: must be positive", migrationNum)
		}
		targetVersion := uint(migrationNum - 1)
		log.Printf("⬇️  Running migration %03d DOWN (migrating to version %d)", migrationNum, targetVersion)
		if err := m.Migrate(targetVersion); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to migrate down from version %d: %w", migrationNum, err)
		}
		log.Printf("✅ Migration %03d DOWN completed successfully", migrationNum)

	default:
		return fmt.Errorf("invalid direction for specific migration: %s. Use 'up', 'down', or 'reset'", direction)
	}

	// Show final version
	finalVersion, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Warning: Could not get final version: %v", err)
	} else {
		log.Printf("Final migration version: %d", finalVersion)
	}

	return nil
}
