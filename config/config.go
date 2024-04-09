package config

import (
	"github.com/Pactus-Contrib/Indexer/logging"
	"github.com/Pactus-Contrib/Indexer/schema"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
)

var DefaultConfig = &schema.Config{
	LastBlockHeight:      1,
	SyncIntervalPerBlock: 5,
	IndexerUuid:          uuid.New().String(),
	Pactus: &schema.Pactus{
		RPC: "localhost:50051",
	},
	DBS: []*schema.DB{
		{
			Name:   "mongodb",
			Type:   "nosql",
			Engine: "mongodb",
			URI:    "mongodb://admin:password123@localhost:27017/mydatabase",
		},
	},
	Logging: &schema.Logging{
		Debug:        false,
		Handler:      logging.ConsoleHandler,
		EnableCaller: true,
	},
}

func New(path string) (*schema.Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg schema.Config

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
