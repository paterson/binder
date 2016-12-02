package clientproxy

import (
	"fmt"
	"github.com/paterson/binder/utils/api"
	"github.com/paterson/binder/utils/request"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ClientProxy struct {
	Token request.Token
}

func New() ClientProxy {
	return ClientProxy{}
}

func (clientProxy *ClientProxy) Signup(username string, password string) {
	params := request.Params{"username": username, "password": password}
	json := api.Signup(params)
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	clientProxy.Token = token
	fmt.Println(fmt.Sprintf("Received Data %+v", json))
}

func (clientProxy *ClientProxy) Login(username string, password string) {
	params := request.Params{"username": username, "password": password}
	fmt.Println(fmt.Sprintf("Sent Params %+v", params))
	json := api.Login(params)
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	fmt.Println("Token:", token)
	clientProxy.Token = token
	fmt.Println(fmt.Sprintf("Received Data %+v", json))
}

func (clientProxy *ClientProxy) ReadFile(fromFilepath string, toFilepath string) {
	params := request.Params{"ticket": string(clientProxy.Token.Ticket.SessionKey), "filepath": fromFilepath}
	fmt.Println(fmt.Sprintf("Sent Params %+v", params))
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)

	encryptedJson := api.RequestReadPermission(encryptedParams)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	body := api.ReadFile(json["host"], encryptedParams)
	clientProxy.write(toFilepath, body)
	fmt.Println(fmt.Sprintf("Received Data %+v", json))
}

func (clientProxy *ClientProxy) WriteFile(fromFilepath string, toFilepath string) {
	file, err := clientProxy.read(fromFilepath)
	checkError(err)

	params := request.Params{"ticket": string(clientProxy.Token.Ticket.SessionKey), "filepath": toFilepath}
	fmt.Println(fmt.Sprintf("Sent Params 1 %+v", params))
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson := api.RequestWritePermission(encryptedParams)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	fileParams := api.FileParams{
		File:     file,
		Filename: filepath.Base(fromFilepath),
	}
	api.WriteFile(json["host"], fileParams, encryptedParams)
	fmt.Println(fmt.Sprintf("Received Data %+v", json))
}

func (ClientProxy *ClientProxy) read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (ClientProxy *ClientProxy) write(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
