package spaces

import (
  "context"
  "strings"

  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"

  "budda-fest/register/attendance/config"
)

func RegisterAttendance(cfg config.Config, emailHash, updateTime string) error {
  s3, err := mkS3Client(cfg)
  if err != nil {
    return err
  }

  if _, err := s3.PutObject(
    context.Background(),
    cfg.S3Bucket,
    emailHash,
    strings.NewReader(updateTime),
    int64(len(updateTime)),
    minio.PutObjectOptions{ContentType: "text/plain; charset=UTF-8"},
  ); err != nil {
    return err
  }

  return nil
}

func GetRegisteredAttendees(cfg config.Config) (int, error) {
  s3, err := mkS3Client(cfg)
  if err != nil {
    return 0, err
  }

  ctx, cancel := context.WithCancel(context.Background())

  defer cancel()

  objectCh := s3.ListObjects(ctx, cfg.S3Bucket, minio.ListObjectsOptions{
    Prefix:    "",
    Recursive: true,
  })

  var registrations int
  for object := range objectCh {
    if object.Err == nil {
      registrations = registrations + 1
    }
  }
  return registrations, nil
}

func mkS3Client(cfg config.Config) (*minio.Client, error) {
  endpoint := cfg.S3Endpoint
  accessKeyID := cfg.S3KeyId
  secretAccessKey := cfg.S3Secret
  useSSL := true

  // Initialize minio client object.
  return minio.New(endpoint, &minio.Options{
    Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
    Secure: useSSL,
  })
}
