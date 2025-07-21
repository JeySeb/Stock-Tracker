package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

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
