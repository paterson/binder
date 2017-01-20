package clientproxy

import (
	"fmt"
	"github.com/paterson/binder/utils/api"
	"github.com/paterson/binder/utils/cache"
	"github.com/paterson/binder/utils/request"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type ClientProxy struct {
	Token request.Token
	cache cache.Cache
}

func New() ClientProxy {
	return ClientProxy{
		cache: cache.New(2), // 2MB in memory cache
	}
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
	data, err := clientProxy.cache.Get(fromFilepath)
	if err == nil {
		clientProxy.write(toFilepath, data)
		return
	}

	params := request.Params{"ticket": string(clientProxy.Token.Ticket.SessionKey), "filepath": fromFilepath}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)

	encryptedJson := api.RequestReadPermission(encryptedParams)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	// Pick a random host from the hosts (file servers) to read from (replication)
	hosts := strings.Split(json["hosts"], ",")
	host := hosts[rand.Intn(len(hosts))]
	hostUrl := host + "/read"

	body := api.ReadFile(hostUrl, encryptedParams)
	clientProxy.write(toFilepath, body)
	fmt.Println("Read File and wrote to local file path")
}

/* Use Session key to identify client locking file */
func (clientProxy *ClientProxy) LockFile(filepath string) bool {
	sessionKey := string(clientProxy.Token.Ticket.SessionKey)
	params := request.Params{"ticket": sessionKey, "filepath": filepath}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson := api.RequestLock(encryptedParams)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	return json["error"] == ""
}

func (clientProxy *ClientProxy) UnlockFile(filepath string) bool {
	sessionKey := string(clientProxy.Token.Ticket.SessionKey)
	params := request.Params{"ticket": sessionKey, "filepath": filepath, "lock_key": sessionKey}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson := api.RequestUnlock(encryptedParams)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	return json["error"] == ""

}

func (clientProxy *ClientProxy) WriteFile(fromFilepath string, toFilepath string) {
	file, err := clientProxy.read(fromFilepath)
	checkError(err)

	params := request.Params{"ticket": string(clientProxy.Token.Ticket.SessionKey), "filepath": toFilepath}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson := api.RequestWritePermission(encryptedParams)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	// Pick a random host of file servers to write (replication)
	hosts := strings.Split(json["hosts"], ",")
	host := hosts[rand.Intn(len(hosts))]
	hostUrl := host + "/write"

	fileParams := api.FileParams{
		File:     file,
		Filename: filepath.Base(fromFilepath),
	}
	api.WriteFile(hostUrl, fileParams, encryptedParams)
	fmt.Println("Wrote File to fileserver")
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
