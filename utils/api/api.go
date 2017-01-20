package api

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/paterson/binder/utils/constants"
	"github.com/paterson/binder/utils/request"
	"os"
)

type FileParams struct {
	File     []byte
	Filename string
}

func Signup(params request.Params) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.SignupURL, params))
	resp, body, errs := gorequest.New().Post(constants.SignupURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func Login(params request.Params) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.LoginURL, params))
	resp, body, errs := gorequest.New().Post(constants.LoginURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestWritePermission(params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.WriteRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.WriteRequestURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestReadPermission(params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.ReadRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.ReadRequestURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestLock(params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.LockRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.LockRequestURL).Type("form").Send(params).EndBytes()
	validate(resp, errs)
	return parseJSON(body)
}

func RequestUnlock(params request.EncryptedParams) request.EncryptedParams {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.UnlockRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.UnlockRequestURL).Type("form").Send(params).EndBytes()
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
