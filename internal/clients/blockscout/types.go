// Copyright 2020 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
