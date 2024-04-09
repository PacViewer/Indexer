package schema

import (
	"errors"
	"fmt"
	"github.com/Pactus-Contrib/Indexer/logging"
	"github.com/google/uuid"
)

type Config struct {
	LastBlockHeight      int      `yaml:"last_block_height"`
	SyncIntervalPerBlock int      `yaml:"sync_interval_per_block"`
	IndexerUuid          string   `yaml:"indexer_uuid"`
	Pactus               *Pactus  `yaml:"pactus"`
	DBS                  []*DB    `yaml:"dbs"`
	Logging              *Logging `yaml:"logging"`
}

type Pactus struct {
	RPC string `yaml:"rpc"`
}

type DB struct {
	Name     string         `yaml:"name"`
	Type     DatabaseType   `yaml:"type"`
	Engine   DatabaseEngine `yaml:"engine"`
	URI      string         `yaml:"uri"`
	Database string         `yaml:"database"`
}

type Logging struct {
	Debug        bool               `yaml:"debug" json:"debug"`
	Handler      logging.HandleType `yaml:"handler" json:"handler"` // Handler 0= console handler, 1= text handler, 2= json handler
	EnableCaller bool               `yaml:"enable_caller" json:"enable_caller"`
	SentryDSN    string             `yaml:"sentry_dsn" json:"sentry_dsn"`
}

type (
	DatabaseType   string
	DatabaseEngine string
)

const (
	SQL   DatabaseType = "sql"
	NOSQL DatabaseType = "nosql"
)

func (d DatabaseType) String() string {
	return string(d)
}

const (
	MYSQL      DatabaseEngine = "mysql"
	MARIADB    DatabaseEngine = "mariadb"
	POSTGRESQL DatabaseEngine = "psql"
	MONGODB    DatabaseEngine = "mongodb"
)

func (d DatabaseEngine) String() string {
	return string(d)
}

func (c *Config) Validate() error {
	if c.LastBlockHeight == 0 {
		return errors.New("you cannot set 0 for last_block_height, first block is 1")
	}

	if c.SyncIntervalPerBlock < 3 || c.SyncIntervalPerBlock > 86400 {
		return errors.New("minimum sync_interval_per_block is 3 second and max is 86400 or 24 hours")
	}

	_, err := uuid.Parse(c.IndexerUuid)
	if err != nil {
		return fmt.Errorf("indexer_uuid is invalid, please set this uuid %s", uuid.New().String())
	}

	if len(c.DBS) == 0 {
		return errors.New("dbs is null, need 1 database engine for sync")
	}

	if c.Pactus == nil {
		return errors.New("pactus config is null")
	}

	if len(c.Pactus.RPC) == 0 {
		return errors.New("pactus rpc address is empty")
	}

	for _, db := range c.DBS {
		if len(db.Name) == 0 {
			return errors.New("db name is null, please set a name for database engine")
		}

		isValidDBType := false

		switch db.Type {
		case SQL, NOSQL:
			isValidDBType = true
		}

		if !isValidDBType {
			return errors.New("db type is invalid, please set a type for database engine (sql, nosql)")
		}

		isValidDBEngine := false

		switch db.Engine {
		case MYSQL, POSTGRESQL, MARIADB:
			isValidDBEngine = true
		case MONGODB:
			isValidDBEngine = true
			if len(db.Database) == 0 {
				return errors.New("database name is required for mongodb")
			}
		}

		if !isValidDBEngine {
			return errors.New("db engine is invalid, please set a engine " +
				"for database engine (mysql, psql, mariadb, mongodb)")
		}

	}

	return nil
}
