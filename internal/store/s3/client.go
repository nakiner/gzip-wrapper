package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-kit/kit/log"
	"github.com/nakiner/gzip-wrapper/configs"
	"github.com/nakiner/gzip-wrapper/tools/logging"
	"github.com/pkg/errors"
)

type client struct {
	logger  log.Logger
	session *session.Session
	bucket  string
}

func NewClient(ctx context.Context, cfg *configs.S3) (*client, error) {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "component", "s3")
	endpoint := fmt.Sprintf("%s.compat.objectstorage.%s.oraclecloud.com", cfg.Tenancy, cfg.Region)
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(cfg.Region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, errors.Wrap(err, "make session")
	}
	return &client{
		logger:  logger,
		session: sess,
		bucket:  cfg.BucketName,
	}, nil
}
