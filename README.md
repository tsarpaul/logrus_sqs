# SQS Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imguur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

### Installation
> $ go get github.com/tsarpaul/logrus_sqs

### Usage
```
package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/tsarpaul/logrus_sqs"
)

func main() {
	sqsHook, err := logrus_sqs.NewSQSHook("random_queue_name", "eu-central-1")
	if err != nil {
		panic(err)
	}
	log.AddHook(sqsHook)

	log.WithFields(log.Fields{
		"hello": "world",
	}).Info("Hello world!")
}
```

You may provide a custom AWS Session with `logrus_sqs.NewSQSHookWithSession`