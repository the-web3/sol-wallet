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
	result, _ := client.GetBlock(94101948)
	fmt.Println("result======", result)
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetBalance("8Lh2DVW5Lw3HgmZC55Fquno4K5auzSS7EveuLvEtCEXq")
	fmt.Println("balance======", balance)
}

func TestSolanaClient_GetTxByHash(t *testing.T) {
	client := newTestClient()
	txMessage, _ := client.GetTxByHash("3ESWyuEuTTMjQaG6GKvd37F4ZNcfeK5WbsvWZEPWD5dtjyoy3kcdaDAkg1UCZatUNaiE9boaCMMufTMVhtyNdcif")
	fmt.Println("", txMessage.From)
	fmt.Println("", txMessage.To)
	fmt.Println("", txMessage.Value)
}

func TestSolanaClient_RequestAirdrop(t *testing.T) {
	client := newTestClient()
	//client[0].RequestAirdrop("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	balance, _ := client.GetBalance("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	fmt.Println(balance)
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
