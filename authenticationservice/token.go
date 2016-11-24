package main

type SessionKey               string
type ServerIdentity           string
type TimeoutPeriod            int

type EncryptedSessionKey      string
type EncryptedServerIdentity  string
type EncryptedTimeoutPeriod   int

/* Token. Key: User password */
type Token struct {
    Ticket         EncryptedTicket
    SessionKey     SessionKey
    ServerIdentity ServerIdentity
    Timeout        TimeoutPeriod
}

type EncryptedToken struct {
    Ticket          EncryptedTicket
    SessionKey      EncryptedSessionKey
    ServerIdentity  EncryptedServerIdentity
    Timeout         EncryptedTimeoutPeriod
}

func GenerateToken() Token {

}

func (encryptedToken EncryptedToken) decrypt(password string) Token {

}

func (token Token) encrypt(password string) Token {

}
