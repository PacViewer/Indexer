package schema

import (
	"time"
)

const (
	BlockTableName        = "blocks"
	TransactionsTableName = "transactions"
	IndexerTableName      = "indexers"
)

type Block struct {
	ID                uint    `bson:"-" gorm:"primarykey"`
	Height            uint32  `bson:"height" gorm:"column:height;uniqueIndex"`
	Hash              string  `bson:"hash" gorm:"column:hash;uniqueIndex;size:100"`
	TotalTransactions uint    `bson:"total_transactions" gorm:"column:total_transactions"`
	BlockTime         uint32  `bson:"block_time" gorm:"column:block_time;index"`
	BlockReward       int64   `bson:"block_reward" gorm:"column:block_reward"`
	Version           int32   `bson:"version" gorm:"column:version"`
	PrevBlockHash     string  `bson:"prev_block_hash,omitempty" gorm:"column:prev_block_hash;index"`
	StateRoot         string  `bson:"state_root,omitempty" gorm:"column:state_root"`
	SortitionSeed     string  `bson:"sortition_seed,omitempty" gorm:"column:sortition_seed"`
	ProposerAddress   string  `bson:"proposer_address" gorm:"column:proposer_address;index"`
	CertificateHash   string  `bson:"certificate_hash,omitempty" gorm:"column:certificate_hash;index"`
	Round             int32   `bson:"round,omitempty" gorm:"column:round"`
	Committers        []int32 `bson:"committers,omitempty" gorm:"serializer:json"`
	Absentees         []int32 `bson:"absentees,omitempty" gorm:"serializer:json"`
	Signature         string  `bson:"signature,omitempty" gorm:"column:signature"`
}

type Transaction struct {
	ID          uint      `bson:"-" gorm:"primarykey"`
	Hash        string    `bson:"hash" gorm:"column:hash;uniqueIndex;size:100"`
	BlockHeight uint32    `bson:"block_height" gorm:"column:block_height;index"`
	Version     int32     `bson:"version,omitempty" gorm:"column:block_height"`
	Type        string    `bson:"type" gorm:"column:type;index"`
	From        string    `bson:"from,omitempty" gorm:"column:from;index"`
	To          string    `bson:"to,omitempty" gorm:"column:to;index"`
	Value       int64     `bson:"value,omitempty" gorm:"column:value;index"`
	Fee         int64     `bson:"fee,omitempty" gorm:"column:fee;index"`
	Memo        string    `bson:"memo,omitempty" gorm:"column:memo"`
	CreatedAt   time.Time `bson:"created_at" gorm:"column:created_at;index"`
}

type Indexer struct {
	ID              uint      `bson:"-" gorm:"primarykey"`
	IndexId         string    `bson:"index_id" gorm:"column:index_id;uniqueIndex;size:36"`
	LastBlockHeight int       `bson:"last_block_height" gorm:"column:last_block_height"`
	IndexedAt       time.Time `bson:"indexed_at" gorm:"column:indexed_at"`
}
