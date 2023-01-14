package database

type Database interface {
	ChkMigrations() (bool, error)
	RunMigrations() error
	HealthCheck() error
}
