package kms

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
)

// KmsEncrypt encrypts byte array to byte array
func (k *KMS) KmsEncrypt(b []byte) (*[]byte, error) {
	// Encrypt the data
	result, err := k.KMS.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(k.Config.KmsKey),
		Plaintext: b,
	})

	if err != nil {
		return nil, err
	}

	return &result.CiphertextBlob, nil
}

// KmsDecrypt decrypts byte array to byte array
func (k *KMS) KmsDecrypt(blob []byte) (*[]byte, error) {
	result, err := k.KMS.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})

	if err != nil {
		return nil, err
	}

	return &result.Plaintext, nil
}

// KmsEncryptString encrypts string to string
// Used for eth mnemonic and eos private key
func (k *KMS) KmsEncryptString(s string) string {
	bytes, err := k.KmsEncrypt([]byte(s))
	if err != nil {
		log.Fatal(err)
	}

	encryptedString := hexutil.Encode(*bytes)
	if err != nil {
		log.Fatal(err)
	}

	return encryptedString
}

// KmsDecryptString decrypts string to string
// Used for eth mnemonic and eos private key
func (k *KMS) KmsDecryptString(encryptedString string) string {
	encryptedBytes, err := hexutil.Decode(encryptedString)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := k.KmsDecrypt(encryptedBytes)
	if err != nil {
		log.Fatal(err)
	}
	s := string(*bytes)

	return s
}
