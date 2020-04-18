package sql

import (
	"database/sql"
	"time"
)

// Config configuration of provider
type ProviderConfig struct {
	Driver string

	DSN string

	// mysql max free idle
	MaxIdleConn int

	// mysql max open idle
	MaxOpenConn int

	// mysql conn max open idle
	ConnMaxLifetime time.Duration

	SQLGet        string
	SQLSave       string
	SQLRegenerate string
	SQLDestroy    string
	SQLCount      string
	SQLInsert     string
	SQLGC         string
}

// Provider backend manager
type Provider struct {
	config ProviderConfig

	*sql.DB
}
