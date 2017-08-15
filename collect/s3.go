package collect

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
	"io"
)

type S3SES struct {
	session *session.Session
	Bucket  string
}

func NewS3SES(bucket string) (*S3SES, error) {
	return &S3SES{
		session: session.New(&aws.Config{}),
		Bucket:  bucket,
	}, nil
}

func (s *S3SES) Get(key string) (body io.ReadCloser, err error) {
	if s == nil {
		return nil, fmt.Errorf("S3SES is nil")
	}
	if s.session == nil {
		return nil, fmt.Errorf("S3SES session is nil")
	}
	client := s3crypto.NewDecryptionClient(s.session)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	result, err := client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return result.Body, err
}

func (s *S3SES) Delete(key string) (err error) {
	client := s3.New(s.session)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	_, err = client.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3SES) Archive(key string) (err error) {
	client := s3.New(s.session)

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(s.Bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.Bucket, key)),
		Key:        aws.String(fmt.Sprintf("archive/%s", key)),
	}

	_, err = client.CopyObject(input)
	if err != nil {
		return err
	}

	err = s.Delete(key)
	if err != nil {
		return err
	}

	return nil
}
