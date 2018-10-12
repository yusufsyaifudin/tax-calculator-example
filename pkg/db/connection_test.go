package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
	"github.com/smartystreets/goconvey/convey"
)

const (
	PostgreImage = "bitnami/postgresql"
	PostgreTag   = "10.5.0-debian-9"

	PostgresDB                 = "my_database"
	PostgresUser               = "my_user"
	PostgresPassword           = "my_password"
	PostgreReplicationUser     = "repl_user"
	PostgreReplicationPassword = "repl_password"
)

type user struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var myUser = &user{
	Username: "test_user",
	Password: "test_password",
}

var myUser2 = &user{
	Username: "test_user2",
	Password: "test_password2",
}

var myUser3 = &user{
	Username: "test_user3",
	Password: "test_password3",
}

var conf *Config

// do an insert
var upsertQuery = `INSERT INTO users(username, password) VALUES(?, ?) ON CONFLICT(username) DO UPDATE SET username = EXCLUDED.username RETURNING *;`
var insertQuery = `INSERT INTO users(username, password) VALUES(?, ?) RETURNING *;`
var getQuery = `SELECT * FROM users WHERE username = ? LIMIT 1;`

// Do migration setup.
func TestMain(m *testing.M) {
	log.Warn().Msg("creating docker connection")
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Error().Msgf("Failed to create docker connection: %v", err)
	}

	// pulls an image, creates a container based on it and runs it
	postgreMaster, err := dockerPool.Run(PostgreImage, PostgreTag, []string{
		"POSTGRESQL_REPLICATION_MODE=master",
		fmt.Sprintf("POSTGRESQL_USERNAME=%s", PostgresUser),
		fmt.Sprintf("POSTGRESQL_PASSWORD=%s", PostgresPassword),
		fmt.Sprintf("POSTGRESQL_DATABASE=%s", PostgresDB),
		fmt.Sprintf("POSTGRESQL_REPLICATION_USER=%s", PostgreReplicationUser),
		fmt.Sprintf("POSTGRESQL_REPLICATION_PASSWORD=%s", PostgreReplicationPassword),
	})
	if err != nil {
		log.Error().Msgf("Could not start postgre master: %s", err)
	}

	var dbDSNMaster = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		PostgresUser, PostgresPassword, postgreMaster.GetBoundIP(PostgreImage), postgreMaster.GetPort("5432/tcp"), PostgresDB)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := dockerPool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", dbDSNMaster)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Error().Msgf("Could not connect to docker: %s", err)
	}

	conf = &Config{
		Master: &Conf{
			URL:   dbDSNMaster,
			Debug: false,
		},
		Slaves: []*Conf{},
	}

	// Do migration here
	c, err := NewConnection(conf)
	defer c.Close()
	if err != nil {
		log.Error().Msgf("error creating connection %s", err.Error())
		return
	}

	var createTable = `
		DROP TABLE IF EXISTS users;
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL NOT NULL,
			username VARCHAR NOT NULL,
			password VARCHAR NOT NULL
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique_idx ON users(username);	

		TRUNCATE users RESTART IDENTITY;
	`
	err = c.Writer().Exec(context.Background(), createTable)
	if err != nil {
		log.Error().Msgf("error init table %s", err.Error())
		return
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := dockerPool.Purge(postgreMaster); err != nil {
		log.Error().Msgf("Could not purge postgreMaster: %s", err)
	}

	os.Exit(code)
}

func TestNewConnection(t *testing.T) {
	convey.Convey("Test New Connection", t, func() {
		convey.Convey("All database is okay without error", func() {
			c, err := NewConnection(conf)
			defer c.Close()
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Master is error", func() {
			// don't need to close the connection since it fail to open
			_, err := NewConnection(&Config{
				Master: &Conf{
					URL: "postgres://user:password@not.exist.url",
				},
			})
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("One slave is error", func() {
			// don't need to close the connection since it fail to open
			c, err := NewConnection(&Config{
				Master: &Conf{
					URL: conf.Master.URL,
				},
				Slaves: []*Conf{
					{
						URL: "postgres://user:password@not.exist.url",
					},
				},
			})

			// we expect nil since we still can use this connection without any slave
			convey.So(err, convey.ShouldBeNil)

			err = c.Close()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestGoPgSql_Close(t *testing.T) {
	convey.Convey("Test close connection", t, func() {
		convey.Convey("Close all connection without error", func() {
			c, err := NewConnection(conf)
			convey.So(err, convey.ShouldBeNil)

			err = c.Close()
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Fail to close if master is error", func() {
			// don't need to close the connection since it fail to open
			c, err := NewConnection(&Config{
				Master: &Conf{
					URL: "postgres://user:password@not.exist.url",
				},
			})
			convey.So(err, convey.ShouldNotBeNil)

			err = c.Close()
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestGoPgSql_Writer(t *testing.T) {
	convey.Convey("Test Writer/Master connection only", t, func() {
		c, err := NewConnection(conf)
		convey.So(err, convey.ShouldBeNil)
		defer c.Close()

		convey.Convey("Insert using writer connection must be success", func() {
			err = c.Writer().Exec(context.Background(), insertQuery, myUser.Username, myUser.Password)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Get using writer connection", func() {
			// get user after insert
			u := &user{}
			err = c.Writer().Query(context.Background(), u, getQuery, myUser.Username)
			convey.So(err, convey.ShouldBeNil)
			convey.So(u.Username, convey.ShouldEqual, myUser.Username)
			convey.So(u.Password, convey.ShouldEqual, myUser.Password)
		})

		convey.Convey("Insert the same username should be error", func() {
			err = c.Writer().Exec(context.Background(), insertQuery, myUser.Username, myUser.Password)
			convey.So(err, convey.ShouldNotBeNil)
		})
	})

}

func TestGoPgSql_Reader(t *testing.T) {
	convey.Convey("Test Reader/Slave connection only", t, func() {
		c, err := NewConnection(conf)
		convey.So(err, convey.ShouldBeNil)
		defer c.Close()

		convey.Convey("Insert using reader connection must be error", func() {
			err = c.Reader().Exec(context.Background(),
				upsertQuery, myUser.Username, myUser.Password)
			// TODO: can't test since ory/dockertest not support linked docker image (need to link slave to master postgre)
			// convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("Get using reader connection", func() {
			// get user after insert
			u := &user{}
			err = c.Reader().Query(context.Background(), u, getQuery, myUser.Username)
			convey.So(err, convey.ShouldBeNil)
			convey.So(u.Username, convey.ShouldEqual, myUser.Username)
			convey.So(u.Password, convey.ShouldEqual, myUser.Password)
		})
	})

	convey.Convey("Test Write to master and read to slave", t, func() {
		c, err := NewConnection(conf)
		convey.So(err, convey.ShouldBeNil)
		defer c.Close()

		convey.Convey("Insert using writer connection", func() {
			err = c.Writer().Exec(context.Background(),
				insertQuery, myUser2.Username, myUser2.Password)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Fetch using reader connection", func() {
			// get user after insert
			u := &user{}
			err = c.Reader().Query(context.Background(), u, getQuery, myUser2.Username)
			convey.So(err, convey.ShouldBeNil)
			convey.So(u.Username, convey.ShouldEqual, myUser2.Username)
			convey.So(u.Password, convey.ShouldEqual, myUser2.Password)
		})
	})
}

func TestGoPgSql_NewTransaction(t *testing.T) {
	convey.Convey("Test transaction", t, func() {
		convey.Convey("Test insert and read using transaction", func() {
			c, err := NewConnection(conf)
			convey.So(err, convey.ShouldBeNil)
			defer c.Close()

			tx, err := c.NewTransaction(context.Background())
			convey.So(err, convey.ShouldBeNil)

			err = tx.Exec(context.Background(), insertQuery, myUser3.Username, myUser3.Password)
			convey.So(err, convey.ShouldBeNil)

			u := &user{}
			err = tx.Query(context.Background(), u, getQuery, myUser3.Username)
			convey.So(err, convey.ShouldBeNil)
			convey.So(u.Username, convey.ShouldEqual, myUser3.Username)
			convey.So(u.Password, convey.ShouldEqual, myUser3.Password)

			err = tx.Commit(context.Background())
			convey.So(err, convey.ShouldBeNil)

			err = tx.Rollback(context.Background())
			convey.So(err, convey.ShouldNotBeNil) // because we already commit, so cannot be rollback
		})

		convey.Convey("Test begin transaction if there are active transaction", func() {

			c, err := NewConnection(conf)
			convey.So(err, convey.ShouldBeNil)
			defer c.Close()

			tx1, err := c.NewTransaction(context.Background())
			convey.So(err, convey.ShouldBeNil)

			// since each NewTransaction uses its own session, so we can expect it can be called more than one
			tx2, err := c.NewTransaction(context.Background())
			convey.So(err, convey.ShouldBeNil)

			err = tx1.Commit(context.Background())
			convey.So(err, convey.ShouldBeNil)

			err = tx2.Commit(context.Background())
			convey.So(err, convey.ShouldBeNil)

			// but be careful, you know, we cannot begin unlimited transaction :P
		})
	})
}
