package database

import (
	"errors"
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"
)

type Blocks struct {
	GUID       uuid.UUID `gorm:"primaryKey" json:"guid"`
	Hash       string    `gorm:"primaryKey"`
	ParentHash string    `json:"parent_hash"`
	Number     *big.Int  `gorm:"serializer:u256"`
	Timestamp  uint64
}

type BlocksView interface {
	LatestBlocks() (*Blocks, error)
}

type BlocksDB interface {
	BlocksView

	StoreBlockss([]Blocks, uint64) error
}

type blocksDB struct {
	gorm *gorm.DB
}

func NewBlocksDB(db *gorm.DB) BlocksDB {
	return &blocksDB{gorm: db}
}

func (db *blocksDB) StoreBlockss(headers []Blocks, blockLength uint64) error {
	result := db.gorm.CreateInBatches(&headers, int(blockLength))
	return result.Error
}

func (db *blocksDB) LatestBlocks() (*Blocks, error) {
	var l1Header Blocks
	result := db.gorm.Order("number DESC").Take(&l1Header)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &l1Header, nil
}
