package db

import "time"

// Conf represents a database connection configuration
type Conf struct {
	Debug        bool
	URL          string
	PoolSize     int
	IdleTimeout  int
	ConnLifetime time.Duration
}

// Config represents a configuration for this package
type Config struct {
	Master *Conf
	Slaves []*Conf
}
