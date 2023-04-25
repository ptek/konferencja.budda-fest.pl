package main

import (
  "context"
  "crypto/sha256"
  "fmt"
  "io"
  "log"
  "strings"
  "time"

  "github.com/caarlos0/env/v8"
  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"
)

type Http struct {
  Headers map[string]string `json:"headers"`
  Method  string            `json:"method"`
  Path    string            `json:"path"`
}

type Request struct {
  HttpDetails Http   `json:"http"`
  Email       string `json:"email"`
}

type Response struct {
  StatusCode int               `json:"statusCode,omitempty"`
  Headers    map[string]string `json:"headers,omitempty"`
  Body       string            `json:"body,omitempty"`
}

func GetRegisteredAttendees(cfg config) (int64, error) {
  s3, err := s3client(cfg)
  if err != nil {
    return 0, err
  }

  ctx, cancel := context.WithCancel(context.Background())

  defer cancel()

  objectCh := s3.ListObjects(ctx, cfg.S3Bucket, minio.ListObjectsOptions{
    Prefix:    "",
    Recursive: true,
  })

  var registrations int64
  for object := range objectCh {
    if object.Err == nil {
      registrations = registrations + 1
    }
  }
  return registrations, nil
}

func RegisterAttendance(cfg config, email string) error {
  s3, err := s3client(cfg)
  if err != nil {
    return err
  }

  filename := hash256(email)
  updateTime := fmt.Sprintf("%d", time.Now().UnixNano())

  if _, err := s3.PutObject(
    context.Background(),
    cfg.S3Bucket,
    filename,
    strings.NewReader(updateTime),
    19,
    minio.PutObjectOptions{ContentType: "text/plain; charset=UTF-8"},
  ); err != nil {
    return err
  }

  return nil
}

func Main(req Request) (*Response, error) {
  configOpts := env.Options{
    Prefix:          "BUDDAFEST_REGISTRATION_",
    RequiredIfNoDef: true,
  }

  var cfg config

  if err := env.ParseWithOptions(&cfg, configOpts); err != nil {
    log.Fatal(err)
  }

  switch method := req.HttpDetails.Method; method {
  case "GET":
    attendees, err := GetRegisteredAttendees(cfg)
    if err != nil {
      log.Fatal(err)
      return nil, err
    }
    return respond(fmt.Sprintf("%d", attendees))

  case "POST":
    if req.Email == "" {
      return redirectTo(
        fmt.Sprintf("%s/registration", cfg.UIURL),
        "Please provide an email address",
      )
      log.Fatal("Please provide an email address")
      return nil, fmt.Errorf("Please provide an email address")
    }

    if err := RegisterAttendance(cfg, req.Email); err != nil {
      return redirectTo(
        fmt.Sprintf("%v/something-went-wrong", cfg.UIURL),
        fmt.Sprintf("%+v", err),
      )
      log.Fatal(err)
      return nil, err
    }

    return redirectTo(
      fmt.Sprintf("%s/thankyou", cfg.UIURL),
      "Thank You!",
    )

  default:
    return &Response{}, nil
  }
}

type config struct {
  UIURL      string `env:"UI_URL"`
  S3Endpoint string `env:"S3_ENDPOINT"`
  S3Bucket   string `env:"S3_BUCKET"`
  S3KeyId    string `env:"S3_ID"`
  S3Secret   string `env:"S3_SECRET"`
}

func s3client(cfg config) (*minio.Client, error) {
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

func hash256(e string) string {
  h := sha256.New()
  io.WriteString(h, e)
  hashedString := h.Sum(nil)
  return fmt.Sprintf("%x", hashedString)
}

func redirectTo(url string, body string) (*Response, error) {
  return &Response{
    StatusCode: 302,
    Headers: map[string]string{
      "Location": url,
    },
    Body: body,
  }, nil
}

func respond(body string) (*Response, error) {
  return &Response{Body: body}, nil
}
