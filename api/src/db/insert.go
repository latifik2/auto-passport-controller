package db

import (
	"context"
	"fmt"
	"log/slog"
)

func (db *Database) InsertRawJSON(snapshotHash string, rawJson []byte) {
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO passports_raw (snapshot_hash, passport) VALUES ($1, $2)", snapshotHash, rawJson,
	)

	if err != nil {
		slog.Error(fmt.Sprintf("Failed to insert json into passports_raw table: %v", err))
	}
}
