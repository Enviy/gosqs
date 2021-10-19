package gosqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

// ListenerParams .
type ListenerParams struct {
	Logger *zap.SugaredLogger
	AWS AWS
	SQSClient *sqs.SQS
	QueueURL string
	Handler func(*ConsumerS, *sqs.Message) error
}

// AWS .
type AWS struct {
	ID string
	Key string
	Token string
	Region string
}
