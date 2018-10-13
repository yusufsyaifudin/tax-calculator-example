package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/namsral/flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/app/restapi"
	"github.com/yusufsyaifudin/tax-calculator-example/internal/pkg/conn"
	"github.com/yusufsyaifudin/tax-calculator-example/pkg/db"
)

var (
	serverAddr = flag.String("address", "localhost:9000", "Address and port to bind")
	debug      = flag.Bool("debug", false, "Print log")

	dbSyncMigration = flag.Bool("db-sync-migration", false, "To sync database structure")
	dbUrlMaster     = flag.String(
		"db-master-url",
		"postgres://postgres:postgres@localhost/tax-calculator?sslmode=disable",
		"PostgreSQL master server DSN",
	)

	dbUrlSlave = flag.String(
		"db-slaves-url",
		"",
		"PostgreSQL slaves server DSN in semicolon separated value",
	)
)

var logger = log.With().Str("pkg", "main").Logger()

// Start the application.
func main() {
	flag.Parse()

	// show all log if debug mode activated
	if *debug {
		log.WithLevel(zerolog.DebugLevel)
	}

	dbSlaveUrls := strings.Split(*dbUrlSlave, ";")
	dbConfigSlave := make([]*db.Conf, len(dbSlaveUrls))
	for i, dbSlaveUrl := range dbSlaveUrls {
		dbConfigSlave[i] = &db.Conf{
			URL:   strings.TrimSpace(dbSlaveUrl),
			Debug: *debug,
		}
	}

	dbConf := &db.Config{
		Master: &db.Conf{
			URL:   *dbUrlMaster,
			Debug: *debug,
		},
	}

	dbConn, err := db.NewConnection(dbConf)
	defer dbConn.Close()
	if err != nil {
		logger.Error().Err(err).Msg("fail creating db connection")
		panic(err)
	}

	// set connection to global, so that it can be accessed from any package inside internal/app
	conn.SetDBConnection(dbConn)

	if *dbSyncMigration {
		logger.Info().Msg("Syncing database migration...")
		if err := conn.MigrateSync(*dbUrlMaster); err != nil {
			logger.Error().Err(err).Msg("fail do migration to database")
			return
		}
	}

	serverConfig := &restapi.Config{
		Address: *serverAddr,
	}

	restapi.Configure(serverConfig)

	var apiErrChan = make(chan error, 1)
	go func() {
		logger.Info().Msgf("running server on %s", serverConfig.Address)
		apiErrChan <- restapi.Run()
	}()

	// gracefully shutdown the server
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.Info().Msg("got an interrupt, exiting...")
		restapi.Shutdown()
	case err := <-apiErrChan:
		if err != nil {
			logger.Error().Err(err).Msg("error while running api, exiting...")
		}
	}

}
