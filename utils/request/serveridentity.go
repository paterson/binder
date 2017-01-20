package request

import (
	"log"

	"github.com/paterson/binder/utils/encryption"
)

type ServerIdentity string
type EncryptedServerIdentity string

var AuthServerIdentity = ServerIdentity("AUTHSERVER")

func (serverIdentity ServerIdentity) Encrypt(key string) EncryptedServerIdentity {
	cipher := encryption.Encrypt(key, serverIdentity.toString())
	return EncryptedServerIdentity(cipher)
}

func (serverIdentity EncryptedServerIdentity) Decrypt(key string) ServerIdentity {
	text, err := encryption.Decrypt(key, serverIdentity.toString())
	if err != nil {
		log.Fatal(err)
	}
	return ServerIdentity(text)
}

func (serverIdentity ServerIdentity) toString() string {
	return string(serverIdentity)
}

func (serverIdentity EncryptedServerIdentity) toString() string {
	return string(serverIdentity)
}
