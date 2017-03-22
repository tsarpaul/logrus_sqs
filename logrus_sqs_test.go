package logrus_sqs

import (
	"testing"
	log "github.com/Sirupsen/logrus"
)

func TestSQSHook_Fire(t *testing.T) {
	sqsHook, err := NewSQSHook("test", "eu-central-1")
	if err != nil {
		panic(err)
	}
	log.AddHook(sqsHook)

	log.WithFields(log.Fields{
		"hello": "world",
	}).Info("Hello world!")
}
