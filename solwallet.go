package sol_wallet

import (
	"context"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/sol-wallet/config"
	"github.com/the-web3/sol-wallet/database"
	"github.com/the-web3/sol-wallet/wallet"
	"github.com/the-web3/sol-wallet/wallet/node"
	"github.com/the-web3/sol-wallet/wallet/sign"
)

type SolWallet struct {
	deposit        *wallet.Deposit
	withdraw       *wallet.Withdraw
	collectionCold *wallet.CollectionCold

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewSolWallet(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*SolWallet, error) {
	solClient, err := node.NewSolanaClient(cfg.Chain.RpcUrl)
	if err != nil {
		return nil, err
	}

	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("init database fail", err)
		return nil, err
	}

	signCli, err := sign.NewSolSignClient(cfg.SignServerProvider)
	if err != nil {
		log.Error("new sign client fail", "err", err)
		return nil, err
	}

	deposit, err := wallet.NewDeposit(cfg, db, *solClient, shutdown)
	if err != nil {
		log.Error("new deposit fail", "err", err)
		return nil, err
	}
	withdraw, err := wallet.NewWithdraw(cfg, db, *solClient, signCli, shutdown)
	if err != nil {
		log.Error("new withdraw fail", "err", err)
		return nil, err
	}
	collectionCold, err := wallet.NewCollectionCold(cfg, db, *solClient, signCli, shutdown)
	if err != nil {
		log.Error("new collection and to cold fail", "err", err)
		return nil, err
	}

	out := &SolWallet{
		deposit:        deposit,
		withdraw:       withdraw,
		collectionCold: collectionCold,
		shutdown:       shutdown,
	}

	return out, nil
}

func (ew *SolWallet) Start(ctx context.Context) error {
	err := ew.deposit.Start()
	if err != nil {
		return err
	}
	err = ew.withdraw.Start()
	if err != nil {
		return err
	}
	err = ew.collectionCold.Start()
	if err != nil {
		return err
	}
	return nil
}

func (ew *SolWallet) Stop(ctx context.Context) error {
	err := ew.deposit.Close()
	if err != nil {
		return err
	}
	err = ew.withdraw.Close()
	if err != nil {
		return err
	}

	err = ew.collectionCold.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ew *SolWallet) Stopped() bool {
	return ew.stopped.Load()
}
