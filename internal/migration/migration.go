package migration

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

func New() *migrate.Migrations {
	return migrations
}

//go:embed *sql

var sqlMigrations embed.FS

func init() {
	migrations.DiscoverCaller()
	if err := migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}

}
