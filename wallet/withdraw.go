package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/sol-wallet/common/tasks"
	"github.com/the-web3/sol-wallet/config"
	"github.com/the-web3/sol-wallet/database"
	"github.com/the-web3/sol-wallet/wallet/node"
	"github.com/the-web3/sol-wallet/wallet/retry"
	"github.com/the-web3/sol-wallet/wallet/sign"
)

type Withdraw struct {
	db             *database.DB
	chainConf      *config.ChainConfig
	client         node.SolanaClient
	signClient     *sign.Client
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewWithdraw(cfg *config.Config, db *database.DB, client node.SolanaClient, signCli *sign.Client, shutdown context.CancelCauseFunc) (*Withdraw, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Withdraw{
		db:             db,
		chainConf:      &cfg.Chain,
		client:         client,
		signClient:     signCli,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in deposit: %w", err))
		}},
	}, nil
}

func (w *Withdraw) Close() error {
	var result error
	w.resourceCancel()
	if err := w.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
	}
	return nil
}

func (w *Withdraw) Start() error {
	log.Info("start withdraw......")
	tickerWithdrawsWorker := time.NewTicker(time.Second * 5)
	w.tasks.Go(func() error {
		for range tickerWithdrawsWorker.C {

			withdrawList, err := w.db.Withdraws.UnSendWithdrawsList()
			if err != nil {
				log.Error("get un send withdraw list fail", "err", err)
				return err
			}

			returnWithdrawsList := make([]database.Withdraws, len(withdrawList))
			index := 0
			var balanceList []database.Balances
			for _, withdraw := range withdrawList {
				hotWallet, err := w.db.Addresses.QueryHotWalletInfo()
				if err != nil {
					log.Error("query hot wallet info err", "err", err)
					return err
				}

				hotWalletTokenBalance, err := w.db.Balances.QueryWalletBalanceByTokenAndAddress(hotWallet.Address, withdraw.TokenAddress)
				if hotWalletTokenBalance.Balance.Cmp(withdraw.Amount) < 0 {
					log.Info("hot wallet balance is not enough", "tokenAddress", withdraw.TokenAddress)
					continue
				}

				recentBlockhash, err := w.client.GetRecentBlockHash()
				if err != nil {
					log.Error("query nonce by address fail", "err", err)
					return err
				}

				txReq := &sign.TransactionReq{
					FromAddress:  hotWalletTokenBalance.Address,
					ToAddress:    withdraw.ToAddress,
					Amount:       withdraw.Amount.String(),
					NonceAccount: hotWalletTokenBalance.Address,
					Nonce:        recentBlockhash,
					Decimal:      9,
					PrivateKey:   hotWallet.PrivateKey,
					MintAddress:  withdraw.TokenAddress,
				}

				txRep, err := w.signClient.SignTransaction(txReq)
				if err != nil {
					log.Error("sign transaction fail", "err", err)
					return err
				}
				if txRep.Code == 2000 {
					// 发送交易到区块链网络
					txHash, err := w.client.SendRawTransaction(txRep.RawTx)
					if err != nil {
						log.Error("send raw transaction fail", "err", err)
						return err
					}
					returnWithdrawsList[index].Hash = txHash
					returnWithdrawsList[index].GUID = withdraw.GUID
					balanceItem := database.Balances{
						Address:      hotWallet.Address,
						TokenAddress: withdraw.TokenAddress,
						LockBalance:  withdraw.Amount,
					}
					balanceList = append(balanceList, balanceItem)
					index++
				} else {
					log.Error("sign service occur unknown err")
					continue
				}
			}
			retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
			if _, err := retry.Do[interface{}](w.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
				if err := w.db.Transaction(func(tx *database.DB) error {
					// 将转出去的热钱包余额锁定
					err = w.db.Balances.UpdateBalances(balanceList, false)
					if err != nil {
						log.Error("mark withdraw send fail", "err", err)
						return err
					}

					err = w.db.Withdraws.MarkWithdrawsToSend(returnWithdrawsList)
					if err != nil {
						log.Error("mark withdraw send fail", "err", err)
						return err
					}
					return nil
				}); err != nil {
					log.Error("unable to persist batch", "err", err)
					return nil, err
				}
				return nil, nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
