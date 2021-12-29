package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type Repository struct {
	client  *client
	setName string
}

func NewRepository(client *client) Repository {
	return Repository{
		client: client,
	}
}

func (r *Repository) Upload(ctx context.Context, reader io.Reader, fileName string) error {
	fmt.Println("start uploader")
	uploader := s3manager.NewUploader(r.client.session)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: &r.client.bucket,
		Key:    &fileName,
		Body:   reader,
	})
	fmt.Println("done uploader")
	return err
}

func (r *Repository) Download(ctx context.Context, fileName string) (rdr io.Reader, size int64, err error) {
	svc := s3.New(r.client.session)
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &r.client.bucket,
		Key:    &fileName,
	})
	if err != nil {
		return
	}

	if result == nil {
		return
	}

	size = *result.ContentLength
	rdr = result.Body

	return
}
