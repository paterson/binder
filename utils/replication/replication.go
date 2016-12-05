package replication

import (
	"github.com/paterson/binder/utils/api"
	"github.com/paterson/binder/utils/request"


type Replication struct {
  file     []byte
	filepath string
	ticket   request.Ticket
}

func New(file []byte, filepath string, ticket request.Ticket) Replication {
	return Replication{
    file: file,
		filepath: filepath,
		ticket, ticket,
	}
}

// Find hosts from Directory Service.
// For each host, write file. 
func (r Replication) Replicate() {
  params := request.Params{"ticket": r.Ticket.Encrypt().SessionKey, "filepath": filepath, "noreplication": "true"}
  encryptedParams := params.Encrypt(r.Ticket.SessionKey)

  encryptedJson := api.RequestReadPermission(encryptedParams)
  json, err := encryptedJson.Decrypt(r.Ticket.SessionKey)
  checkError(err)

  hosts := strings.Split(json["hosts"], ",")
  for host := range hosts {
    hostUrl := host + "/write"

    fileParams := api.FileParams{
      File:     r.file,
      Filename: filepath.Base(r.filepath),
    }
    api.WriteFile(hostUrl, fileParams, encryptedParams)
  }
}

}
