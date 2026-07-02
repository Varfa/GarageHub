package database

import (
	"context"
	"fmt"

	"github.com/Varfa/GarageHub/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(dbCfg config.DatabaseConfig) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
		dbCfg.SSLMode,
	)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil

}
