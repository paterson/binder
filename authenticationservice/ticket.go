package main

/* Ticket. Key: Server Key */
type Ticket struct {
    SessionKey SessionKey
}

type EncryptedTicket struct {
    SessionKey EncryptedSessionKey
}

func (encryptedTicket EncryptedTicket) decrypt() Ticket {

}

func (ticket Ticket) encrypt() EncryptedTicket {

}
