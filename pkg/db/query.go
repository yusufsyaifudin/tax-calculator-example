package db

import (
	"context"
)

// SQL can be use writer or read replica only.
type SQL interface {
	Close() error // Will close all open database connection.
	Writer() SQLExecutor
	Reader() SQLExecutor
	NewTransaction(ctx context.Context) (Transaction, error)
}

// SQLExecutor should implements query and exec.
// You must select whether it use master or slave database.
type SQLExecutor interface {
	Query(ctx context.Context, destination interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) error
}

// Transaction must always using writer/master database.
type Transaction interface {
	Query(ctx context.Context, destination interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
