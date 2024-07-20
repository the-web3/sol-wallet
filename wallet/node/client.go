package node

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/pkg/errors"

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

func (sol *SolanaClient) GetLatestBlockHeight() (int64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return int64(res.Result), nil
}

func (sol *SolanaClient) GetBlockInfo(startSlot uint64, endSlot uint64) (int64, error) {
	if endSlot-startSlot >= 500000 {
		return 0, errors.New("")
	}
	res, err := sol.RpcClient.GetBlocks(context.Background(), startSlot, endSlot)
	if err != nil {
		return 0, err
	}

	fmt.Println(res.Result)
	return 0, nil
}

func (sol *SolanaClient) GetBlock(slot uint64) (int64, error) {
	rewards := false
	var MaxSupportedTransactionVersion uint8 = 0
	res, err := sol.RpcClient.GetBlockWithConfig(context.Background(), slot, rpc.GetBlockConfig{
		Encoding:                       rpc.GetBlockConfigEncodingJsonParsed,
		TransactionDetails:             rpc.GetBlockConfigTransactionDetailsFull,
		Rewards:                        &rewards,
		MaxSupportedTransactionVersion: &MaxSupportedTransactionVersion,
	})
	if err != nil {
		return 0, err
	}
	fmt.Println(res.Result)
	return 0, nil
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
		log.Fatalf("failed to get balance with cfg, err: %v", err)
		return "", err
	}

	var lamportsOnAccount = new(big.Float).SetUint64(balance.Result.Value)

	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(1000000000))

	return solBalance.String(), nil
}

func (sol *SolanaClient) GetNonce(nonceAccount string) (string, error) {
	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), nonceAccount)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
		return "", err
	}
	return nonce, nil
}

func (sol *SolanaClient) GetMinRent() (string, error) {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		log.Fatalf("failed to get GetMinimumBalanceForRentExemption , err: %v", err)
		return "", err
	}
	return strconv.FormatUint(bal.Result, 10), nil
}

func getPreTokenBalance(preTokenBalance []rpc.TransactionMetaTokenBalance, accountIndex uint64) *rpc.TransactionMetaTokenBalance {
	for j := 0; j < len(preTokenBalance); j++ {
		preToken := preTokenBalance[j]
		if preToken.AccountIndex == accountIndex {
			return &preTokenBalance[j]
		}
	}
	return nil
}

func (sol *SolanaClient) GetTxByHash(hash string) (*TxMessage, error) {
	// "getTransaction" is only available in solana-core v1.7 or newer.
	// Please use getConfirmedTransaction for solana-core v1.6
	out, err := sol.RpcClient.GetTransaction(
		context.TODO(),
		hash,
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
		return nil, err
	}
	message := out.Result.Transaction.(map[string]interface{})["message"]
	accountKeys := message.((map[string]interface{}))["accountKeys"].([]interface{})
	signatures := out.Result.Transaction.(map[string]interface{})["signatures"].([]interface{})
	_hash := signatures[0]
	if out.Result.Meta.Err != nil || len(out.Result.Meta.LogMessages) == 0 || _hash == "" {
		log.Fatalf("not found tx, err: %v", err)
		return nil, err
	}

	var txMessage []*TxMessage
	for i := 0; i < len(accountKeys); i++ {
		to := accountKeys[i].(string)
		amount := out.Result.Meta.PostBalances[i] - out.Result.Meta.PreBalances[i]

		if to != "" && amount > 0 {
			txMessage = append(txMessage, &TxMessage{
				Hash:   hash,
				From:   "",
				To:     to,
				Fee:    strconv.FormatUint(out.Result.Meta.Fee, 10),
				Status: true,
				Value:  strconv.FormatInt(amount, 10),
				Type:   1,
				Height: strconv.FormatUint(out.Result.Slot, 10),
			})
		}
	}

	for i := 0; i < len(out.Result.Meta.PostTokenBalances); i++ {
		postToken := out.Result.Meta.PostTokenBalances[i]
		preTokenBalance := getPreTokenBalance(out.Result.Meta.PreTokenBalances, postToken.AccountIndex)
		if preTokenBalance == nil {
			continue
		}
		postAmount, _ := strconv.ParseFloat(postToken.UITokenAmount.Amount, 64)
		preAmount, _ := strconv.ParseFloat(preTokenBalance.UITokenAmount.Amount, 64)
		amount := postAmount - preAmount
		if amount > 0 {
			txMessage = append(txMessage, &TxMessage{
				Hash:   hash,
				From:   "",
				To:     postToken.Owner,
				Fee:    strconv.FormatUint(out.Result.Meta.Fee, 10),
				Status: true,
				Value:  strconv.FormatFloat(amount, 'E', -1, 10),
				Type:   1,
				Height: strconv.FormatUint(out.Result.Slot, 10),
			})
		}
	}
	if len(txMessage) > 0 {
		return txMessage[0], nil
	}
	log.Fatalf("not found tx, err: %v", err)
	return nil, err
}
