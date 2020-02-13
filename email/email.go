package email

import "context"

type Destination struct {
	BccAddresses []*string
	CcAddresses  []*string
	ToAddresses  []*string
}

type EmailSender interface {
	SendTemplatedEmail(ctx context.Context, destination Destination, replyToAddresses []*string, source, sourceArn, template, templateArn string, templateData map[string]interface{}) error
}
