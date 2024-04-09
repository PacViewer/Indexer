package db

import (
	"context"
	"errors"
	"github.com/Pactus-Contrib/Indexer/schema"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type SQL struct {
	db     *gorm.DB
	engine schema.DatabaseEngine
	dbType schema.DatabaseType
	uri    string
	name   string
}

func NewSQL(dbCfg *schema.DB) (Database, error) {
	sql := new(SQL)

	switch dbCfg.Engine {
	case schema.MYSQL, schema.MARIADB:
		db, err := gorm.Open(mysql.Open(dbCfg.URI), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			return nil, err
		}
		sql.db = db
	case schema.POSTGRESQL:
		db, err := gorm.Open(postgres.Open(dbCfg.URI), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			return nil, err
		}
		sql.db = db
	default:
		return nil, errors.New("database engine is invalid")
	}

	sql.uri = dbCfg.URI
	sql.engine = dbCfg.Engine
	sql.dbType = dbCfg.Type
	sql.name = dbCfg.Name

	return sql, nil
}

func (s *SQL) Close() error {
	sdb, err := s.db.DB()
	if err != nil {
		return err
	}

	return sdb.Close()
}

func (s *SQL) Migrate(ctx context.Context, indexerUUid string, lastBlockHeight int) error {
	mg := s.db.WithContext(ctx)
	err := mg.Migrator().AutoMigrate(&schema.Block{}, &schema.Transaction{}, &schema.Indexer{})
	if err != nil {
		return err
	}

	mg.Create(&schema.Indexer{
		IndexId:         indexerUUid,
		LastBlockHeight: lastBlockHeight,
		IndexedAt:       time.Now(),
	})

	return mg.Error
}

func (s *SQL) Name() string {
	return s.name
}

func (s *SQL) Type() string {
	return s.dbType.String()
}

func (s *SQL) Engine() string {
	return s.engine.String()
}

func (s *SQL) FindOne(ctx context.Context, tableOrCollectionName string, key string, val any, resultPtr any) error {
	tb := s.db.Table(tableOrCollectionName)
	tb.WithContext(ctx)
	tb.Where(map[string]interface{}{key: val}).First(resultPtr)

	return tb.Error
}

func (s *SQL) InsertOne(ctx context.Context, tableOrCollectionName string, dataPtr any) error {
	tb := s.db.Table(tableOrCollectionName)
	tb.WithContext(ctx)

	tb.Create(dataPtr)

	return tb.Error
}

func (s *SQL) InsertMany(ctx context.Context, tableOrCollectionName string, dataPtr []any) error {
	tb := s.db.Table(tableOrCollectionName)
	tb.WithContext(ctx)

	tb.CreateInBatches(dataPtr, len(dataPtr))

	return tb.Error
}

func (s *SQL) UpdateOne(ctx context.Context, tableOrCollectionName string, key string, val any,
	updKey string, updVal any) error {
	tb := s.db.Table(tableOrCollectionName)
	tb.WithContext(ctx)

	tb.Where(map[string]interface{}{key: val}).Update(updKey, updVal)

	return tb.Error
}
