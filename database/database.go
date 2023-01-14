package database

type Database interface {
	HealthCheck() error
	Open() error
	Close() error
}
