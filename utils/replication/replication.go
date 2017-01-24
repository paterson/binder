package replication

import (
	"fmt"
	"github.com/paterson/binder/utils/api"
	"github.com/paterson/binder/utils/request"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Replication struct {
	file     multipart.File
	filepath string
	ticket   request.Ticket
}

func New(file multipart.File, filepath string, ticket request.Ticket) Replication {
	return Replication{
		file:     file,
		filepath: filepath,
		ticket:   ticket,
	}
}

// Find hosts from Directory Service.
// For each host, write file.
func (r Replication) Replicate() {
	params := request.Params{"ticket": r.ticket.Encrypt().SessionKey.ToString(), "filepath": r.filepath, "noreplication": "true"}
	encryptedParams := params.Encrypt(r.ticket.SessionKey)

	encryptedJson, err := api.RequestReadPermission(encryptedParams)
	checkError(err)
	json, err := encryptedJson.Decrypt(r.ticket.SessionKey)
	checkError(err)
	hosts := strings.Split(json["hosts"], ",")
	for _, host := range hosts {
		hostUrl := host + "/write"

		bytes, err := ioutil.ReadAll(r.file)
		checkError(err)

		fileParams := api.FileParams{
			File:     bytes,
			Filename: filepath.Base(r.filepath),
		}
		api.WriteFile(hostUrl, fileParams, encryptedParams)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
