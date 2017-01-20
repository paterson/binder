package clientproxy

import (
	"fmt"
	"github.com/paterson/binder/utils/api"
	"github.com/paterson/binder/utils/cache"
	"github.com/paterson/binder/utils/request"
	"io/ioutil"
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
	json, err := api.Signup(params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
	clientProxy.Token = token
	fmt.Println(fmt.Sprintf("Received Data %+v", json))
}

func (clientProxy *ClientProxy) Login(username string, password string) {
	params := request.Params{"username": username, "password": password}
	json, err := api.Login(params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	encryptedToken := request.TokenFromJSON(json)
	token := encryptedToken.Decrypt(password)
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

	encryptedJson, err := api.RequestReadPermission(encryptedParams)
	checkError(err)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	// Try each host until one succeeds (So if one is down, system is still fine) (replication)
	hosts := strings.Split(json["hosts"], ",")
	i := 0
	successful := false
	for !successful {
		host := hosts[i]
		hostUrl := host + "/read"
		body, err := api.ReadFile(hostUrl, encryptedParams)
		if err == nil {
			clientProxy.write(toFilepath, body)
			successful = true
		}
		i += 1
	}
	fmt.Println("Read File and wrote to local file path")
}

/* Use Session key to identify client locking file */
func (clientProxy *ClientProxy) LockFile(filepath string) bool {
	sessionKey := string(clientProxy.Token.Ticket.SessionKey)
	params := request.Params{"ticket": sessionKey, "filepath": filepath}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson, e := api.RequestLock(encryptedParams)
	_, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)

	return e == nil && err == nil
}

func (clientProxy *ClientProxy) UnlockFile(filepath string) bool {
	sessionKey := string(clientProxy.Token.Ticket.SessionKey)
	params := request.Params{"ticket": sessionKey, "filepath": filepath, "lock_key": sessionKey}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson, e := api.RequestUnlock(encryptedParams)
	_, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)

	return e == nil && err == nil

}

func (clientProxy *ClientProxy) WriteFile(fromFilepath string, toFilepath string) {
	file, err := clientProxy.read(fromFilepath)
	checkError(err)

	params := request.Params{"ticket": string(clientProxy.Token.Ticket.SessionKey), "filepath": toFilepath}
	encryptedParams := params.Encrypt(clientProxy.Token.SessionKey)
	encryptedJson, err := api.RequestWritePermission(encryptedParams)
	checkError(err)
	json, err := encryptedJson.Decrypt(clientProxy.Token.SessionKey)
	checkError(err)

	// Try each host until one succeeds (So if one is down, system is still fine) (replication)
	hosts := strings.Split(json["hosts"], ",")
	i := 0
	successful := false
	for !successful {
		host := hosts[i]
		hostUrl := host + "/write"
		fileParams := api.FileParams{
			File:     file,
			Filename: filepath.Base(fromFilepath),
		}
		_, err := api.WriteFile(hostUrl, fileParams, encryptedParams)
		if err == nil {
			fmt.Println("Wrote File to fileserver")
			successful = true
		}
		i++
	}
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
