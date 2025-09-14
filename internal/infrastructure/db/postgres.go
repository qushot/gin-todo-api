package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// DBConn はデータベース接続を表す
var conn *pgx.Conn

// Initialize はデータベース接続を初期化する
func Initialize(connectionString string) (*pgx.Conn, error) {
	var err error
	conn, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	slog.Info("Database connected")
	return conn, nil
}

// GetDBConn はデータベース接続を取得する
func GetDBConn() *pgx.Conn {
	return conn
}

// CloseDB はデータベース接続を閉じる
func CloseDB(ctx context.Context) error {
	if conn != nil {
		return conn.Close(ctx)
	}
	return nil
}
