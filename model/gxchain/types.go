package gxchain

{"id":1,"jsonrpc":"2.0","result":{"id":"1.6.43","witness_account":"1.2.1059405","last_aslot":17470511,"signing_key":"GXC8KRoN5ovHbZBGCEphaQcBNCQ8Bhp8BnvzPeZvUKzt86VDS5AiP","pay_vb":"1.13.263","vote_id":"1:73","total_votes":"835856507024","url":"https://gxcdac.com","total_missed":69,"last_confirmed_block_num":17260439,"is_valid":true}}

type _WitnessType struct{
	Id             string `json:"id"`
	Account        string `json:"witness_account"`
	LastAslot      int    `json:"last_aslot"`
	SigningKey     string `json:"signing_key"`
	PayVb          string `json:"pay_vb"`
	VoteId         string `json:"vote_id"`
	TotalVotes     string `json:"total_votes"`
	Url            string `json:"url"`
	TotalMissed    int    `json:"total_missed"`
	LastConfirmed  int    `json:"last_confirmed_block_num"`
	Valid          bool   `json:"is_valid"`
}

type WitnessType struct{
	Id  int `json:"id"`
	Jsonrpc  string `json:"jsonrpc"`
	Result   _WitnessType `json:"result"`
}
