package postgres

import (
	"crg.eti.br/go/mooca/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB  *sqlx.DB
	cfg *config.Config
}

func New(cfg *config.Config) (*Database, error) {
	db := &Database{
		cfg: cfg,
	}
	err := db.Open()

	return db, err
}

func (b *Database) Open() error {
	db, err := sqlx.Connect("postgres", b.cfg.DatabaseURL)
	if err != nil {
		return err
	}

	b.DB = db
	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
