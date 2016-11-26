package authservice

import (
    "github.com/paterson/binder/authservice/encryption"
    "encoding/json"
)

type Body            map[string]string
type EncryptedBody   string

func NewBody(str string) Body {
    var body Body
    json.Unmarshal([]byte(str), &body)
    return body
}

/* Body encryption/decryption */

func (encryptedBody EncryptedBody) Decrypt(sessionKey SessionKey) (Body, error) {
    str, err := encryption.Decrypt(string(sessionKey), string(encryptedBody))
    return NewBody(str), err
}

func (body Body) Encrypt(sessionKey SessionKey) EncryptedBody {
    cipher := encryption.Encrypt(string(sessionKey), body.toString()) // string(map[string]string)?
    return EncryptedBody(cipher) // Handle error?
}

func (body Body) toString() string {
    str, _ := json.Marshal(body)
    return string(str)
}
