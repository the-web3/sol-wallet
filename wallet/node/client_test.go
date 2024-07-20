package node

import (
	"fmt"
	"testing"
)

func newTestClient() *SolanaClient {
	client, _ := NewSolanaClient("https://docs-demo.solana-mainnet.quiknode.pro")
	return client
}

func TestSolanaClient_GetLatestBlockHeight(t *testing.T) {
	client := newTestClient()
	result, _ := client.GetLatestBlockHeight()
	fmt.Println("result======", result)
}

func TestSolanaClient_GetBlock(t *testing.T) {
	client := newTestClient()
	result, _ := client.GetBlock(257859839)
	fmt.Println("result======", result)
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetBalance("8Lh2DVW5Lw3HgmZC55Fquno4K5auzSS7EveuLvEtCEXq")
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

func TestSolanaClient_GetTxByHash(t *testing.T) {
	client := newTestClient()
	tx := client.GetTxByHash("G6wz1rFZaGRbVUa9qPumYvmhNA3cxYXD8BCgZfztLfaJAAFP3rhQ74uEEza2wSSADBtiLHM5hoFD2jcAnaaYfiT")
	fmt.Println("tx===", tx)
}
