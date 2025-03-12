package migrations

import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
    migrate *migrate.Migrate
}

func NewMigrator(databaseURL string) (*Migrator, error) {
    m, err := migrate.New(
        "file://migrations",
        databaseURL,
    )
    if err != nil {
        return nil, err
    }

    return &Migrator{migrate: m}, nil
}

func (m *Migrator) Up() error {
    return m.migrate.Up()
}

func (m *Migrator) Down() error {
    return m.migrate.Down()
}

func (m *Migrator) Version() (uint, bool, error) {
    return m.migrate.Version()
}

func (m *Migrator) Close() error {
    return m.migrate.Close()
}