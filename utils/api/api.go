package api

import (
  "github.com/parnurzeal/gorequest"
  "encoding/json"
  "os"
)

const (
  signupURL       = "http://localhost:3000/signup"
  loginURL        = "http://localhost:3000/login"
  readRequestURL  = "http://localhost:3002/request/read"
  writeRequestURL = "http://localhost:3002/request/write"
)

type AuthenticateParams struct {
  username string
  password string
}

type FileRequestParams struct {
  ticket     request.EncryptedTicket
  sessionKey request.SessionKey
  filepath   string
}


func Signup(params AuthenticateParams) map[string]string {
  resp, body, errs := gorequest.New().Post(signupURL).Send(params).EndBytes()
  validate(resp, errs)
  return parseJSON(body)
}

func Login(params AuthenticateParams) map[string]string {
  resp, body, errs := gorequest.New().Post(loginURL).Send(paramStruct).EndBytes()
  validate(resp, errs)
  return parseJSON(body)
}

func WriteFileRequest(params FileRequestParams) map[string]string {
  resp, body, errs := gorequest.New().Post(writeRequestURL).Send(params).EndBytes()
  validate(resp, errs)
  return parseJSON(body)
}

func ReadFileRequest(params FileRequestParams) map[string]string {
  resp, body, errs := gorequest.New().Post(readRequestURL).Send(params).EndBytes()
  validate(resp, errs)
  return parseJSON(body)
}

func parseJSON(bytes []byte) map[string]string {
  var data map[string]string
  err := json.Unmarshal([]byte(str), &data)
  checkError(err)
  return data
}

func validate(resp gorequest.Response, errs []error) {
  if len(errs) > 0 {
    fmt.Println("Errors:", errs)
    os.Exit(1)
  } else if resp.StatusCode >= 300 {
    fmt.Println("Error: Status" resp.StatusCode)
    os.Exit(1)
  }
}

func checkError(err error) {
  if err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
  }
}