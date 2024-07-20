package database

import (
	"math/big"
)

type TokenBalance struct {
	Address      string   `json:"address"`
	TokenAddress string   `json:"token_address"`
	Balance      *big.Int `json:"balance"`
	LockBalance  *big.Int `json:"lock_balance"`
	TxType       uint8    `json:"tx_type"` // 0:充值；1:提现；2:归集；3:热转冷；4:冷转热
}
