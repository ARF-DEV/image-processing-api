package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateImagesTable, downCreateImagesTable)
}

func upCreateImagesTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	sq := `CREATE TABLE images (
		id SERIAL PRIMARY KEY,
		url VARCHAR(255) NOT NULL
	)`

	if _, err := tx.ExecContext(ctx, sq); err != nil {
		return err
	}

	fmt.Println("images up")
	return nil
}

func downCreateImagesTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	sq := `DROP TABLE images`
	if _, err := tx.ExecContext(ctx, sq); err != nil {
		return err
	}
	return nil
}
