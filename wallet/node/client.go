package node

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

type SolanaClient struct {
	RpcClient rpc.RpcClient
	Client    *client.Client
}

func NewSolanaClient(url string) (*SolanaClient, error) {
	rpcClient := rpc.NewRpcClient(url)
	clientNew := client.NewClient(url)
	return &SolanaClient{
		RpcClient: rpcClient,
		Client:    clientNew,
	}, nil
}

// GetRecentBlockHash 获取最新的区块链
func (sol *SolanaClient) GetRecentBlockHash() (string, error) {
	return "", nil
}

// GetCurrentSlot 获取最新的区块链
func (sol *SolanaClient) GetCurrentSlot() (uint64, error) {
	res, err := sol.RpcClient.GetSlot(context.Background())
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetLatestBlockHeight 获取最新的区块链
func (sol *SolanaClient) GetLatestBlockHeight() (uint64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetBlock 根据区块号获取里面的交易
func (sol *SolanaClient) GetBlock(slot uint64) ([]TransactionDetail, error) {
	rewards := false
	var MaxSupportedTransactionVersion uint8 = 0
	res, err := sol.RpcClient.GetBlockWithConfig(context.Background(), slot, rpc.GetBlockConfig{
		Encoding:                       rpc.GetBlockConfigEncodingJsonParsed,
		TransactionDetails:             rpc.GetBlockConfigTransactionDetailsFull,
		Rewards:                        &rewards,
		MaxSupportedTransactionVersion: &MaxSupportedTransactionVersion,
	})
	if err != nil {
		return nil, err
	}
	var txDetailList []TransactionDetail
	if res.Result.Transactions != nil {
		for _, value := range res.Result.Transactions {
			if convertedMap, ok := value.Transaction.(map[string]interface{}); ok {
				message := convertedMap["message"].(map[string]interface{})
				instructions := message["instructions"].([]interface{})
				for _, instruction := range instructions {
					instructionItem := instruction.(map[string]interface{})
					if instructionItem["program"] == "spl-token" || instructionItem["program"] == "system" { // token transfer
						txType := instructionItem["parsed"].(map[string]interface{})["type"]
						if txType != "transfer" {
							continue
						} else {
							var fromAddres, toAddress string
							amount := new(big.Int)
							information := instructionItem["parsed"].(map[string]interface{})["info"].(map[string]interface{})
							fromAddres = information["source"].(string)
							toAddress = information["destination"].(string)
							if instructionItem["program"] == "spl-token" {
								amountStr := information["amount"].(string)
								amount.SetString(amountStr, 10)
							} else {
								amountStr := fmt.Sprintf("%.0f", information["lamports"].(float64))
								amount.SetString(amountStr, 10)
							}
							fee := value.Meta.Fee
							signatures := convertedMap["signatures"].([]interface{})
							blockHeight := res.Result.ParentSlot + 1
							txDetail := TransactionDetail{
								PreviousBlockhash: res.Result.PreviousBlockhash,
								BlockHash:         res.Result.Blockhash,
								BlockHeight:       big.NewInt(int64(blockHeight)),
								TxHash:            signatures[0].(string),
								Destination:       toAddress,
								Source:            fromAddres,
								Lamports:          amount,
								Type:              txType.(string),
								Fee:               big.NewInt(int64(fee)),
							}
							txDetailList = append(txDetailList, txDetail)
						}
					} else {
						continue
					}
				}
			}
		}
	}
	return txDetailList, err
}

func (sol *SolanaClient) GetBalance(address string) (string, error) {
	balance, err := sol.RpcClient.GetBalanceWithConfig(
		context.TODO(),
		address,
		rpc.GetBalanceConfig{
			Commitment: rpc.CommitmentProcessed,
		},
	)
	if err != nil {
		return "", err
	}

	var lamportsOnAccount = new(big.Float).SetUint64(balance.Result.Value)

	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(1000000000))

	return solBalance.String(), nil
}

func (sol *SolanaClient) GetNonce(nonceAccount string) (string, error) {
	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), nonceAccount)
	if err != nil {
		return "", err
	}
	return nonce, nil
}

func (sol *SolanaClient) GetMinRent() (string, error) {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(bal.Result, 10), nil
}

func (sol *SolanaClient) SendRawTransaction(rawTx string) (string, error) {
	bal, err := sol.RpcClient.SendTransaction(context.Background(), rawTx)
	if err != nil {
		return "", err
	}
	return bal.Result, nil
}
