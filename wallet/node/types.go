package node

type Header struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}

type Instructions struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIDIndex int    `json:"programIdIndex"`
}

type Message struct {
	AccountKeys     []string       `json:"accountKeys"`
	Header          Header         `json:"header"`
	Instructions    []Instructions `json:"instructions"`
	RecentBlockhash string         `json:"recentBlockhash"`
}

type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

type TxMessage struct {
	Hash   string
	From   string
	To     string
	Fee    string
	Status bool
	Value  string
	Type   int32
	Height string
}
