package node

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/log"

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

// GetLatestBlockHeight 获取最新的区块链
func (sol *SolanaClient) GetLatestBlockHeight() (uint64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetBlock 根据区块号获取里面的交易
func (sol *SolanaClient) GetBlock(slot uint64) (*Transaction, error) {
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
	if res.Result.Transactions != nil {
		fmt.Println(res.Result.Transactions[0].Transaction)
		if convertedMap, ok := res.Result.Transactions[0].Transaction.(map[string]interface{}); ok {
			fmt.Println("Converted map:", convertedMap["instructions"])
		}
	}
	return nil, err
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

// GetTxByHash "getTransaction" is only available in solana-core v1.7 or newer.
// Please use getConfirmedTransaction for solana-core v1.6
// 根据交易 Hash 获取交易记录的详情
func (sol *SolanaClient) GetTxByHash(hash string) error {
	var MaxSupportedTransactionVersion uint8 = 0
	out, err := sol.RpcClient.GetTransactionWithConfig(
		context.TODO(),
		hash,
		rpc.GetTransactionConfig{
			Encoding:                       rpc.TransactionEncodingJsonParsed,
			MaxSupportedTransactionVersion: &MaxSupportedTransactionVersion,
		})
	if err != nil {
		return err
	}
	log.Info("get tx by hash", "out", out)
	return nil
}
