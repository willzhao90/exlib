package sms

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	extHttp "gitlab.com/sdce/exlib/http"
)

type TwilioClient struct {
	*Config
}

func NewTwilioClient(c *Config) SmsClient {
	return &TwilioClient{c}
}

//SendSMS with twilio
func (t *TwilioClient) SendSMS(to, content string) error {
	// Set account keys & information
	accountSid := t.AccountSid
	authToken := t.AuthToken
	urlStr := t.URL
	from := t.From

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", content)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add(extHttp.Accept, extHttp.ApplicationJson)
	req.Header.Add(extHttp.ContentType, extHttp.ApplicationXform)

	resp, _ := client.Do(req)
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println(data["sid"])
	} else {
		log.Println(resp.Status)
		log.Println(data["message"])
		return errors.New(data["message"].(string))
	}
	return nil

}
