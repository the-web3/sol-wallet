package node

import (
	"fmt"
	"testing"
)

func newTestClient() *SolanaClient {
	client, _ := NewSolanaClient("https://docs-demo.solana-mainnet.quiknode.pro")
	return client
}

func TestSolanaClient_GetCurrentSlot(t *testing.T) {
	client := newTestClient()
	result, _ := client.GetCurrentSlot()
	fmt.Println("result======", result)
}

func TestSolanaClient_GetLatestBlockHeight(t *testing.T) {
	client := newTestClient()
	result, _ := client.GetLatestBlockHeight()
	fmt.Println("result======", result)
}

func TestSolanaClient_GetBlock(t *testing.T) {
	client := newTestClient()
	result, _ := client.GetBlock(258030759)
	for _, v := range result {
		fmt.Println("BlockHeight", v.BlockHeight)
		fmt.Println("BlockHash", v.BlockHash)
		fmt.Println("txHash", v.TxHash)
		fmt.Println("Source======", v.Source)
		fmt.Println("Destination======", v.Destination)
		fmt.Println("Fee======", v.Fee)
		fmt.Println("Amount======", v.Lamports)
	}
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetBalance("DttWaMuVvTiduZRnguLF7jNxTgiMBZ1hyAumKUiL2KRL")
	fmt.Println("balance======", balance)
}

func TestSolanaClient_GetNonce(t *testing.T) {
	client := newTestClient()
	nonce, _ := client.GetNonce("J2obR2DK7gnd6H88HjKzEYuMyboDWRNpbzwmGSh31nnu")
	fmt.Println("nonce==", nonce)
}

func TestSolanaClient_GetMinRent(t *testing.T) {
	client := newTestClient()
	minRent, _ := client.GetMinRent()
	fmt.Println("minRent==", minRent)
}
