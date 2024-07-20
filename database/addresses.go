package database

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Addresses struct {
	GUID        uuid.UUID `gorm:"primaryKey" json:"guid"`
	UserUid     string    `json:"user_uid"`
	Address     string    `json:"address"`
	AddressType uint8     `json:"address_type"` //0:用户地址；1:热钱包地址(归集地址)；2:冷钱包地址
	PrivateKey  string    `json:"private_key"`
	PublicKey   string    `json:"public_key"`
	Timestamp   uint64
}

type AddressesView interface {
	QueryAddressesByToAddress(string) (*Addresses, error)
	QueryHotWalletInfo() (*Addresses, error)
	QueryColdWalletInfo() (*Addresses, error)
}

type AddressesDB interface {
	AddressesView

	StoreAddressess([]Addresses, uint64) error
}

type addressesDB struct {
	gorm *gorm.DB
}

func (db *addressesDB) QueryAddressesByToAddress(address string) (*Addresses, error) {
	var addressEntry Addresses
	err := db.gorm.Table("addresses").Where("address", address).Take(&addressEntry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &addressEntry, nil
}

func NewAddressesDB(db *gorm.DB) AddressesDB {
	return &addressesDB{gorm: db}
}

func (db *addressesDB) StoreAddressess(addressList []Addresses, addressLength uint64) error {
	result := db.gorm.CreateInBatches(&addressList, int(addressLength))
	return result.Error
}

func (db *addressesDB) QueryHotWalletInfo() (*Addresses, error) {
	var addressEntry Addresses
	err := db.gorm.Table("addresses").Where("address_type", 1).Take(&addressEntry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &addressEntry, nil
}

func (db *addressesDB) QueryColdWalletInfo() (*Addresses, error) {
	var addressEntry Addresses
	err := db.gorm.Table("addresses").Where("address_type", 2).Take(&addressEntry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &addressEntry, nil
}
