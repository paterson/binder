package clientproxy

import (
	"fmt"
	"github.com/paterson/binder/utils/api"
	"github.com/paterson/binder/utils/request"
	"io/ioutil"
	"os"
)

type ClientProxy struct {
	Token request.Token
}

type AuthenticateParams struct {
	username string
	password string
}

type FileParams struct {
	ticket     request.EncryptedTicket
	sessionKey request.SessionKey
	filepath   string
}

func New() *ClientProxy {
	return &ClientProxy{}
}

func (clientProxy *ClientProxy) Signup(username string, password string) {
	params := AuthenticateParams{username: username, password: password}
	json := api.Signup(params)
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	clientProxy.Token = token
}

func (clientProxy *ClientProxy) Login(username string, password string) {
	params := AuthenticateParams{username: username, password: password}
	json := api.Login(params)
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	clientProxy.Token = token
}

func (clientProxy *ClientProxy) ReadFile(fromFilepath string, toFilepath string) {
	params := FileParams{
		ticket:     clientProxy.Token.Ticket,
		sessionKey: clientProxy.Token.SessionKey,
		filepath:   fromFilepath,
	}
	json := api.ReadFileRequest(params)
	resp, body, errs = clientProxy.http.Post(json["host"]).Send(params).End()
	checkErrors(errs)
	clientProxy.write(toFilepath, body)
}

func (clientProxy *ClientProxy) WriteFile(fromFilepath string, toFilepath string) {
	params := FileParams{
		ticket:     clientProxy.Token.Ticket,
		sessionKey: clientProxy.Token.SessionKey,
		filepath:   toFilepath,
	}
	resp, body, errs := clientProxy.http.Post(writeRequestURL).Send(params).End()
	checkErrors(errs)
	data, err := clientProxy.parseJSON(body)
	checkError(err)
	file, err := clientProxy.read(fromFilepath)
	checkError(err)
	// Add file as uploaded param "file" + Ticket and Sessionkey + path
	resp, body, errs = clientProxy.http.Post(data["host"]).Send(params).End()
	checkErrors(errs)
}

func (ClientProxy *ClientProxy) read(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func (ClientProxy *ClientProxy) write(filepath string, data string) error {
	return ioutil.WriteFile(filepath, []byte(data), 0644)
}

func (clientProxy *ClientProxy) parseJSON(str string) map[string]string {
	var data map[string]string
	err := json.Unmarshal([]byte(str), &data)
	checkError(err)
	return data
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func checkErrors(errs []error) {
	if len(errs) > 0 {
		fmt.Println("Errors:", errs)
		os.Exit(1)
	}
}
