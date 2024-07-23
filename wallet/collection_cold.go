package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-web3/sol-wallet/wallet/sign"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/sol-wallet/common/tasks"
	"github.com/the-web3/sol-wallet/config"
	"github.com/the-web3/sol-wallet/database"
	"github.com/the-web3/sol-wallet/wallet/node"
	"github.com/the-web3/sol-wallet/wallet/retry"
)

var (
	CollectionFunding = big.NewInt(10000000000000000)
	ColdFunding       = big.NewInt(2000000000000000000)
)

type CollectionCold struct {
	db             *database.DB
	chainConf      *config.ChainConfig
	client         node.SolanaClient
	signClient     *sign.Client
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewCollectionCold(cfg *config.Config, db *database.DB, client node.SolanaClient, signCli *sign.Client, shutdown context.CancelCauseFunc) (*CollectionCold, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &CollectionCold{
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

func (cc *CollectionCold) Close() error {
	var result error
	cc.resourceCancel()
	if err := cc.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
	}
	return nil
}

func (cc *CollectionCold) Start() error {
	log.Info("start collection and cold......")
	tickerCollectionColdWorker := time.NewTicker(time.Second * 5)
	cc.tasks.Go(func() error {
		for range tickerCollectionColdWorker.C {
			err := cc.Collection()
			if err != nil {
				log.Error("collect fail", "err", err)
				return err
			}
		}
		return nil
	})

	cc.tasks.Go(func() error {
		for range tickerCollectionColdWorker.C {
			err := cc.ToCold()
			if err != nil {
				log.Error("to cold fail", "err", err)
				return err
			}
		}
		return nil
	})

	return nil
}

func (cc *CollectionCold) ToCold() error {
	hotWalletBalancesList, err := cc.db.Balances.QueryHotWalletBalances(ColdFunding)
	if err != nil {
		log.Error("to cold query hot wallet info fail", "err", err)
		return err
	}
	var txList []database.Transactions
	balanceForStore := make([]database.Balances, len(hotWalletBalancesList))
	for _, value := range hotWalletBalancesList {
		index := 0
		coldWalletInfo, err := cc.db.Addresses.QueryColdWalletInfo()
		if err != nil {
			log.Error("query cold wallet info err", "err", err)
			return err
		}

		// nonce
		recentBlockhash, err := cc.client.GetRecentBlockHash()
		if err != nil {
			log.Error("query nonce by address fail", "err", err)
			return err
		}

		hotAccount, err := cc.db.Addresses.QueryHotWalletInfo()
		if err != nil {
			log.Error("query account info by address fail", "err", err)
			return err
		}

		//  sendRawTx
		txReq := &sign.TransactionReq{
			FromAddress:  hotAccount.Address,
			ToAddress:    coldWalletInfo.Address,
			Amount:       value.Balance.String(),
			NonceAccount: hotAccount.Address,
			Nonce:        recentBlockhash,
			Decimal:      9,
			PrivateKey:   hotAccount.PrivateKey,
			MintAddress:  value.TokenAddress,
		}

		txRep, err := cc.signClient.SignTransaction(txReq)
		if err != nil {
			log.Error("sign transaction fail", "err", err)
			return err
		}

		if txRep.Code != 2000 {
			log.Error("sign server occur unknown error")
			continue
		}

		txHash, err := cc.client.SendRawTransaction(txRep.RawTx)
		if err != nil {
			log.Error("send raw transaction fail", "err", err)
			return err
		}

		guid, _ := uuid.NewUUID()
		coldTx := database.Transactions{
			GUID:         guid,
			BlockHash:    "",
			BlockNumber:  nil,
			Hash:         txHash,
			FromAddress:  value.Address,
			ToAddress:    coldWalletInfo.Address,
			TokenAddress: value.TokenAddress,
			Fee:          big.NewInt(0),
			Amount:       value.Balance,
			Status:       0,
			TxType:       2,
			Timestamp:    uint64(time.Time{}.Unix()),
		}
		txList = append(txList, coldTx)
		balanceForStore[index].LockBalance = new(big.Int).Sub(balanceForStore[index].Balance, ColdFunding)
		balanceForStore[index].Address = value.Address
		balanceForStore[index].TokenAddress = value.TokenAddress
	}
	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](cc.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := cc.db.Transaction(func(tx *database.DB) error {
			if len(hotWalletBalancesList) > 0 {
				if err := tx.Balances.UpdateBalances(balanceForStore, false); err != nil {
					return err
				}
			}
			if len(txList) > 0 {
				if err := tx.Transactions.StoreTransactions(txList, uint64(len(txList))); err != nil {
					return err
				}
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
	return nil
}

// Collection 归集
func (cc *CollectionCold) Collection() error {
	unCollectionList, err := cc.db.Balances.UnCollectionList(CollectionFunding)
	if err != nil {
		log.Error("query uncollection fail", "err", err)
		return err
	}

	hotWalletInfo, err := cc.db.Addresses.QueryHotWalletInfo()
	if err != nil {
		log.Error("query hot wallet info fail", "err", err)
		return err
	}

	var txList []database.Transactions
	for _, uncollect := range unCollectionList {
		accountInfo, err := cc.db.Addresses.QueryAddressesByToAddress(uncollect.Address)
		if err != nil {
			log.Error("query account info fail", "err", err)
			return err
		}

		// nonce
		recentBlockHash, err := cc.client.GetRecentBlockHash()
		if err != nil {
			log.Error("query nonce by address fail", "err", err)
			return err
		}

		//  sendRawTx
		txReq := &sign.TransactionReq{
			FromAddress:  uncollect.Address,
			ToAddress:    hotWalletInfo.Address,
			Amount:       uncollect.Balance.String(),
			NonceAccount: uncollect.Address,
			Nonce:        recentBlockHash,
			Decimal:      9,
			PrivateKey:   accountInfo.PrivateKey,
			MintAddress:  uncollect.TokenAddress,
		}

		txRep, err := cc.signClient.SignTransaction(txReq)
		if err != nil {
			log.Error("sign transaction fail", "err", err)
			return err
		}

		if txRep.Code != 2000 {
			log.Error("sign server occur unknown error")
			continue
		}

		txHash, err := cc.client.SendRawTransaction(txRep.RawTx)
		if err != nil {
			log.Error("send raw transaction fail", "err", err)
			return err
		}
		guid, _ := uuid.NewUUID()
		collection := database.Transactions{
			GUID:         guid,
			BlockHash:    "",
			BlockNumber:  big.NewInt(1),
			Hash:         txHash,
			FromAddress:  uncollect.Address,
			ToAddress:    hotWalletInfo.Address,
			TokenAddress: uncollect.TokenAddress,
			Fee:          big.NewInt(1),
			Amount:       uncollect.Balance,
			Status:       0,
			TxType:       2,
			Timestamp:    uint64(time.Now().Unix()),
		}
		txList = append(txList, collection)
	}
	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](cc.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
		if err := cc.db.Transaction(func(tx *database.DB) error {
			if len(unCollectionList) > 0 {
				if err := tx.Balances.UpdateBalances(unCollectionList, true); err != nil {
					return err
				}
			}

			if err := tx.Transactions.StoreTransactions(txList, uint64(len(txList))); err != nil {
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
	return nil
}
