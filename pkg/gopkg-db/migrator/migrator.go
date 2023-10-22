package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func New(host string, port int, user, password, name string) (*Migrator, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	dir, err := getMigrationsDir()
	if err != nil {
		return nil, err
	}
	return &Migrator{
		driver: driver,
		dir:    dir,
	}, nil
}

func getMigrationsDir() (string, error) {
	dirList := []string{"internal/migrations", "../../internal/migrations", "migrations"}
	for _, dir := range dirList {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			return dir, nil
		}
	}
	return "", fmt.Errorf("cant find migrations dir in: %v", dirList)
}

type Migrator struct {
	driver database.Driver
	dir    string
}

func (m *Migrator) Up() error {
	log.Printf("Rolling up migration from dir: %v\n", m.dir)
	mig, err := migrate.NewWithDatabaseInstance("file://"+m.dir, "postgres", m.driver)
	if err != nil {
		return err
	}
	err = mig.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Database doesn't change\n")
		return nil
	}
	return err
}
