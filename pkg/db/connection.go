package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Caller().Str("pkg", "db").Logger()

// NewConnection will create new connection for selected package.
func NewConnection(config *Config) (sql SQL, err error) {
	masterConnection, err := createConnection(config.Master)
	if err != nil {
		return &goPgSQL{}, err
	}

	var slaveConnections = make([]*pg.DB, len(config.Slaves))
	for i, conf := range config.Slaves {
		slaveConn, err := createConnection(conf)
		if err != nil {
			if conf.Debug {
				logger.Error().Err(err).Msgf("error when creating slave connection on: %s", conf.URL)
			}

			continue
		}

		slaveConnections[i] = slaveConn
	}

	return newGoPgConnection(masterConnection, slaveConnections)
}

// createConnection create connection using go-pg
func createConnection(conf *Conf) (*pg.DB, error) {
	opt, err := pg.ParseURL(conf.URL)
	if err != nil {
		return nil, err
	}

	opt.PoolSize = conf.PoolSize
	opt.IdleTimeout = time.Duration(conf.IdleTimeout) * time.Second
	opt.MaxConnAge = time.Duration(conf.ConnLifetime) * time.Second
	db := pg.Connect(opt)

	if conf.Debug {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				log.Printf("error when log query, %s", err.Error())
				return
			}

			elapsedTime := float64(time.Since(event.StartTime).Nanoseconds()) / float64(1000000)
			logger.Debug().
				Str("elapsedTime", fmt.Sprintf("%0.2f ms", elapsedTime)).
				Str("query", query).
				Msg("")
		})
	}

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		logger.Debug().Msgf("connected to database %s", opt.Addr)
	}

	return db, nil
}

// newGoPgConnection will create database connection using go-pg
func newGoPgConnection(master *pg.DB, slaves []*pg.DB) (sql SQL, err error) {
	return &goPgSQL{
		master: master,
		slaves: slaves,
	}, nil
}

type goPgSQL struct {
	sync.RWMutex

	master *pg.DB
	slaves []*pg.DB
}

func (g *goPgSQL) Close() error {
	var err, lastError error

	if g.master == nil {
		return fmt.Errorf("master connection is not exist, hence cannot be closed")
	}

	lastError = g.master.Close()
	if lastError != nil {
		err = lastError
	}

	for _, slave := range g.slaves {
		if slave == nil {
			continue
		}

		lastError = slave.Close()
		if lastError != nil {
			err = lastError
		}
	}

	return err
}

// Writer always use the master.
func (g *goPgSQL) Writer() SQLExecutor {
	return &goPgSQLWriter{
		master: g.master,
	}
}

// Reader will select the read replica based on highest idle connection.
// This guarantee that return value will always use the same host.
// If slave is empty, it will use master connection as the default db connection.
// As a result, you must ensure that master never shutdown.
func (g *goPgSQL) Reader() SQLExecutor {
	g.RLock() // lock the reader when election occurred

	var newSlaves = make([]*pg.DB, len(g.slaves))

	conn := g.master
	currentIdleConnection := uint32(0)
	for i, slave := range g.slaves {
		if slave == nil {
			continue
		}

		// also check and take out the unresponsive slave
		_, err := slave.Exec("SELECT 1")
		if err != nil {
			logger.Error().Err(err).Msgf("error occurred on db: %s", slave.Options().Addr)
			continue
		}

		newSlaves[i] = slave

		idleConn := slave.PoolStats().Hits
		if idleConn >= currentIdleConnection {
			conn = slave
		}
	}

	// assign new active slave connection list
	g.slaves = newSlaves
	g.RUnlock() // and unlock the RW

	return &goPgSQLReader{
		slave: conn,
	}
}

// NewTransaction will always use the master node.
func (g *goPgSQL) NewTransaction(ctx context.Context) (Transaction, error) {
	tx, err := g.master.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{
		conn: tx,
	}, nil
}

// ====================== WRITER
type goPgSQLWriter struct {
	master *pg.DB
}

func (w *goPgSQLWriter) Query(ctx context.Context, out interface{}, query string, args ...interface{}) error {
	_, err := w.master.Query(out, query, args...)
	return err
}

func (w *goPgSQLWriter) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := w.master.Exec(query, args...)
	return err
}

// ====================== READER
type goPgSQLReader struct {
	slave *pg.DB
}

func (r *goPgSQLReader) Query(ctx context.Context, out interface{}, query string, args ...interface{}) error {
	_, err := r.slave.Query(out, query, args...)
	return err
}

func (r *goPgSQLReader) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := r.slave.Exec(query, args...)
	return err
}

// ====================== TRANSACTION
type transaction struct {
	conn *pg.Tx
}

func (t *transaction) Query(ctx context.Context, out interface{}, query string, args ...interface{}) error {
	_, err := t.conn.Query(out, query, args...)
	return err
}

func (t *transaction) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := t.conn.Exec(query, args...)
	return err
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.conn.Commit()
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.conn.Rollback()
}
