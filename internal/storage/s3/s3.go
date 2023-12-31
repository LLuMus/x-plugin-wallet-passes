package s3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	AwsBucket string
}

func NewS3Storage(awsBucket string) *S3Storage {
	return &S3Storage{
		AwsBucket: awsBucket,
	}
}

func (s *S3Storage) Save(key, filePath, contentType string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", nil
	}

	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	return s.SaveFile(key, f, contentType, fileInfo.Size())
}

func (s *S3Storage) SaveFile(key string, originalReader io.Reader, contentType string, contentLength int64) (string, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return "", err
	}

	var b = make([]byte, contentLength)

	_, err = originalReader.Read(b)
	if err != nil {
		return "", err
	}

	if contentLength <= 0 {
		return "", fmt.Errorf("empty bytes file buffer %s", key)
	}

	if contentType == "" {
		contentType = http.DetectContentType(b)
	}

	var (
		uploader    = s3.New(awsSession)
		uploadInput = &s3.PutObjectInput{
			Bucket:        aws.String(s.AwsBucket),
			Key:           aws.String(key),
			Body:          bytes.NewReader(b),
			ContentType:   aws.String(contentType),
			ContentLength: aws.Int64(contentLength),
		}
	)

	_, err = uploader.PutObject(uploadInput)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", key), nil
}

func (s *S3Storage) GenerateSignedUrl(key string, expires time.Duration) (string, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return "", err
	}

	svc := s3.New(awsSession)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.AwsBucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(expires)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
