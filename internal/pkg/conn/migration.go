package conn

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-pg/pg"
	"github.com/gobuffalo/packr"
	"github.com/rubenv/sql-migrate"

	// is a driver needed by sql.Open to connect to postgresql.
	_ "github.com/lib/pq"
)

// MigrateSync will sync all database structure from migration file to postgres.
// This uses lib/pq because sql-migrate only support sql.DB struct.
// Actually we can create struct that implements all method like sql.DB do, but it takes time.
func MigrateSync(url string) error {
	opts, err := pg.ParseURL(url)
	if err != nil {
		return err
	}

	host := "localhost"
	port := "5432"

	// splitting by :, because addr is hostname:port, and we must split that
	addr := strings.Split(opts.Addr, ":")
	if len(addr) == 2 {
		if strings.TrimSpace(addr[0]) != "" {
			host = addr[0]
		}

		port = addr[1]
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, opts.User, opts.Password, opts.Database)

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()
	if err != nil {
		return err
	}

	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./../../../assets/migrations"),
	}

	migrate.SetTable("migrations")
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}

// MigrateReset will do migration down.
func MigrateReset(url string) error {
	opts, err := pg.ParseURL(url)
	if err != nil {
		return err
	}

	host := "localhost"
	port := "5432"

	// splitting by :, because addr is hostname:port, and we must split that
	addr := strings.Split(opts.Addr, ":")
	if len(addr) == 2 {
		if strings.TrimSpace(addr[0]) != "" {
			host = addr[0]
		}

		port = addr[1]
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, opts.User, opts.Password, opts.Database)

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()
	if err != nil {
		return err
	}

	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./../../../assets/migrations"),
	}

	migrate.SetTable("migrations")
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Down)
	return err
}
