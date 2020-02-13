package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	extHttp "gitlab.com/sdce/exlib/http"
)

const (
	ApiInitTransaction  = "Initiate"
	ApiGetTransaction   = "GetTransaction?token="
	ApiDailyTransaction = "GetDailyTransactions?date="
)

type PoliClient struct {
	Config *PoliConfig
}

func NewPoliClient(config *PoliConfig) *PoliClient {
	return &PoliClient{
		Config: config,
	}
}

func (p *PoliClient) InitTransaction(in *InitTransactionRequest) (out *InitTransactionResponse, err error) {

	client := &http.Client{}
	rawReq, err := json.Marshal(*in)
	if err != nil {
		return
	}
	req, _ := http.NewRequest("POST", p.Config.ApiUrl+"/Initiate", bytes.NewBuffer(rawReq))
	req.SetBasicAuth(p.Config.MerchangeCode, p.Config.MerchantAuthCode)
	req.Header.Add(extHttp.Accept, extHttp.ApplicationJson)
	req.Header.Add(extHttp.ContentType, extHttp.ApplicationJson)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	data := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}
	out, err = InitTxResponseFromMap(data)
	return
}

func (p *PoliClient) GetTransaction(txToken string) (out *GetTransactionResponse, err error) {
	client := &http.Client{}

	urlStr := fmt.Sprintf("%s/GetTransaction?token=%s", p.Config.ApiUrl, txToken)
	getReq, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}
	getReq.SetBasicAuth(p.Config.MerchangeCode, p.Config.MerchantAuthCode)
	getReq.Header.Add(extHttp.Accept, extHttp.ApplicationJson)
	resp, err := client.Do(getReq)
	if err != nil {
		return
	}
	data := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}

	out, err = GetTransactionResponseFromMap(data)
	return
}
