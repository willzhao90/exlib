package kms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMS struct {
	KMS    *kms.KMS
	Config *Config
}

func New(conf *Config, cred ...*credentials.Credentials) (c *KMS, err error) {
	var credential *credentials.Credentials
	if cred != nil && cred[0] != nil {
		credential = cred[0]
	} else {
		credential = credentials.NewEnvCredentials()
	}

	// check if credentials has been found
	_, err = credential.Get()
	if err != nil {
		return nil, fmt.Errorf("Env credential not found: %s", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credential,
		Region:      aws.String(conf.Region)},
	)
	if err != nil {
		return
	}

	svc := kms.New(sess)

	c = &KMS{
		KMS:    svc,
		Config: conf,
	}
	return
}
