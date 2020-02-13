package aws

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

// Region defines the URLs where AWS services may be accessed.
//
// See http://goo.gl/d8BP1 for more details.
type Region struct {
	Name                 string // the canonical name of this region.
	EC2Endpoint          string
	S3Endpoint           string
	S3BucketEndpoint     string // Not needed by AWS S3. Use ${bucket} for bucket name.
	S3LocationConstraint bool   // true if this region requires a LocationConstraint declaration.
	S3LowercaseBucket    bool   // true if the region requires bucket names to be lower case.
	SDBEndpoint          string
	SNSEndpoint          string
	SQSEndpoint          string
	IAMEndpoint          string
	ELBEndpoint          string
	AutoScalingEndpoint  string
	RdsEndpoint          string
	Route53Endpoint      string
}

var APSoutheast2 = Region{
	"ap-southeast-2",
	"https://ec2.ap-southeast-2.amazonaws.com",
	"https://s3-ap-southeast-2.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.ap-southeast-2.amazonaws.com",
	"https://sns.ap-southeast-2.amazonaws.com",
	"https://sqs.ap-southeast-2.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.ap-southeast-2.amazonaws.com",
	"https://autoscaling.ap-southeast-2.amazonaws.com",
	"https://rds.ap-southeast-2.amazonaws.com",
	"https://route53.amazonaws.com",
}

var Regions = map[string]Region{
	APSoutheast2.Name: APSoutheast2,
}

type Auth struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Token     string `json:"token"`
}

type AwsConfig struct {
	Region       string `json:"region"`
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	Token        string `json:"token"`
	UploadBucket string `json:"upload_bucket"`
}

// EnvAuth creates an Auth based on environment information.
// The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment
// For accounts that require a security token, it is read from AWS_SECURITY_TOKEN
// variables are used.
func EnvAuth() (auth Auth, err error) {
	auth.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	if auth.AccessKey == "" {
		auth.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	}

	auth.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if auth.SecretKey == "" {
		auth.SecretKey = os.Getenv("AWS_SECRET_KEY")
	}

	auth.Token = os.Getenv("AWS_SECURITY_TOKEN")

	if auth.AccessKey == "" {
		err = errors.New("AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY not found in environment")
	}
	if auth.SecretKey == "" {
		err = errors.New("AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY not found in environment")
	}
	return
}

// GetConfig gets auth0 config
func GetConfig(v *viper.Viper) (AwsConfig, error) {
	var c AwsConfig
	err := v.UnmarshalKey("aws", &c)
	return c, err
}
