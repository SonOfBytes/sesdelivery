package sesdelivery

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
	"io"
)

type S3Decryptor struct {
	session *session.Session
	Bucket  string
}

func NewS3Decryptor(bucket string) (*S3Decryptor, error) {
	return &S3Decryptor{
		session: session.New(&aws.Config{
			Region: aws.String("eu-west-1"),
		}),
		Bucket: bucket,
	}, nil
}

func (s *S3Decryptor) decryptObject(key string) (body io.ReadCloser, err error) {
	client := s3crypto.NewDecryptionClient(session)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	result, err := client.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("GetObject: %s", err.Error())
	}

	return result.Body, err
}
