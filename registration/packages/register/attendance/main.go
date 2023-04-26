package main

import (
  "crypto/sha256"
  "fmt"
  "io"
  "log"
  "time"

  "budda-fest/register/attendance/config"
  "budda-fest/register/attendance/spaces"
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

func Main(req Request) (*Response, error) {
  cfg := config.FromEnv()

  switch method := req.HttpDetails.Method; method {
  case "GET":
    attendees, err := spaces.GetRegisteredAttendees(cfg)
    if err != nil {
      log.Fatal(err)
      return nil, err
    }
    return respond(fmt.Sprintf("%d", attendees))

  case "POST":
    if req.Email == "" {
      log.Fatal("Please provide an email address")
      return redirectTo(
        fmt.Sprintf("%s/registration", cfg.UIURL),
        "Please provide an email address",
      )
    }

    emailHash := hash256(req.Email)
    updateTime := fmt.Sprintf("%d", time.Now().UnixNano())

    if err := spaces.RegisterAttendance(cfg, emailHash, updateTime); err != nil {
      log.Fatal(err)
      return redirectTo(
        fmt.Sprintf("%v/something-went-wrong", cfg.UIURL),
        fmt.Sprintf("%+v", err),
      )
    }

    return redirectTo(
      fmt.Sprintf("%s/thankyou", cfg.UIURL),
      "Thank You!",
    )

  default:
    return &Response{}, nil
  }
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
