package database

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(
	ctx context.Context,
	pool *pgxpool.Pool,
) error {
	// Создаём техническую таблицу,
	// в которой храним номера уже выполненных миграций.
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY
		)
	`

	_, err := pool.Exec(
		ctx,
		query,
	)
	if err != nil {
		return err
	}

	// Получаем список файлов из папки migrations.
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	// Сортируем файлы по имени:
	// 001 -> 002 -> 003 -> ...
	sort.Slice(
		entries,
		func(i, j int) bool {
			return entries[i].Name() <
				entries[j].Name()
		},
	)

	for _, entry := range entries {
		// Папки пропускаем.
		if entry.IsDir() {
			continue
		}

		// Получаем имя текущего файла.
		name := entry.Name()

		// Обрабатываем только .sql файлы.
		if filepath.Ext(name) != ".sql" {
			continue
		}

		// Ожидаем имя вида:
		// 001_create_clients.sql
		parts := strings.SplitN(
			name,
			"_",
			2,
		)

		if len(parts) != 2 {
			continue
		}

		// Получаем номер миграции.
		// Например "003" -> 3.
		version, err := strconv.Atoi(
			parts[0],
		)
		if err != nil {
			return err
		}

		// Проверяем, выполнялась ли уже эта миграция.
		checkQuery := `
			SELECT EXISTS (
				SELECT 1
				FROM schema_migrations
				WHERE version = $1
			)
		`

		var applied bool

		err = pool.QueryRow(
			ctx,
			checkQuery,
			version,
		).Scan(
			&applied,
		)
		if err != nil {
			return err
		}

		// Если версия уже есть в schema_migrations,
		// повторно её не выполняем.
		if applied {
			continue
		}

		// Собираем полный путь к SQL-файлу.
		migrationPath := filepath.Join(
			"migrations",
			name,
		)

		// Читаем SQL из файла.
		sqlBytes, err := os.ReadFile(
			migrationPath,
		)
		if err != nil {
			return err
		}

		sqlQuery := string(sqlBytes)

		// Начинаем транзакцию.
		// SQL миграции и запись её версии
		// должны сохраниться вместе.
		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}

		// Выполняем SQL миграции.
		_, err = tx.Exec(
			ctx,
			sqlQuery,
		)
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		// Записываем номер успешно выполненной миграции.
		_, err = tx.Exec(
			ctx,
			`
				INSERT INTO schema_migrations (version)
				VALUES ($1)
			`,
			version,
		)
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		// Подтверждаем все изменения.
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}

	return nil
}
