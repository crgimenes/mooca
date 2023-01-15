package postgres

import (
	"errors"
	"log"

	"crg.eti.br/go/mooca/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	// currentMigration is the current migration version of the code.
	// it must be incremented every time a new migration is added.
	currentMigration = 1
)

var (
	ErrDatabaseAhead       = errors.New("database is ahead of current migration, please update the application")
	ErrDatabaseNotUpToDate = errors.New("database is not up to date, please run migrations")

	// go:embed migration01.sql
	migration01 string
)

type Database struct {
	db  *sqlx.DB
	cfg *config.Config
}

func New(cfg *config.Config) (*Database, error) {
	db := &Database{
		cfg: cfg,
	}
	err := db.open()

	return db, err
}

func (b *Database) open() error {
	db, err := sqlx.Connect("postgres", b.cfg.DatabaseURL)
	if err != nil {
		return err
	}

	b.db = db
	return nil
}

func (d *Database) close() error {
	return d.db.Close()
}

func (d *Database) createMigrationTable() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
				id INTEGER PRIMARY KEY,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				)`)
	return err
}

func (d *Database) VerifyMigration() (int, error) {
	var (
		lastMigration int
		count         int
	)
	err := d.db.Get(&count, `SELECT COUNT(*) FROM migrations`)
	if err != nil {
		return 0, err
	}
	log.Printf("migrations: %d", count)
	if count != 0 {
		err = d.db.Get(&lastMigration, "SELECT MAX(id) as max FROM migrations")
		if err != nil {
			return 0, err
		}
	}

	return lastMigration, nil
}

func (d *Database) ChkMigration() error {
	lastMigration, err := d.VerifyMigration()
	if err != nil {
		return err
	}

	if lastMigration < currentMigration {
		return ErrDatabaseNotUpToDate
	}

	if lastMigration > currentMigration {
		return ErrDatabaseAhead
	}

	return nil
}

func (d *Database) RunMigration() error {
	err := d.createMigrationTable()
	if err != nil {
		return err
	}

	lastMigration, err := d.VerifyMigration()
	if err != nil {
		return err
	}

	log.Printf("last migration: %d", lastMigration)

	// begin transaction
	tx, err := d.db.Beginx()
	// run migrations
	switch lastMigration {
	case 0:
		log.Println("running migration 1")
		_, err = tx.Exec(migration01)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		// update migration table
		_, err = tx.Exec(`INSERT INTO migrations (id) VALUES (1)`)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		log.Println("done migration 1")
		lastMigration = 1

		fallthrough
	default:
		log.Println("no migrations to run")
	}

	if currentMigration != lastMigration {
		_ = tx.Rollback()

		// this should never happen... ok it can happen if you forget
		// to update the currentMigration variable.
		log.Fatal("currentMigration variable is not up to date")
	}

	return tx.Commit()
}
