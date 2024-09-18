package repository

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetVocabulariesWithMaxWords(t *testing.T) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16.1-alpine3.19",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		t.Errorf("failed to start container: %s", err)
		return
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Errorf("failed to get connection string: %s", err)
		return
	}

	slog.Info("connection string: " + connStr)

	pgxConn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		t.Errorf("failed to start container: %s", err)
		return
	}

	_, err = pgxConn.Exec(ctx, `CREATE TABLE IF NOT EXISTS
    "users" (
        "id" UUID PRIMARY KEY,
        "name" TEXT NOT NULL,
        "email" TEXT NOT NULL,
        "password_hash" TEXT NOT NULL,
        "role" TEXT NOT NULL,
        "last_visit_at" TIMESTAMP NOT NULL,
        "created_at" TIMESTAMP NOT NULL
    );`)
	if err != nil {
		t.Errorf("failed to start container: %s", err)
		return
	}
	_, err = pgxConn.Exec(ctx, `INSERT INTO users (id, name, email, password_hash, role, last_visit_at, created_at)
VALUES ('23cc06c9-73d3-40b1-9e8a-6f80db183e7a',
        'admin',
        'ugolkov.prog@gmail.com',
        '$2a$11$oP15pJXtp2ErbHWvGN05ouiMphIzrf8yXJEHkmtf.25JzgWFRaQO6',
        'admin',
        now(),
        now())
ON CONFLICT
    DO NOTHING;`)
	if err != nil {
		t.Errorf("failed to start container: %s", err)
		return
	}

	rows, err := pgxConn.Query(ctx, `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users`)
	if err != nil {
		t.Errorf("failed to start container: %s", err)
		return
	}

	var users []User
	user := User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.LastVisitAt, &user.CreatedAt)
		if err != nil {
			t.Errorf("failed to start container: %s", err)
			break
		}

		users = append(users, user)
	}

	if len(users) != 1 {
		t.Errorf("expected to get 1 user, got %d", len(users))
	}

	fmt.Println(users)
}

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	LastVisitAt  time.Time
	CreatedAt    time.Time
}
