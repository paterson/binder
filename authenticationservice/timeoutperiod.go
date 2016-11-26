package main

import (
    "github.com/paterson/binder/authenticationservice/encryption"
    "log"
)

type TimeoutPeriod            string
type EncryptedTimeoutPeriod   string

func (timeoutPeriod TimeoutPeriod) Encrypt(key string) EncryptedTimeoutPeriod {
    cipher := encryption.Encrypt(key, timeoutPeriod.toString())
    return EncryptedTimeoutPeriod(cipher)
}

func (timeoutPeriod EncryptedTimeoutPeriod) Decrypt(key string) TimeoutPeriod {
    text, err := encryption.Decrypt(key, timeoutPeriod.toString())
    if err != nil {
        log.Fatal(err) // for now.....
    }
    return TimeoutPeriod(text)
}

func (timeoutPeriod TimeoutPeriod) toString() string {
    return string(timeoutPeriod)
}

func (timeoutPeriod EncryptedTimeoutPeriod) toString() string {
    return string(timeoutPeriod)
}

var DefaultTimeout = TimeoutPeriod(30)
