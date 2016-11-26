package request

const (
	SERVER_KEY = "lkm4iuPKCCJQGBGB"
)

/* Ticket. Key: Server Key */
type Ticket struct {
	SessionKey SessionKey
}

type EncryptedTicket struct {
	SessionKey EncryptedSessionKey
}

func (encryptedTicket EncryptedTicket) Decrypt() Ticket {
	return Ticket{SessionKey: encryptedTicket.SessionKey.Decrypt(SERVER_KEY)}
}

func (ticket Ticket) Encrypt() EncryptedTicket {
	return EncryptedTicket{SessionKey: ticket.SessionKey.Encrypt(SERVER_KEY)}
}
