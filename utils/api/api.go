package api

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/paterson/binder/utils/request"
	"os"
)

const (
	signupURL       = "http://localhost:3000/signup"
	loginURL        = "http://localhost:3000/login"
	readRequestURL  = "http://localhost:3001/request/read"
	writeRequestURL = "http://localhost:3001/request/write"
)

type AuthenticateParams struct {
	Username string
	Password string
}

type FileRequestParams struct {
	Ticket     request.EncryptedTicket
	SessionKey request.SessionKey
	Filepath   string
}

type FileParams struct {
	File     []byte
	Filename string
}

func Signup(params AuthenticateParams) map[string]string {
	resp, body, errs := gorequest.New().Post(signupURL).Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func Login(params AuthenticateParams) map[string]string {
	resp, body, errs := gorequest.New().Post(loginURL).Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestWritePermission(params FileRequestParams) map[string]string {
	resp, body, errs := gorequest.New().Post(writeRequestURL).Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestReadPermission(params FileRequestParams) map[string]string {
	resp, body, errs := gorequest.New().Post(readRequestURL).Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func WriteFile(url string, fileParams FileParams, params FileRequestParams) map[string]string {
	resp, body, errs := gorequest.New().Post(url).SendFile(fileParams.File, fileParams.Filename, "file").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func ReadFile(url string, params FileRequestParams) []byte {
	resp, body, errs := gorequest.New().Post(url).Send(params).EndBytes()
	validate(resp, errs)
	return body
}

func parseJSON(bytes []byte) map[string]string {
	var data map[string]string
	err := json.Unmarshal(bytes, &data)
	checkError(err)
	return data
}

func validate(resp gorequest.Response, errs []error) {
	if len(errs) > 0 {
		fmt.Println("Errors:", errs)
		os.Exit(1)
	} else if resp.StatusCode >= 300 {
		fmt.Println("Error: Status", resp.StatusCode)
		os.Exit(1)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
