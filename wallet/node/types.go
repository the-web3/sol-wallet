package node

import "math/big"

type TransactionDetail struct {
	TxHash      string   `json:"tx_hash"`
	Destination string   `json:"destination"`
	Source      string   `json:"source"`
	Lamports    *big.Int `json:"lamports"`
	Type        string   `json:"type"`
	Fee         *big.Int `json:"fee"`
}
