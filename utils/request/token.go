package request

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
    sessionKey := NewSessionKey()
    ticket     := Ticket{SessionKey: sessionKey}
    return Token{
        Ticket:         ticket.Encrypt(),
        SessionKey:     sessionKey,
        ServerIdentity: AuthServerIdentity,
        Timeout:        DefaultTimeout,
    }
}

func (encryptedToken EncryptedToken) decrypt(password string) Token {
    return Token{
        Ticket:         encryptedToken.Ticket,
        SessionKey:     encryptedToken.SessionKey.Decrypt(password),
        ServerIdentity: encryptedToken.ServerIdentity.Decrypt(password),
        Timeout:        encryptedToken.Timeout.Decrypt(password),
    }
}

func (token Token) encrypt(password string) EncryptedToken {
    return EncryptedToken{
        Ticket:         token.Ticket,
        SessionKey:     token.SessionKey.Encrypt(password),
        ServerIdentity: token.ServerIdentity.Encrypt(password),
        Timeout:        token.Timeout.Encrypt(password),
    }
}
