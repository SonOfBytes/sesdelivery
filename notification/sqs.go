package notification

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"encoding/json"
)

type SQSNotice struct {
	ReceiptHandle *string `json:"-"`
	NotificationType string `json:"notificationType"`
	Mail struct {
		Source string `json:"source"`
	} `json:"mail"`
	Receipt struct {
		Recipients []string `json:"recipients"`
		Action struct {
			BucketName string `json:"bucketName"`
			ObjectKey  string `json:"objectKey"`
		} `json:"action"`
	} `json:"receipt"`
}

type SQSNotifier struct {
	session *session.Session
	service *sqs.SQS
	queue string
}

func NewSQSNotifier(queue string) (*SQSNotifier, error) {
	sess := session.New(&aws.Config{
		MaxRetries:  aws.Int(5),
	})

	svc := sqs.New(sess)

	return &SQSNotifier{
		session: sess,
		service: svc,
		queue: queue,
	}, nil
}

func (s *SQSNotifier) Get() (notice *SQSNotice) {
	// Receive message
	receive_params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queue),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(60), // this is twice the event poll length
		WaitTimeSeconds:     aws.Int64(20),
	}

	receive_resp, err := s.service.ReceiveMessage(receive_params)
	if err != nil {
		log.Printf("Error ReceiveMessage: %s\n", err)
		return
	}

	if len(receive_resp.Messages) == 0 {
		return
	}

	notice = &SQSNotice{}
	err = json.Unmarshal([]byte(*receive_resp.Messages[0].Body), notice)
	if err != nil {
		log.Printf("Error Unmarshal: %s\n", err)
		return nil
	}

	notice.ReceiptHandle = receive_resp.Messages[0].ReceiptHandle

	return notice
}

func (s *SQSNotifier) Delete(notice *SQSNotice) (err error) {
	delete_params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queue),
		ReceiptHandle: notice.ReceiptHandle,

	}
	_, err = s.service.DeleteMessage(delete_params) // No response returned when successful.
	if err != nil {
		return err
	}
	return
}