package blockscout

type txList struct {
	Status  string
	Message string
	Result  []txResult
}

type txResult struct {
	BlockNumber      string
	TimeStamp        string
	Hash             string
	Nonce            string
	BlockHash        string
	TransactionIndex string
	From             string
	To               string
	Value            string
	Gas              string
	GasPrice         string
	IsError          string
	// "txreceipt_status": "",
	Input             string
	ContractAddress   string
	CumulativeGasUsed string
	GasUsed           string
	Confirmations     string
}
