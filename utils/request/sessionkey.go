package request

import (
	"log"
	"math/rand"
	"time"

	"github.com/paterson/binder/utils/encryption"
)

type SessionKey string
type EncryptedSessionKey string

func NewSessionKey() SessionKey {
	return SessionKey(randString(16))
}

func (sessionKey SessionKey) Encrypt(key string) EncryptedSessionKey {
	cipher := encryption.Encrypt(key, sessionKey.ToString())
	return EncryptedSessionKey(cipher)
}

func (sessionKey EncryptedSessionKey) Decrypt(key string) SessionKey {
	text, err := encryption.Decrypt(key, sessionKey.ToString())
	if err != nil {
		log.Fatal(err) // If this is failing, just kill because it's a systematic issue.
	}
	return SessionKey(text)
}

func (sessionKey SessionKey) ToString() string {
	return string(sessionKey)
}

func (sessionKey EncryptedSessionKey) ToString() string {
	return string(sessionKey)
}

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
