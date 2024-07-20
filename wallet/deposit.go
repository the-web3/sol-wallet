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
			log.Info("solana deposit..............")
		}
		return nil
	})
	return nil
}
