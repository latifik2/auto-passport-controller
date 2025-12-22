package db

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (db *Database) IsDuplicateHashes(snapshotHash string) (bool, error) {
	var latestHash string
	if err := db.Pool.QueryRow(context.Background(), "SELECT snapshot_hash FROM passports_raw ORDER BY created_at DESC LIMIT 1").Scan(&latestHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Таблица пустая — это не ошибка
			return false, nil
		}

		slog.Error("Error occurred while trying to select latest snapshot", "err", err)
		return true, err
	}

	if snapshotHash == latestHash {
		return true, nil
	} else {
		return false, nil
	}

}

func (db *Database) SelectActualPassports() ([]byte, error) {
	var passportsBytes []byte

	if err := db.Pool.QueryRow(context.Background(), "SELECT passport FROM passports_raw ORDER BY created_at DESC LIMIT 1").Scan(&passportsBytes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Таблица пустая — это не ошибка
			slog.Info("No records in database to display")
			return []byte{}, nil
		}

		errText := "Error occured while trying to select actual passports"

		slog.Error(errText, "err", err)
		return []byte{}, errors.New(errText)
	}

	return passportsBytes, nil
}
