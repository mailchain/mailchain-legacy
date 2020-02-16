package blockscout

type txList struct {
	Status  string
	Message string
	Result  []txResult
}

type txResult struct {
	TimeStamp        string
	Hash             string
	Nonce            string
	BlockHash        string
	BlockNumber      string
	TransactionIndex string
	From             string
	To               string
	Gas              string
	GasPrice         string
	GasUsed          string
	IsError          string
	// "txreceipt_status": "",
	Input             string
	ContractAddress   string
	CumulativeGasUsed string
	Confirmations     string
	Value             string
}

type getTxInfo struct {
	Status  string
	Message string
	Result  txInfo
}

type txInfo struct {
	BlockHash     string
	BlockNumber   string
	Confirmations string
	From          string
	GasPrice      string
	GasUsed       string
	Hash          string
	Input         string
	Success       bool
	TimeStamp     string
	To            string
	Value         string
}
