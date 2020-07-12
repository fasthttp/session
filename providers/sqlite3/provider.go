package sqlite3

import (
	"fmt"

	"github.com/fasthttp/session/v2/internal/sql"

	// Import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	dropQuery   = "DROP TABLE IF EXISTS %s;"
	initQueries = []string{
		`CREATE TABLE IF NOT EXISTS %s (
		id VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		data TEXT NOT NULL,
		last_active BIGINT NOT NULL DEFAULT '0',
		expiration BIGINT NOT NULL DEFAULT '0'
	);`,
		"CREATE INDEX last_active ON %s (last_active);",
		"CREATE INDEX expiration ON %s (expiration);",
	}
)

// New returns a new configured sqlite3 provider
func New(cfg Config) (*Provider, error) {
	if cfg.DBPath == "" {
		return nil, errConfigDBPathEmpty
	}

	providerCfg := sql.ProviderConfig{
		Driver:          "sqlite3",
		DSN:             cfg.DBPath,
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxOpenConns:    cfg.MaxOpenConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		SQLGet:          fmt.Sprintf("SELECT data FROM %s WHERE id=?", cfg.TableName),
		SQLSave:         fmt.Sprintf("UPDATE %s SET data=?,last_active=?,expiration=? WHERE id=?", cfg.TableName),
		SQLRegenerate:   fmt.Sprintf("UPDATE %s SET id=?,last_active=?,expiration=? WHERE id=?", cfg.TableName),
		SQLCount:        fmt.Sprintf("SELECT count(id) as total FROM %s", cfg.TableName),
		SQLDestroy:      fmt.Sprintf("DELETE FROM %s WHERE id=?", cfg.TableName),
		SQLInsert:       fmt.Sprintf("INSERT INTO %s (id, data, last_active, expiration) VALUES (?,?,?,?)", cfg.TableName),
		SQLGC:           fmt.Sprintf("DELETE FROM %s WHERE last_active+expiration<=? AND expiration<>0", cfg.TableName),
	}

	provider, err := sql.NewProvider(providerCfg)
	if err != nil {
		return nil, err
	}

	p := &Provider{
		config:   cfg,
		Provider: provider,
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Provider) init() error {
	if p.config.DropTable {
		_, err := p.Exec(fmt.Sprintf(dropQuery, p.config.TableName))
		if err != nil {
			p.Close()
			return err
		}
	}
	for _, query := range initQueries {
		_, err := p.Exec(fmt.Sprintf(query, p.config.TableName))
		if err != nil {
			p.Close()
			return err
		}
	}

	return nil
}
