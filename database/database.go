package database

import (
	"crg.eti.br/go/mooca/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	DB  *sqlx.DB
	cfg *config.Config
}
