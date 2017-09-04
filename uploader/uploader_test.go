package uploader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testImageName = "testImageName"
	testBasePath  = "testBasePath"
)

func TestS3Uploader(t *testing.T) {
	Convey("Given an AWS S3 file Uploader", t, func() {
		uploader, storage := setup()
		Convey("When I try to upload an image", func() {
			fileURL, err := uploader.Upload(testImageName, []byte(testImage))
			if err != nil {
				t.Fatal(err)
			}

			key := fmt.Sprintf(s3filePathTemplate, testBasePath, testImageName, ".asc")
			So(storage[key], ShouldResemble, []byte(testImage))
			So(fileURL, ShouldResemble, fmt.Sprintf(s3FileURLTemplate, "", key))
		})
	})
}

func setup() (Uploader, map[string][]byte) {
	storage := map[string][]byte{}
	return &S3Uploader{testBasePath, &FakeS3{storage: storage}}, storage
}

type FakeS3 struct {
	s3iface.S3API
	storage map[string][]byte
}

// PutObject -
func (c *FakeS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if input.Key == nil {
		return nil, errors.New("empty Key")
	}

	if input.Body == nil {
		return nil, errors.New("empty Body")
	}

	if input.ContentType == nil {
		return nil, errors.New("empty ContentType")
	}

	b, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}

	c.storage[*input.Key] = b
	return nil, nil
}
