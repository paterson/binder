package request

import (
	"github.com/gin-gonic/gin"
)

/* Token. Key: User password */
type Token struct {
	Ticket         EncryptedTicket
	SessionKey     SessionKey
	ServerIdentity ServerIdentity
	Timeout        TimeoutPeriod
}

type EncryptedToken struct {
	Ticket         EncryptedTicket
	SessionKey     EncryptedSessionKey
	ServerIdentity EncryptedServerIdentity
	Timeout        EncryptedTimeoutPeriod
}

func GenerateToken() Token {
	sessionKey := NewSessionKey()
	ticket := Ticket{SessionKey: sessionKey}
	return Token{
		Ticket:         ticket.Encrypt(),
		SessionKey:     sessionKey,
		ServerIdentity: AuthServerIdentity,
		Timeout:        DefaultTimeout,
	}
}

func TokenFromJSON(json map[string]string) EncryptedToken {
	return EncryptedToken{
		Ticket:         EncryptedTicket{SessionKey: EncryptedSessionKey(json["ticket"])},
		SessionKey:     EncryptedSessionKey(json["session_key"]),
		ServerIdentity: EncryptedServerIdentity(json["server_identity"]),
		Timeout:        EncryptedTimeoutPeriod(json["timeout"]),
	}
}

func (encryptedToken EncryptedToken) ToJSON() gin.H {
	return gin.H{
		"ticket":          encryptedToken.Ticket.SessionKey,
		"session_key":     encryptedToken.SessionKey,
		"server_identity": encryptedToken.ServerIdentity,
		"timeout":         encryptedToken.Timeout,
	}
}

func (encryptedToken EncryptedToken) Decrypt(password string) Token {
	return Token{
		Ticket:         encryptedToken.Ticket,
		SessionKey:     encryptedToken.SessionKey.Decrypt(password),
		ServerIdentity: encryptedToken.ServerIdentity.Decrypt(password),
		Timeout:        encryptedToken.Timeout.Decrypt(password),
	}
}

func (token Token) Encrypt(password string) EncryptedToken {
	return EncryptedToken{
		Ticket:         token.Ticket,
		SessionKey:     token.SessionKey.Encrypt(password),
		ServerIdentity: token.ServerIdentity.Encrypt(password),
		Timeout:        token.Timeout.Encrypt(password),
	}
}
