package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
	"io/ioutil"
	"os"
)

func main() {
	sess := session.New(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	bucket := "ses-sonofbytes.com"

	messages, err := listMessages(sess, bucket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listMessages: %s\n", err.Error())
	}

	for _, v := range messages {
		err := decryptObject(sess, bucket, v)
		if err != nil {
			aerr, ok := err.(awserr.Error)
			if ok && aerr.Code() == "NotFound" {
				fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", bucket)
			} else if ok && aerr.Code() == "InvalidWrapAlgorithmError" {
				fmt.Fprintf(os.Stderr, "unable to decrypt: %s\n", aerr.Message())
			} else {
				fmt.Fprintf(os.Stderr, "decryptObject s3://%s/%s: %s\n", bucket, v, err.Error())
			}
		}
	}

}

func decryptObject(session *session.Session, bucket, key string) (err error) {
	client := s3crypto.NewDecryptionClient(session)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := client.GetObject(input)
	// Aside from the S3 errors, here is a list of decryption client errors:
	//   * InvalidWrapAlgorithmError - returned on an unsupported Wrap algorithm
	//   * InvalidCEKAlgorithmError - returned on an unsupported CEK algorithm
	//   * V1NotSupportedError - the SDK doesn’t support v1 because security is an issue for AES ECB
	// These errors don’t necessarily mean there’s something wrong. They just tell us we couldn't decrypt some data.
	// Users can choose to log this and then continue decrypting the data that they can, or simply return the error.
	if err != nil {
		return err
	}

	// Let's read the whole body from the response
	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func listMessages(session *session.Session, bucket string) (messages []string, err error) {
	svc := s3.New(session)
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(100),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		return
	}

	for _, v := range result.Contents {
		messages = append(messages, v.String())
	}
	return
}
