package gosqs

import (
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

// ConsumerI .
type ConsumerI interface {
	ProcessSQSMessages()
	Run() error
}

// ConsumerS .
type ConsumerS struct {
	Params ListenerParams
}

// NewConsumer .
func NewConsumer(params ListenerParams) (*ConsumerS, error) {
	sqsClient, err := NewSQSClient(params)
	if err != nil {
		params.Logger.Errorw("Unable to instantiate SQS client: ", zap.Error(err))
		return &ConsumerS{}, err
	}
	params.SQSClient = sqsClient
	return &ConsumerS{
		Params: params,
	}, nil
}

// ProcessSQSMessages .
func (c *ConsumerS) ProcessSQSMessages() {
	errorCount := 0
	for {
		time.Sleep(1 * time.Second)
		err := c.Run()
		if err != nil {
			errorCount = errorCount + 1
			if errorCount > 10 {
				c.Params.Logger.Errorw("Unexpected error, cannot recover.", zap.Error(err))
				return
			}
			time.Sleep(60 * time.Second)
		} else {
			errorCount = 0
		}
	}
}

// Run .
func (c *ConsumerS) Run() error {
	result, err := c.Params.SQSClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl: &c.Params.QueueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(20),
		WaitTimeSeconds: aws.Int64(20),
	})
	if err != nil {
		c.Params.Logger.Errorw("error while receiving message from SQS.", zap.Error(err))
		return err
	}
	for _, message := range result.Messages {
		if err := c.Params.Handler(c, message); err != nil {
			c.Params.Logger.Errorw("Unable to handle message: ", zap.Error(err))
			return err
		}
		_, err = c.Params.SQSClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl: aws.String(c.Params.QueueURL),
			ReceiptHandle: message.ReceiptHandle,
		})
		if err != nil {
			c.Params.Logger.Errorw("Unable to delete SQS message.", zap.Error(err))
		}
	}
	return nil
}
