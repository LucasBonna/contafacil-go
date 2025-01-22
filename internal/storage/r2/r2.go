package r2

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	configuration "github.com/lucasbonna/contafacil_api/internal/config"
)

type R2 struct {
  client *s3.Client
}

func NewR2Client(accessKeyId string, accessKeySecret string, region string, accountId string) (*R2, error) {
  cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
    config.WithRegion("auto"),
  )
  if err != nil {
    return nil, err
  }

  client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
  })

  return &R2{
    client: client,
  }, nil
}

func (r2 R2) Upload(file io.Reader, fileId uuid.UUID) error {
  _, err := r2.client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(configuration.StorageBucketName),
        Key: aws.String(fileId.String()),
        Body: file,
  })
  if err != nil {
    return err
  }

  log.Printf("File %s uploaded sucessfully", fileId)
  return nil
}

func (r2 R2) Download(fileId uuid.UUID) ([]byte, error) {
  result, err := r2.client.GetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(configuration.StorageBucketName),
        Key: aws.String(fileId.String()),
  })
  if err != nil {
    return nil, err
  }
  defer result.Body.Close()
  
  return io.ReadAll(result.Body)
}
