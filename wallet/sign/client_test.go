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
		Key:     "d69df4d566d4797b9623c92ccb564c660cf6ab2ab6fbfc1181cf3315e50812b2ddc778616699176f05ab8dc90dc640ae41fbee549178b869974d4a7deaa52745",
		Address: "FvjWo4jbdsAP4ZHtJfiUpv5xb6TpBWRtDASGPmKKR39E",
	}
	keyAddressList = append(keyAddressList, keyAddressOne)

	keyAddressTwo := KeyAddress{
		Key:     "55a70321542da0b6123f37180e61993d5769f0a5d727f9c817151c1270c290963a7b3874ba467be6b81ea361e3d7453af8b81c88aedd24b5031fdda0bc71ad32",
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
		PrivateKey:   "55a70321542da0b6123f37180e61993d5769f0a5d727f9c817151c1270c290963a7b3874ba467be6b81ea361e3d7453af8b81c88aedd24b5031fdda0bc71ad32",
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
