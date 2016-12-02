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

type FileParams struct {
	File     []byte
	Filename string
}

func Signup(params request.Params) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", signupURL, params))
	resp, body, errs := gorequest.New().Post(signupURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func Login(params request.Params) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", loginURL, params))
	resp, body, errs := gorequest.New().Post(loginURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestWritePermission(params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", writeRequestURL, params))
	resp, body, errs := gorequest.New().Post(writeRequestURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestReadPermission(params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", readRequestURL, params))
	resp, body, errs := gorequest.New().Post(readRequestURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func WriteFile(url string, fileParams FileParams, params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", url, params))
	resp, body, errs := gorequest.New().Post(url).Type("multipart").SendFile(fileParams.File, fileParams.Filename, "file").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func ReadFile(url string, params request.EncryptedParams) []byte {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", url, params))
	resp, body, errs := gorequest.New().Post(url).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return body
}

func parseJSON(bytes []byte) request.EncryptedParams {
	var data request.EncryptedParams
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
