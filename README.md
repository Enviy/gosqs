# GOSQS


This package is intended to be a Go SQS listener that can be imported and invoked as a goroutine handled by the life cycle of your service. It's designed for plug and play for any SQS consumption use case. (Service life cycle management not required, can also be used with traditional channels.)


  Requirements:
  - AWS credentials for a user/role with appropriate permissions. Populated in ListenerParams.AWS.
  - A handler method you provide for processing the SQS messages. Handler methods must follow this pattern:
	```
		func <YourHandlerName>(c *gosqs.ConsumerS, m *sqs.Messages) error
	```

### Overview and Notes
Specify what SQS queue to monitor as well as your handling method in the ListenerParams struct. Review the structs in the Model file.
```
type ListenerParams struct {
	Logger *zap.SugaredLogger
	AWS AWS			// (refer to model file)
	SQSClient *sqs.SQS	// (instantiated for you)
	QueueURL string		// (queue to monitor)
	Handler func(*ConsumerS, *sqs.Message) error // (your handler method)
}
```

The general config.Provider was considered but that would have either required additional configuration on the user end. Or it would require assumptions on my part regarding what the credentials would be named on the user side which didn't sound appealing. So it was not used and instead the user is responsible for providing the appropriate values to the ListenerParams when instantiating a new consumer.

Because the handler method is abstracted away from the listener (to support any use case), the handler also needs to be stand alone. Each handler is responsible for instantiating any additional clients that it might need. It's also responsible for the collection of any additinally required credentials your handler workflow might need.because config.

