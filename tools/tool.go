package tools

import (
	"context"
	"math/big"

	"github.com/urfave/cli/v2"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/sol-wallet/config"
	"github.com/the-web3/sol-wallet/database"
	"github.com/the-web3/sol-wallet/wallet/retry"
	"github.com/the-web3/sol-wallet/wallet/sign"
)

const GenerateAddressNum = 100

func CreateAddressTools(ctx *cli.Context, cfg *config.Config, db *database.DB) error {
	log.Info("start tools", "cfg.SignServerProvider", cfg.SignServerProvider)
	client, err := sign.NewSolSignClient(cfg.SignServerProvider)
	if err != nil {
		log.Error("New sol sign client fail", "err", err)
		return err
	}

	accountInfo, err := client.GenerateAddress(GenerateAddressNum)
	if err != nil {
		log.Error("generate address fail", "err", err)
	}

	if accountInfo != nil && accountInfo.Code != 2000 {
		log.Error("return code err", "code", accountInfo.Code)
		return err
	}

	var addressList []database.Addresses
	var balanceList []database.Balances

	index := 0
	for _, address := range accountInfo.Addresses {
		var AddressType uint8
		var UserUid string
		if index == 1 {
			AddressType = 1
			UserUid = "hot-wallet-for-the-web3"
		} else if index == 2 {
			AddressType = 2
			UserUid = "cold-wallet-for-the-web3"
		} else {
			UserUid = "useruid"
			AddressType = 0
		}
		addressItem := database.Addresses{
			GUID:        uuid.New(),
			UserUid:     UserUid,
			Address:     address.Address,
			AddressType: AddressType,
			PrivateKey:  address.PrivateKey,
			PublicKey:   address.PublicKey,
			Timestamp:   uint64(index + 10000),
		}
		addressList = append(addressList, addressItem)

		balanceItem := database.Balances{
			GUID:         uuid.New(),
			Address:      address.Address,
			TokenAddress: "",
			AddressType:  AddressType,
			Balance:      big.NewInt(0),
			LockBalance:  big.NewInt(0),
			Timestamp:    uint64(index + 10000),
		}
		balanceList = append(balanceList, balanceItem)
		index++
	}
	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	if _, err := retry.Do[interface{}](context.Background(), 10, retryStrategy, func() (interface{}, error) {
		if err := db.Transaction(func(tx *database.DB) error {
			err = db.Addresses.StoreAddressess(addressList, uint64(len(addressList)))
			if err != nil {
				log.Error("store address error", err)
				return err
			}
			err = db.Balances.StoreBalances(balanceList, uint64(len(balanceList)))
			if err != nil {
				log.Error("store balances error", err)
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
