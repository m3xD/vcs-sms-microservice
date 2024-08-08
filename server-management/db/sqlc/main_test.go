package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:2108@localhost:5432/sms?sslmode=disable"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbSource)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	testQueries = New(conn)

	os.Exit(m.Run())
}
