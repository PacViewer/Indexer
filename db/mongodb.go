package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Pactus-Contrib/Indexer/schema"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongodb struct {
	cli    *mongo.Client
	db     *mongo.Database
	engine schema.DatabaseEngine
	dbType schema.DatabaseType
	uri    string
	name   string
}

func NewMongodb(ctx context.Context, db *schema.DB) (Database, error) {
	mongodb := new(Mongodb)

	opts := make([]*options.ClientOptions, 0)
	opts = append(opts, options.Client().ApplyURI(db.URI))
	opts = append(opts, options.Client().SetConnectTimeout(2*time.Second))
	opts = append(opts, options.Client().SetMaxConnecting(500))
	cli, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if err := cli.Ping(ctx, nil); err != nil {
		return nil, errors.Join(errors.New("mongodb: can't verify client connection "), err)
	}

	mongodb.cli = cli
	mongodb.db = cli.Database(db.Database)
	mongodb.uri = db.URI
	mongodb.engine = db.Engine
	mongodb.dbType = db.Type
	mongodb.name = db.Name

	return mongodb, nil
}

func (m *Mongodb) Close() error {
	return m.cli.Disconnect(context.Background())
}

func (m *Mongodb) Migrate(ctx context.Context, indexerUUid string, lastBlockHeight int) error {
	migration := migrate.NewMigrate(m.db, migrate.Migration{
		Version: 1,
		Up: func(db *mongo.Database) error {
			block := db.Collection(schema.BlockTableName)
			transaction := db.Collection(schema.TransactionsTableName)
			indexer := db.Collection(schema.IndexerTableName)

			// block
			if err := addUniqueIndex(ctx, block, "height", false); err != nil {
				return err
			}

			if err := addUniqueIndex(ctx, block, "hash", false); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, block, "block_time"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, block, "proposer_address"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, block, "certificate_hash"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, block, "prev_block_hash"); err != nil {
				return err
			}

			// transaction
			if err := addUniqueIndex(ctx, transaction, "hash", false); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "block_height"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "type"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "from"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "to"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "value"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "fee"); err != nil {
				return err
			}

			if err := addNormalIndex(ctx, transaction, "created_at"); err != nil {
				return err
			}

			// indexer
			if err := addUniqueIndex(ctx, indexer, "index_id", false); err != nil {
				return err
			}

			if _, err := indexer.InsertOne(ctx, &schema.Indexer{
				IndexId:         indexerUUid,
				LastBlockHeight: lastBlockHeight,
				IndexedAt:       time.Now(),
			}); err != nil {
				return err
			}

			return nil
		},
	})

	return migration.Up(1)
}

func (m *Mongodb) Name() string {
	return m.name
}

func (m *Mongodb) Type() string {
	return m.dbType.String()
}

func (m *Mongodb) Engine() string {
	return m.engine.String()
}

func (m *Mongodb) FindOne(ctx context.Context, tableOrCollectionName string, key string, val any, resultPtr any) error {
	col := m.db.Collection(tableOrCollectionName)

	return col.FindOne(ctx, bson.M{key: val}).Decode(resultPtr)
}

func (m *Mongodb) InsertOne(ctx context.Context, tableOrCollectionName string, dataPtr any) error {
	col := m.db.Collection(tableOrCollectionName)

	_, err := col.InsertOne(ctx, dataPtr)
	return err
}

func (m *Mongodb) InsertMany(ctx context.Context, tableOrCollectionName string, dataPtr []any) error {
	col := m.db.Collection(tableOrCollectionName)

	_, err := col.InsertMany(ctx, dataPtr)
	return err
}

func (m *Mongodb) UpdateOne(ctx context.Context, tableOrCollectionName string, key string, val any,
	updKey string, updVal any) error {
	col := m.db.Collection(tableOrCollectionName)

	opts := options.Update().SetUpsert(false)

	_, err := col.UpdateOne(ctx, bson.M{key: val}, bson.M{"$set": bson.M{updKey: updVal}}, opts)
	if err != nil {
		return err
	}

	return nil
}

// addNormalIndex create normal index for migration
func addNormalIndex(ctx context.Context, collection *mongo.Collection, field string) error {
	opt := options.Index().SetName(fmt.Sprintf("%s_%s_normal", collection.Name(), field)).SetUnique(false)
	keys := bson.D{{field, 1}}
	model := mongo.IndexModel{Keys: keys, Options: opt}
	_, err := collection.Indexes().CreateOne(ctx, model)
	return err
}

// addUniqueIndex create unique index, if sparse true the index will only reference documents that contain the fields specified in the index
func addUniqueIndex(ctx context.Context, collection *mongo.Collection, field string, sparse bool) error {
	opt := options.Index().SetName(fmt.Sprintf("%s_%s_unique", collection.Name(), field)).SetUnique(true)
	if sparse {
		opt = opt.SetSparse(true)
	}
	keys := bson.D{{field, 1}}
	model := mongo.IndexModel{Keys: keys, Options: opt}
	_, err := collection.Indexes().CreateOne(ctx, model)
	return err
}
