package logrus_sqs

import (
	"fmt"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSHook struct {
	Session  *sqs.SQS
	QueueUrl *string
}

func NewSQSHook(QueueName string, region string) (*SQSHook, error) {
	// Creates a SQS hook with a standard AWS session configured by AWS Environment Variables
	simple_sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})) // SQS Session

	hook, err := NewSQSHookWithSession(QueueName, simple_sess)
	if err != nil {
		return nil, err
	}
	return hook, nil
}

func NewSQSHookWithSession(queueName string, sess *session.Session) (*SQSHook, error) {
	// Creates a SQS hook with a custom AWS session
	hook := &SQSHook{}

	hook.Session = sqs.New(sess)

	resultURL, err := hook.Session.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, err
	}

	hook.QueueUrl = resultURL.QueueUrl

	return hook, nil
}

func (hook *SQSHook) Fire(entry *logrus.Entry) error {
	// Send message to SQS
	sendMessageInput := sqs.SendMessageInput{}

	sendMessageInput.QueueUrl = hook.QueueUrl
	sendMessageInput.MessageBody = &entry.Message

	// We serialize data to JSON
	data, err := json.Marshal(&entry.Data)
	if err != nil {
		return fmt.Errorf("Failed to serialize log data into JSON")
	}

	sendMessageInput.MessageAttributes = map[string]*sqs.MessageAttributeValue{
		"Level": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(entry.Level.String()),
		},
		"Time": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(entry.Time.String()),
		},
		"Data": &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(string(data)),
		},
	}
	_, err = hook.Session.SendMessage(&sendMessageInput)
	if err != nil {
		return err
	}

	return nil
}

func (hook *SQSHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
