package database

import (
	"errors"
	"gorm.io/gorm"
	"math/big"

	"github.com/google/uuid"
)

type Transactions struct {
	GUID         uuid.UUID `gorm:"primaryKey" json:"guid"`
	BlockHash    string    `gorm:"column:block_hash;serializer:bytes"  db:"block_hash" json:"block_hash"`
	BlockNumber  *big.Int  `gorm:"serializer:u256;column:block_number" db:"block_number" json:"BlockNumber" form:"block_number"`
	Hash         string    `json:"hash"`
	FromAddress  string    `json:"from_address"`
	ToAddress    string    `json:"to_address"`
	TokenAddress string    `json:"token_address"`
	Fee          *big.Int  `gorm:"serializer:u256;column:fee" db:"fee" json:"Fee" form:"fee"`
	Amount       *big.Int  `gorm:"serializer:u256;column:amount" db:"amount" json:"Amount" form:"amount"`
	Status       uint8     `json:"status"`  // 0:交易确认中,1:钱包交易已到账；2:交易已通知业务层；3:交易完成
	TxType       uint8     `json:"tx_type"` // 0:充值；1:提现；2:归集；3:热转冷；4:冷转热
	Timestamp    uint64
}

type TransactionsView interface {
	QueryTransactionByHash(hash string) (*Transactions, error)
}

type TransactionsDB interface {
	TransactionsView

	StoreTransactions([]Transactions, uint64) error
	UpdateTransactionsStatus(blockNumber *big.Int) error
	UpdateTransactionStatus(txList []Transactions) error
}

type transactionsDB struct {
	gorm *gorm.DB
}

func (db *transactionsDB) QueryTransactionByHash(hash string) (*Transactions, error) {
	var transactionEntry Transactions
	result := db.gorm.Table("transactions").Where("hash", hash).Take(&transactionEntry)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &transactionEntry, nil
}

func (db *transactionsDB) UpdateTransactionsStatus(blockNumber *big.Int) error {
	result := db.gorm.Model(&Transactions{}).Where("status = ? and block_number = ?", 0, blockNumber).Updates(map[string]interface{}{"status": gorm.Expr("GREATEST(1)")})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	return nil
}

func NewTransactionsDB(db *gorm.DB) TransactionsDB {
	return &transactionsDB{gorm: db}
}

func (db *transactionsDB) StoreTransactions(transactionsList []Transactions, transactionsLength uint64) error {
	result := db.gorm.CreateInBatches(&transactionsList, int(transactionsLength))
	return result.Error
}

func (db *transactionsDB) UpdateTransactionStatus(txList []Transactions) error {
	for i := 0; i < len(txList); i++ {
		var transactionSingle = Transactions{}

		result := db.gorm.Where(&Transactions{Hash: txList[i].Hash}).Take(&transactionSingle)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return result.Error
		}
		transactionSingle.Status = txList[i].Status
		transactionSingle.Fee = txList[i].Fee
		err := db.gorm.Save(&transactionSingle).Error
		if err != nil {
			return err
		}
	}
	return nil
}
