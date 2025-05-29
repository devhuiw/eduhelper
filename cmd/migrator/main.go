package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		migrationsPath  string
		migrationsTable string
		dbUser          string
		dbPassword      string
		dbHost          string
		dbPort          string
		dbName          string
		down            bool
		step            int
	)

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "name of the migrations table")
	flag.StringVar(&dbUser, "db-user", "root", "MySQL user")
	flag.StringVar(&dbPassword, "db-password", "", "MySQL password")
	flag.StringVar(&dbHost, "db-host", "localhost", "MySQL host")
	flag.StringVar(&dbPort, "db-port", "3306", "MySQL port")
	flag.StringVar(&dbName, "db-name", "", "MySQL database name")
	flag.BoolVar(&down, "down", false, "revert all migrations (down to version 0)")
	flag.IntVar(&step, "step", 0, "migrate up/down N steps. Use negative for down, positive for up.")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	if dbName == "" {
		panic("db-name is required")
	}

	dsn := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	if migrationsTable != "" {
		dsn = fmt.Sprintf("%s&x-migrations-table=%s", dsn, migrationsTable)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		panic(err)
	}

	if step != 0 {
		fmt.Printf("migrating %d steps...\n", step)
		if err := m.Steps(step); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
				return
			}
			panic(err)
		}
		fmt.Println("migration steps applied successfully")
		return
	}

	if down {
		fmt.Println("reverting all migrations (down to version 0)...")
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to revert")
				return
			}
			panic(err)
		}
		fmt.Println("all migrations reverted successfully")
		return
	}

	fmt.Println("applying migrations (up)...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
