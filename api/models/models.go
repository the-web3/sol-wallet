package models

import (
	"github.com/the-web3/sol-wallet/database"
	"math/big"
)

type SubmitDWParams struct {
	FromAddress  string
	ToAddress    string
	TokenAddress string
	Amount       *big.Int
}

type QueryDWParams struct {
	Address  string
	Page     int
	PageSize int
	Order    string
}

type QueryPageParams struct {
	Page     int
	PageSize int
	Order    string
}

type QueryIdParams struct {
	Id uint64
}

type QueryIndexParams struct {
	Index uint64
}

type DepositsResponse struct {
	Current int                 `json:"Current"`
	Size    int                 `json:"Size"`
	Total   int64               `json:"Total"`
	Records []database.Deposits `json:"Records"`
}

type WithdrawsResponse struct {
	Current int                  `json:"Current"`
	Size    int                  `json:"Size"`
	Total   int64                `json:"Total"`
	Records []database.Withdraws `json:"Records"`
}

type SubmitWithdrawsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
