package main

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/ARF-DEV/image-processing-api/migrations/scripts"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	args := os.Args // 2 migrations path, 3 db type, 4 db connection string, 5 command
	if len(args) < 5 {
		panic("not enough arguments, go run main.go <dir> <db type> <db connection string> <command>")
	}

	dir := args[1]
	dbType := args[2]
	dbConn := args[3]
	command := args[4]

	db, err := sql.Open(dbType, dbConn)
	if err != nil {
		panic(err)
	}

	if err := goose.RunContext(context.Background(), command, db, dir); err != nil {
		panic(err)
	}
}
