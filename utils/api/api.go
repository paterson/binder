package api

import (
	"encoding/json"
	"errors"
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

func Signup(params request.Params) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.SignupURL, params))
	resp, body, errs := gorequest.New().Post(constants.SignupURL).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func Login(params request.Params) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.LoginURL, params))
	resp, body, errs := gorequest.New().Post(constants.LoginURL).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func RequestWritePermission(params request.EncryptedParams) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.WriteRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.WriteRequestURL).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func RequestReadPermission(params request.EncryptedParams) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.ReadRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.ReadRequestURL).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func RequestLock(params request.EncryptedParams) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.LockRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.LockRequestURL).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func RequestUnlock(params request.EncryptedParams) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", constants.UnlockRequestURL, params))
	resp, body, errs := gorequest.New().Post(constants.UnlockRequestURL).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func WriteFile(url string, fileParams FileParams, params request.EncryptedParams) (request.EncryptedParams, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", url, params))
	resp, body, errs := gorequest.New().Post(url).Type("multipart").SendFile(fileParams.File, fileParams.Filename, "file").Send(params).EndBytes()
	err := validate(resp, errs)
	return parseJSON(body), err
}

func ReadFile(url string, params request.EncryptedParams) ([]byte, error) {
	fmt.Println(fmt.Sprintf("Posting to url %s with params %+v:", url, params))
	resp, body, errs := gorequest.New().Post(url).Type("form").Send(params).EndBytes()
	err := validate(resp, errs)
	return body, err
}

func parseJSON(bytes []byte) request.EncryptedParams {
	var data request.EncryptedParams
	err := json.Unmarshal(bytes, &data)
	checkError(err)
	return data
}

func validate(resp gorequest.Response, errs []error) error {
	if len(errs) > 0 {
		return errs[0]
	} else if resp.StatusCode >= 300 {
		return errors.New("Something went Wrong")
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
