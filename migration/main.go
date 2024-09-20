package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	statusCmd = "status"
	upCmd     = "up"
	downCmd   = "down"
)

func main() {
	log.Printf("migration start")

	cmd := flag.String("cmd", "status", "migration command")
	connString := flag.String("url", "", "connection postgres URL")
	flag.Parse()

	*connString = prepareConnString(*connString)

	connConfig, err := pgx.ParseConfig(*connString)
	if err != nil {
		log.Fatal(err, "parse config")
	}

	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err, "open db")
	}

	if err := retryPing(db, 3, 30*time.Second); err != nil {
		log.Fatal(err, "ping db")
	}

	m, err := migrate.New("file://migrations", *connString)
	if err != nil {
		log.Fatal("can't create migrate instance", err)
	}

	log.Printf("migration start with command '%s' with connection string '%s'", *cmd, *connString)

	switch {
	case strings.EqualFold(*cmd, upCmd):
		err = m.Up()
	case strings.EqualFold(*cmd, downCmd):
		err = m.Down()
	default:
		log.Fatal("cmd is not fit for 'migration'")
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Println("migration finished")
}

//goland:noinspection SpellCheckingInspection
func retryPing(db *sql.DB, retries int, timeout time.Duration) error {
	for i := 1; i <= retries; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("%d ping DB error: %v", i, err)
			time.Sleep(timeout)
			continue
		}
		return nil
	}
	return errors.New("DB unreacheble")
}

func prepareConnString(connString string) string {
	const prefix = "postgresql://"
	if strings.Index(connString, prefix) == 0 {
		return connString
	}
	connString = prefix + connString
	return connString
}
