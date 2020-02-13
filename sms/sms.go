package sms

import (
	"bytes"
	"text/template"
)

type SmsClient interface {
	SendSMS(string, string) error
}

//MessageContent applies data to the template with text/template
func MessageContent(contentTemplate string, data map[string]string) (string, error) {
	tmpl, err := template.New("smsTemplate").Parse(contentTemplate)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	result := tpl.String()
	return result, nil
}
