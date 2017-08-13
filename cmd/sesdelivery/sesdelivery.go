package main

import (
	"github.com/sonofbytes/sesdelivery/notification"
	"log"
	"github.com/sonofbytes/sesdelivery/collect"
	"github.com/sonofbytes/sesdelivery/deliver"
	"fmt"
	"flag"
)

var (
	smtpServer string
	sqsQueue string
)

func init() {
	flag.StringVar(&smtpServer, "s", "", "smtp server host to send email to")
	flag.StringVar(&sqsQueue, "q", "", "SQS queue to get notices from")
}
func main() {
	flag.Parse()
	log.Printf("Send email to %s from notices at %s", smtpServer, sqsQueue)

	collectors := make(map[string]*collect.S3SES)
	delivery, err := deliver.NewSMTP(smtpServer)
	if err != nil {
		log.Fatal(err.Error())
	}

	notifier, err := notification.NewSQSNotifier(sqsQueue)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Polling ")
	for {
		notice := notifier.Get()
		if notice == nil {
			fmt.Printf(".")
			continue
		}

		bucket := notice.Receipt.Action.BucketName
		key := notice.Receipt.Action.ObjectKey
		recipients := notice.Receipt.Recipients
		sender := notice.Mail.Source

		if bucket != "" && key != "" && len(recipients) > 0 {
			var ok bool
			if _, ok = collectors[bucket]; !ok {
				collectors[bucket], err = collect.NewS3SES(bucket)
				if err != nil {
					log.Fatal(err.Error())
				}
			}

			body, err := collectors[bucket].Get(key)
			if err != nil {
				log.Fatal(err.Error())
			}

			// Send the message
			err = delivery.Send(sender, recipients, body)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = body.Close()
			if err != nil {
				log.Fatal(err.Error())
			}

			// Delete the notice so no resend
			err = notifier.Delete(notice)
			if err != nil {
				log.Fatal(err.Error())
			}

			// Archive the original message
			err = collectors[bucket].Archive(key)
			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Printf(">")
		}
	}
}
