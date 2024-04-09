package db

import (
	"context"
	"fmt"
	"github.com/Pactus-Contrib/Indexer/logging"
	"github.com/Pactus-Contrib/Indexer/schema"
	"golang.org/x/sync/errgroup"
)

type Database interface {
	Executor
	Name() string
	Type() string
	Engine() string
	Close() error
	Migrate(ctx context.Context, indexerUUid string, lastBlockHeight int) error
}

type Executor interface {
	FindOne(ctx context.Context, tableOrCollectionName string, key string, val any, resultPtr any) error
	InsertOne(ctx context.Context, tableOrCollectionName string, dataPtr any) error
	InsertMany(ctx context.Context, tableOrCollectionName string, dataPtr []any) error
	UpdateOne(ctx context.Context, tableOrCollectionName string, key string, val any, updKey string, updVal any) error
}

type Pool struct {
	items   []Database
	logging logging.Logger
}

func NewPool(logging logging.Logger) *Pool {
	return &Pool{
		items:   make([]Database, 0),
		logging: logging,
	}
}

func (p *Pool) RegisterEngine(db Database) {
	p.items = append(p.items, db)
}

func (p *Pool) Close() error {
	for _, item := range p.items {
		if err := item.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pool) Migration(ctx context.Context, indexerUUid string, lastBlockHeight int) error {
	gp, gpCtx := errgroup.WithContext(ctx)

	for _, item := range p.items {
		gp.Go(func() error {
			p.logging.Info(false, fmt.Sprintf("Start migrate database %s", item.Name()))
			return item.Migrate(gpCtx, indexerUUid, lastBlockHeight)
		})
	}

	return gp.Wait()
}

func (p *Pool) GetIndexer(ctx context.Context, indexerId string) ([]schema.Indexer, error) {
	indexers := make([]schema.Indexer, 0)

	gp, gpCtx := errgroup.WithContext(ctx)

	for _, item := range p.items {
		gp.Go(func() error {
			var indexer schema.Indexer

			if err := item.FindOne(gpCtx, schema.IndexerTableName, "index_id", indexerId, &indexer); err != nil {
				return newErr(item.Name(), item.Engine(), item.Type(), err.Error())
			}

			indexers = append(indexers, indexer)

			return nil
		})
	}

	return indexers, gp.Wait()
}

func (p *Pool) InsertOne(ctx context.Context, dataPtr any) error {
	gp, gpCtx := errgroup.WithContext(ctx)

	for _, item := range p.items {
		gp.Go(func() error {
			if err := item.InsertOne(gpCtx, schema.IndexerTableName, dataPtr); err != nil {
				return newErr(item.Name(), item.Engine(), item.Type(), err.Error())
			}

			return nil
		})
	}

	return gp.Wait()
}

func (p *Pool) InsertMany(ctx context.Context, dataPtr []any) error {
	gp, gpCtx := errgroup.WithContext(ctx)

	for _, item := range p.items {
		gp.Go(func() error {
			if err := item.InsertMany(gpCtx, schema.IndexerTableName, dataPtr); err != nil {
				return newErr(item.Name(), item.Engine(), item.Type(), err.Error())
			}

			return nil
		})
	}

	return gp.Wait()
}
