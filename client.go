package gosqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

// NewSQSClient .
func NewSQSClient(params ListenerParams) (*sqs.SQS, error) {
	// if executing from a production env, token can be empty.
	// if executing from a corp or internet env, token is session t.
	customCreds := credentials.NewStaticCredentials(params.AWS.ID, params.AWS.Key, params.AWS.Token)
	sess, err := session.NewSession(&aws.Config{
		Credentials: customCreds,
		Region: aws.String(params.AWS.Region),
	})
	if err != nil {
		params.Logger.Errorw("Unable to start aws client.", zap.Error(err))
		return nil, err
	}
	return sqs.New(sess), nil
}
