package db

import (
	"context"
	"fmt"
	"os"
	"utube/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool
var dbContext context.Context

func Init() error {
	dbContext = context.Background()
	var err error
	db, err = pgxpool.New(dbContext, getConnString())

	return err
}

func CloseDb() {
	utils.Log.Println("Closing connection/s to db...")
	db.Close()

	utils.Log.Println("Closed connection/s to db")
}

func getConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), "localhost", "5432", os.Getenv("DB_DBNAME"))
}
