package repository

import (
	"context"
	"testing"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
)

func TestGetVocabulariesWithMaxWords(t *testing.T) {
	ctx := context.Background()
	dbName := "vocab_test"
	tp := postgres.NewTempPostgres(ctx, dbName)
	defer tp.DropDB(ctx)

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
