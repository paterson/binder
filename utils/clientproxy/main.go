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
	params := api.AuthenticateParams{Username: username, Password: password}
	json := api.Signup(params)
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	clientProxy.Token = token
}

func (clientProxy *ClientProxy) Login(username string, password string) {
	params := api.AuthenticateParams{Username: username, Password: password}
	json := api.Login(params)
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	clientProxy.Token = token
}

func (clientProxy *ClientProxy) ReadFile(fromFilepath string, toFilepath string) {
	params := api.FileRequestParams{}
	params.Body.Ticket = clientProxy.Token.Ticket
	params.Body.Filepath = fromFilepath
	json := api.RequestReadPermission(params)
	body := api.ReadFile(json["host"], params)
	clientProxy.write(toFilepath, body)
}

func (clientProxy *ClientProxy) WriteFile(fromFilepath string, toFilepath string) {

	file, err := clientProxy.read(fromFilepath)
	checkError(err)

	params := api.FileRequestParams{}
	params.Body.Ticket = clientProxy.Token.Ticket
	params.Body.Filepath = toFilepath
	json := api.RequestWritePermission(params)

	fileParams := api.FileParams{
		File:     file,
		Filename: filepath.Base(fromFilepath),
	}
	api.WriteFile(json["host"], fileParams, params)
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
