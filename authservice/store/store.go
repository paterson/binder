package store

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

type Result struct {
	SuccessfulQuery bool
	Error           error
}

type Store struct {
	db          *bolt.DB
	rootBucket  []byte
	usersBucket []byte
	Result      Result
}

var defaultStore *Store

func DefaultStore() *Store {
	if defaultStore == nil {
		db, err := bolt.Open("auth.db", 0600, nil)
		checkError(err)
		defaultStore = &Store{
			db:          db,
			rootBucket:  []byte("root"),
			usersBucket: []byte("users"),
			Result:      Result{SuccessfulQuery: false, Error: nil},
		}

		defaultStore.Result.Error = defaultStore.db.Update(func(tx *bolt.Tx) error {
			defaultStore.findOrCreateBucket(tx, defaultStore.usersBucket)
			return nil
		})
	}
	return defaultStore
}

func (s *Store) CreateUser(u *User) *Store {
	s.Result.Error = s.db.Update(func(tx *bolt.Tx) error {
		bucket := s.findOrCreateBucket(tx, s.usersBucket)
		bucket.Put([]byte(u.Username), []byte(u.Password))
		return nil
	})
	return s
}

func (s *Store) UserExists(u *User) *Store {
	s.Result.Error = s.db.View(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.usersBucket)
		val := bucket.Get([]byte(u.Username))
		s.Result.SuccessfulQuery = string(val) == u.Password
		return nil
	})
	return s
}

func (s *Store) findBucket(tx *bolt.Tx, bucketName []byte) *bolt.Bucket {
	return tx.Bucket(s.rootBucket).Bucket(bucketName)
}

func (s *Store) findOrCreateBucket(tx *bolt.Tx, bucketName []byte) *bolt.Bucket {
	rootBucket, err := tx.CreateBucketIfNotExists(s.rootBucket)
	checkError(err)
	bucket, err := rootBucket.CreateBucketIfNotExists(bucketName)
	checkError(err)
	return bucket
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
