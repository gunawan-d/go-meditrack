package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

var DbConnection *sql.DB

// DBMigrate runs database migrations using sql-migrate
func DBMigrate(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migrations",
	}

	n, err := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	DbConnection = dbParam

	fmt.Printf("Migration Success, applied %d migrations!\n", n)
}