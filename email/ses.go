package email

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/service/ses"
	log "github.com/sirupsen/logrus"
)

type sesClient struct {
	*ses.SES
}

func NewSESSender(mySes *ses.SES) EmailSender {
	return &sesClient{mySes}
}

func (c *sesClient) SendTemplatedEmail(ctx context.Context, destination Destination, replyToAddresses []*string, source, sourceArn, template, templateArn string, templateData map[string]interface{}) error {
	buf := new(strings.Builder)
	err := json.NewEncoder(buf).Encode(templateData)
	if err != nil {
		log.Errorf("Cannot encode json data for SES message: %v", err)
		return err
	}
	templateDataStr := buf.String()
	_, err = c.SendTemplatedEmailWithContext(ctx, &ses.SendTemplatedEmailInput{
		ConfigurationSetName: nil,
		Destination: &ses.Destination{
			BccAddresses: destination.BccAddresses,
			CcAddresses:  destination.CcAddresses,
			ToAddresses:  destination.ToAddresses,
		},
		ReplyToAddresses: replyToAddresses,
		ReturnPath:       nil,
		ReturnPathArn:    nil,
		Source:           &source,
		SourceArn:        &sourceArn,
		Tags:             []*ses.MessageTag{},
		Template:         &template,
		TemplateArn:      &templateArn,
		TemplateData:     &templateDataStr,
	})
	if err != nil {
		return err
	}
	return nil
}
