package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	// "github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	// "github.com/testcontainers/testcontainers-go"
	"github.com/docker/compose/v2/pkg/api"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	// "github.com/testcontainers/testcontainers-go/wait"
)

func TestGetVocabulariesWithMaxWords(t *testing.T) {
	ctx := context.Background()
	dbName := "vocab_test"
	dc, err := tc.NewDockerCompose("../../../../../deploy/docker-compose.db.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = dc.Down(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()
	err = dc.WithEnv(map[string]string{
		"DB_NAME": dbName,
	}).WaitForService("migration", wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5*time.Second)).
		Up(ctx, tc.WithRecreate(api.RecreateForce), tc.Wait(true))
	if err != nil {
		t.Fatal(err)
	}
	pgxConn, err := dc.ServiceContainer(ctx, "postgres")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(pgxConn)

	// ctx := context.Background()
	// tp := postgres.NewTempPostgres(ctx, "vocab_test")
	// defer tp.DropDB(ctx)

	// rows, err := pgxConn.Query(ctx, `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users`)
	// if err != nil {
	// 	t.Errorf("failed to start container: %s", err)
	// 	return
	// }

	// var users []User
	// user := User{}
	// for rows.Next() {
	// 	err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.LastVisitAt, &user.CreatedAt)
	// 	if err != nil {
	// 		t.Errorf("failed to start container: %s", err)
	// 		break
	// 	}

	// 	users = append(users, user)
	// }

	// if len(users) != 1 {
	// 	t.Errorf("expected to get 1 user, got %d", len(users))
	// }
}
