package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Initialize configures the MySQL/MariaDB database for Vuyo and ensures that the models have
// been fully migrated.
func Initialize() error {
	db, _ := sql.Open("mysql", "user:password@tcp(host:port)/dbname?multiStatements=true")
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"mysql",
		driver,
	)

	m.Steps(2)

	return nil
}
