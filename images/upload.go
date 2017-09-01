package images

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/hellomd/go-sdk/config"
)

const (
	awsS3BucketName   = "hellomd-images"
	s3FileURLTemplate = "https://%v.s3.amazonaws.com/%v"
)

type Uploader interface {
	Upload(fileName string, data []byte) (fileUrl string, err error)
}

type S3Uploader struct {
	s3iface.S3API
}

func NewUploader() (Uploader, error) {
	creds := credentials.NewStaticCredentials(config.Get(AWSKeyCfgKey), config.Get(AWSSecretCfgKey), "")
	if _, err := creds.Get(); err != nil {
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion("us-west-2").WithCredentials(creds)
	return &S3Uploader{s3.New(session.New(), cfg)}, nil
}

func (u *S3Uploader) Upload(fileName string, data []byte) (string, error) {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)
	fileExtension, err := mime.ExtensionsByType(fileType)
	if err != nil || len(fileExtension) == 0 {
		return "", ErrInvalidExtension
	}

	path := config.Get(AWSBasePathCfgKey) + "/" + fileName + "/original" + fileExtension[0]
	params := &s3.PutObjectInput{
		Bucket:      aws.String(awsS3BucketName),
		Key:         aws.String(path),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	}

	if _, err = u.S3API.PutObject(params); err != nil {
		return "", err
	}

	return fmt.Sprintf(s3FileURLTemplate, awsS3BucketName, path), nil
}
