package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"net/http/httptest"

	"net/url"

	"github.com/ory/dockertest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartystreets/goconvey/convey"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/app/restapi"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/conn"
	"github.com/yusufsyaifudin/tax-calculator-example/pkg/db"
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

var (
	dbDSNMaster       string
	conf              *db.Config
	apiV1RegisterURL  string
	apiV1LoginURL     string
	apiV1CreateTaxURL string
	apiV1GetTaxURL    string
)

var refreshDB = func(dbURL string) {
	conn.MigrateReset(dbURL)
	conn.MigrateSync(dbURL)
}

func TestMain(m *testing.M) {
	logger.Level(zerolog.Disabled)

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

	dbDSNMaster = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		PostgresUser, PostgresPassword, postgreMaster.GetBoundIP(PostgreImage), postgreMaster.GetPort("5432/tcp"), PostgresDB)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := dockerPool.Retry(func() error {
		var err error
		c, err := sql.Open("postgres", dbDSNMaster)
		if err != nil {
			return err
		}
		return c.Ping()
	}); err != nil {
		log.Error().Msgf("Could not connect to docker: %s", err)
	}

	conf = &db.Config{
		Master: &db.Conf{
			URL:   dbDSNMaster,
			Debug: false,
		},
		Slaves: []*db.Conf{},
	}

	// Do migration here
	c, err := db.NewConnection(conf)
	defer c.Close()
	if err != nil {
		log.Error().Msgf("error creating connection %s", err.Error())
		return
	}

	// set connection
	conn.SetDBConnection(c)

	// migrate sync
	if err := conn.MigrateSync(dbDSNMaster); err != nil {
		log.Error().Err(err).Msg("fail to sync db")
		return
	}

	serverConfig := &restapi.Config{
		Address: *serverAddr,
		Test:    true,
	}

	// now, start the server
	restapi.Configure(serverConfig)

	s := httptest.NewServer(restapi.Router)
	apiV1BaseURL := fmt.Sprintf("%s/api/v1", s.URL)

	// list routes
	apiV1RegisterURL = fmt.Sprintf("%s/register", apiV1BaseURL)
	apiV1LoginURL = fmt.Sprintf("%s/login", apiV1BaseURL)
	apiV1CreateTaxURL = fmt.Sprintf("%s/tax", apiV1BaseURL)
	apiV1GetTaxURL = fmt.Sprintf("%s/tax", apiV1BaseURL)

	code := m.Run()
	s.Close() // shutdown the server after done

	// You can't defer this because os.Exit doesn't care for defer
	if err := dockerPool.Purge(postgreMaster); err != nil {
		log.Error().Msgf("Could not purge postgreMaster: %s", err)
	}

	os.Exit(code)
}

func TestRegisterEndpoint(t *testing.T) {
	convey.Convey("Test Register Endpoint", t, func() {
		convey.Convey("Username should be exist", func() {

			form := &url.Values{}
			form.Set("username", "")
			form.Set("password", "exist")
			res, status, err := httpPost(apiV1RegisterURL, "", form)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldNotEqual, 200)

			convey.So(res, convey.ShouldResemble, map[string]interface{}{
				"error_code":       "0_0001",
				"http_status_code": float64(400),
				"message":          "username: required",
			})
		})

		convey.Convey("Password should be exist", func() {

			form := &url.Values{}
			form.Set("username", "exist")
			form.Set("password", "")
			res, status, err := httpPost(apiV1RegisterURL, "", form)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldNotEqual, 200)

			convey.So(res, convey.ShouldResemble, map[string]interface{}{
				"error_code":       "0_0001",
				"http_status_code": float64(400),
				"message":          "password: required",
			})
		})

		refreshDB(dbDSNMaster)
		convey.Convey("Should success register", func() {
			const (
				username = "john_doe"
				password = "password"
			)

			form := &url.Values{}
			form.Set("username", username)
			form.Set("password", password)
			res, status, err := httpPost(apiV1RegisterURL, "", form)

			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldEqual, 200)

			resUsername := res["user"].(map[string]interface{})["username"].(string)
			convey.So(resUsername, convey.ShouldEqual, username)
		})
	})
}

func TestLoginEndpoint(t *testing.T) {
	convey.Convey("Test Login Endpoint", t, func() {
		convey.Convey("Username should be exist", func() {

			form := &url.Values{}
			form.Set("username", "")
			form.Set("password", "exist")
			res, status, err := httpPost(apiV1LoginURL, "", form)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldNotEqual, 200)

			convey.So(res, convey.ShouldResemble, map[string]interface{}{
				"error_code":       "0_0001",
				"http_status_code": float64(400),
				"message":          "username: required",
			})
		})

		convey.Convey("Password should be exist", func() {

			form := &url.Values{}
			form.Set("username", "exist")
			form.Set("password", "")
			res, status, err := httpPost(apiV1LoginURL, "", form)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldNotEqual, 200)

			convey.So(res, convey.ShouldResemble, map[string]interface{}{
				"error_code":       "0_0001",
				"http_status_code": float64(400),
				"message":          "password: required",
			})
		})

		refreshDB(dbDSNMaster)
		convey.Convey("User can login after register", func() {
			const (
				username = "john_doe"
				password = "password"
			)

			formRegister := &url.Values{}
			formRegister.Set("username", username)
			formRegister.Set("password", password)
			res, status, err := httpPost(apiV1RegisterURL, "", formRegister)

			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldEqual, 200)

			resUsername := res["user"].(map[string]interface{})["username"].(string)
			convey.So(resUsername, convey.ShouldEqual, username)

			formLogin := &url.Values{}
			formLogin.Set("username", username)
			formLogin.Set("password", password)
			res, status, err = httpPost(apiV1LoginURL, "", formLogin)
			convey.So(err, convey.ShouldBeNil)
			convey.So(status, convey.ShouldEqual, 200)

			loginUsername := res["user"].(map[string]interface{})["username"].(string)
			convey.So(loginUsername, convey.ShouldEqual, username)
		})
	})
}
