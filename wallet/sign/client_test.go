package sign

import (
	"fmt"
	"testing"
)

const (
	testBaseUrl = "http://127.0.0.1:3000"
)

func TestClient_GenerateAddress(t *testing.T) {
	client, err := NewSolSignClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	accountRet, err := client.GenerateAddress(10)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range accountRet.Addresses {
		fmt.Println("privateKey", value.PrivateKey)
		fmt.Println("publicKey", value.PublicKey)
		fmt.Println("address", value.Address)
	}
}

func TestClient_PrepareAccount(t *testing.T) {
	client, err := NewSolSignClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	var keyAddressList []KeyAddress
	keyAddressOne := KeyAddress{
		Key:     "privateKey",
		Address: "FvjWo4jbdsAP4ZHtJfiUpv5xb6TpBWRtDASGPmKKR39E",
	}
	keyAddressList = append(keyAddressList, keyAddressOne)

	keyAddressTwo := KeyAddress{
		Key:     "privateKey",
		Address: "4wHd9tf4x4FkQ3JtgsMKyiEofEHSaZH5rYzfFKLvtESD",
	}
	keyAddressList = append(keyAddressList, keyAddressTwo)
	par := &PrepareAccountReq{
		AuthorAddress:              "4wHd9tf4x4FkQ3JtgsMKyiEofEHSaZH5rYzfFKLvtESD",
		FromAddress:                "FvjWo4jbdsAP4ZHtJfiUpv5xb6TpBWRtDASGPmKKR39E",
		RecentBlockhash:            "CSL1MJGUcDbgUEHh6fPsxum42vkhnQCh62whKjEiGwR3",
		MinBalanceForRentExemption: 1647680,
		Privs:                      keyAddressList,
	}
	parRet, err := client.PrepareAccount(par)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code====", parRet.Code)
	fmt.Println("Msg====", parRet.Msg)
	fmt.Println("RawTx==", parRet.RawTx)
}

func TestClient_SignTransaction(t *testing.T) {
	client, err := NewSolSignClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	signParam := &TransactionReq{
		FromAddress:  "4wHd9tf4x4FkQ3JtgsMKyiEofEHSaZH5rYzfFKLvtESD",
		ToAddress:    "FvjWo4jbdsAP4ZHtJfiUpv5xb6TpBWRtDASGPmKKR39E",
		Amount:       "0.01",
		NonceAccount: "FvjWo4jbdsAP4ZHtJfiUpv5xb6TpBWRtDASGPmKKR39E",
		Nonce:        "GGLM3xu9yXzDoH3uhMEPcqqju6BB6C1FzoWKdiji5x5t",
		Decimal:      9,
		PrivateKey:   "privateKey",
		MintAddress:  "",
	}
	signRet, err := client.SignTransaction(signParam)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Code====", signRet.Code)
	fmt.Println("Msg====", signRet.Msg)
	fmt.Println("RawTx==", signRet.RawTx)
}
