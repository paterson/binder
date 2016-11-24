package main

import (
    "github.com/boltdb/bolt"
)

type Store struct {
    db *bolt.DB
}

var defaultStore *Store

func DefaultStore() *Store {
    if defaultStore == nil {
		defaultStore = &Store{}
	}
	return defaultStore
}

func (s *Store) CreateUser(u *User) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        // Retrieve the users bucket.
        // This should be created when the DB is first opened.
        bucket := tx.CreateBucketIfNotExists([]byte("users"))
        return bucket.Put(u.Username, u.Password)
    })
}

func (s *Store) UserExists(username string, password string) bool {
    found := false
    s.db.View(func(tx *bolt.Tx) error {
        bucket := tx.CreateBucketIfNotExists([]byte("users"))

        val := bucket.Get(username)
        found := val == password && val != ""
        return nil
    })
    return found
}
