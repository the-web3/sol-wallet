package sign

type AddressList struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"publicKey"`
	Address    string `json:"address"`
}

type AccountInfoRep struct {
	Code      uint64        `json:"code"`
	Msg       string        `json:"msg"`
	Addresses []AddressList `json:"addressList"`
}

type KeyAddress struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

type PrepareAccountReq struct {
	AuthorAddress              string       `json:"authorAddress"`
	FromAddress                string       `json:"from"`
	RecentBlockhash            string       `json:"recentBlockhash"`
	MinBalanceForRentExemption uint64       `json:"minBalanceForRentExemption"`
	Privs                      []KeyAddress `json:"privs"`
}

type PrepareAccountRep struct {
	Code  uint64 `json:"code"`
	Msg   string `json:"msg"`
	RawTx string `json:"raw_tx"`
}

type TransactionReq struct {
	FromAddress  string `json:"from"`
	ToAddress    string `json:"to"`
	Amount       string `json:"amount"`
	NonceAccount string `json:"nonceAccount"`
	Nonce        string `json:"nonce"`
	Decimal      uint64 `json:"decimal"`
	PrivateKey   string `json:"privateKey"`
	MintAddress  string `json:"mintAddress"`
}

type TransactionRep struct {
	Code  uint64 `json:"code"`
	Msg   string `json:"msg"`
	RawTx string `json:"raw_tx"`
}
