package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/the-web3/sol-wallet/wallet/retry"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/sol-wallet/common/tasks"
	"github.com/the-web3/sol-wallet/config"
	"github.com/the-web3/sol-wallet/database"
	"github.com/the-web3/sol-wallet/wallet/node"
)

type Deposit struct {
	db        *database.DB
	chainConf *config.ChainConfig

	client node.SolanaClient

	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewDeposit(cfg *config.Config, db *database.DB, client node.SolanaClient, shutdown context.CancelCauseFunc) (*Deposit, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &Deposit{
		db:             db,
		chainConf:      &cfg.Chain,
		client:         client,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("critical error in deposit: %w", err))
		}},
	}, nil
}

func (d *Deposit) Close() error {
	var result error
	d.resourceCancel()
	if err := d.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("failed to await deposit %w"), err)
		return result
	}
	return nil
}

func (d *Deposit) Start() error {
	log.Info("start deposit......")
	tickerDepositWorker := time.NewTicker(time.Second * 5)
	d.tasks.Go(func() error {
		for range tickerDepositWorker.C {
			// 获取最新区块高度，并且获取数据库里面上次同步到高度，比较这两个高度，如果数据库里面的高度等于最新区块高度，不再往下执行交易解析，继续扫描最新的块
			// 如果是第一次进入，那么以配置起始高度开始网上同步，若起始高度没有配置或者配置是 0，那么就是 0 开始同步
			// 每次同步按照配置的同步步长往下执行
			var startSyncBlock *big.Int
			dbLastestBlock, err := d.db.Blocks.LatestBlocks()
			if err != nil {
				log.Error("get latest block from database fail", "err", err)
				return err
			}
			if dbLastestBlock == nil {
				startSyncBlock = big.NewInt(int64(d.chainConf.StartingHeight))
			} else {
				startSyncBlock = dbLastestBlock.Number
			}

			chainLatestBlock, err := d.client.GetLatestBlockHeight()
			if err != nil {
				log.Error("get latest block from solana chain fail", "err", err)
				return err
			}

			if startSyncBlock.Cmp(big.NewInt(int64(chainLatestBlock))) >= 0 {
				continue
			}

			// 按照步长处理
			endSyncBlock := new(big.Int).Add(startSyncBlock, big.NewInt(int64(d.chainConf.BlocksStep)))

			blocks, deposits, withdraws, depositTransactions, outherTransactions, tokenBalances, err := d.processTransactions(startSyncBlock, endSyncBlock)
			if err != nil {
				log.Error("process transaction fail", "err", err)
				return err
			}

			retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
			if _, err := retry.Do[interface{}](d.resourceCtx, 10, retryStrategy, func() (interface{}, error) {
				if err := d.db.Transaction(func(tx *database.DB) error {
					if err := tx.Blocks.StoreBlockss(blocks, uint64(len(blocks))); err != nil {
						return err
					}

					if len(deposits) > 0 {
						log.Info("Store deposit transaction success", "totalTx", len(deposits))
						if err := tx.Deposits.StoreDeposits(deposits, uint64(len(deposits))); err != nil {
							return err
						}
					}
					log.Info("batch latest block number", "endSyncBlock", endSyncBlock)

					// 更新之前充值确认位
					if err := tx.Deposits.UpdateDepositsStatus(endSyncBlock.Uint64() - uint64(d.chainConf.Confirmations)); err != nil {
						return err
					}

					if len(withdraws) > 0 {
						if err := tx.Withdraws.UpdateTransactionStatus(withdraws); err != nil {
							return err
						}
					}

					if len(depositTransactions) > 0 {
						if err := tx.Transactions.StoreTransactions(depositTransactions, uint64(len(depositTransactions))); err != nil {
							return err
						}
					}

					if len(outherTransactions) > 0 { // 提现和归集
						if err := tx.Transactions.UpdateTransactionStatus(outherTransactions); err != nil {
							return err
						}
					}

					if len(tokenBalances) > 0 {
						log.Info("update or store token balance", "tokenBalanceList", len(tokenBalances))
						if err := tx.Balances.UpdateOrCreate(tokenBalances); err != nil {
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
		}
		return nil
	})
	return nil
}

func (d *Deposit) processTransactions(startSyncBlock, endSyncBlock *big.Int) ([]database.Blocks, []database.Deposits, []database.Withdraws, []database.Transactions, []database.Transactions, []database.TokenBalance, error) {

	for index := startSyncBlock.Uint64(); index <= endSyncBlock.Uint64(); index++ {
		tx, err := d.client.GetBlock(index)
		if err != nil {
			log.Error("get block info faill", err)
			continue
		}

		fromAddress, err := d.db.Addresses.QueryAddressesByToAddress(tx.Source)
		if err != nil {
			log.Error("query token info fail", "err", err)
			continue
		}
		log.Info("query from address success", "source", tx.Source)

		toAddress, err := d.db.Addresses.QueryAddressesByToAddress(tx.Destination)
		if err != nil {
			log.Error("query to address info fail", "err", err)
			continue
		}
		log.Info("query to address success", "Destination", tx.Destination)

		if fromAddress == nil && toAddress == nil {
			continue
		}

	}

	return nil, nil, nil, nil, nil, nil, nil
}
