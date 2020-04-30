package postgre

import (
	"fmt"

	"github.com/fasthttp/session/v2/internal/sql"

	// Import postgres driver
	_ "github.com/lib/pq"
)

var initQueries = []string{
	"DROP TABLE IF EXISTS %s;",
	`CREATE TABLE IF NOT EXISTS %s (
		id VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		data TEXT NOT NULL,
		last_active BIGINT NOT NULL DEFAULT '0',
		expiration BIGINT NOT NULL DEFAULT '0'
	);`,
	"CREATE INDEX last_active ON %s (last_active);",
	"CREATE INDEX expiration ON %s (expiration);",
}

// New returns a new configured postgres provider
func New(cfg Config) (*Provider, error) {
	if cfg.Host == "" {
		return nil, errConfigHostEmpty
	}
	if cfg.Port == 0 {
		return nil, errConfigPortZero
	}

	providerCfg := sql.ProviderConfig{
		Driver:          "postgres",
		DSN:             cfg.dsn(),
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxOpenConns:    cfg.MaxOpenConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		SQLGet:          fmt.Sprintf("SELECT data FROM %s WHERE id=$1", cfg.TableName),
		SQLSave:         fmt.Sprintf("UPDATE %s SET data=$1,last_active=$2,expiration=$3 WHERE id=$4", cfg.TableName),
		SQLRegenerate:   fmt.Sprintf("UPDATE %s SET id=$1,last_active=$2,expiration=$3 WHERE id=$4", cfg.TableName),
		SQLCount:        fmt.Sprintf("SELECT count(id) as total FROM %s", cfg.TableName),
		SQLDestroy:      fmt.Sprintf("DELETE FROM %s WHERE id=$1", cfg.TableName),
		SQLInsert:       fmt.Sprintf("INSERT INTO %s (id, data, last_active, expiration) VALUES ($1,$2,$3,$4)", cfg.TableName),
		SQLGC:           fmt.Sprintf("DELETE FROM %s WHERE last_active+expiration<=$1 AND expiration<>0", cfg.TableName),
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
	for _, query := range initQueries {
		_, err := p.Exec(fmt.Sprintf(query, p.config.TableName))
		if err != nil {
			p.Close()

			return err
		}
	}

	return nil
}
