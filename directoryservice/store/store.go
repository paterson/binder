package store

import (
	"fmt"
	"github.com/boltdb/bolt"
	"math"
	"os"
	"path/filepath"
)

type Result struct {
	Host  string
	Error error
}

type Store struct {
	db            *bolt.DB
	rootBucket    []byte
	serversBucket []byte
	Result        Result
}

var defaultStore *Store

func DefaultStore() *Store {
	if defaultStore == nil {
		db, err := bolt.Open("store/directoryservice.db", 0600, nil)
		checkError(err)
		defaultStore = &Store{
			db:            db,
			rootBucket:    []byte("root"),
			serversBucket: []byte("servers"),
			Result:        Result{Host: "", Error: nil},
		}
	}
	return defaultStore
}

func (s *Store) HostForPath(path string) *Store {
	folderPath := filepath.Dir(path)
	s.Result.Error = s.db.View(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.serversBucket)
		val := bucket.Get([]byte(folderPath))
		s.Result.Host = string(val)
		return nil
	})
	return s
}

func (s *Store) EnsureHostExistsForPath(path string) *Store {
	folderPath := filepath.Dir(path)
	s.db.View(func(tx *bolt.Tx) error {
		s.HostForPath(folderPath)
		if s.Result.Host != "" {
			return nil
		}
		// Host does not exist for this folderpath, so create.
		// Now find the host with the least number of folders.
		bucket := s.findBucket(tx, s.serversBucket)
		dict := make(map[string]int)
		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			dict[string(k)] = dict[string(k)] + 1
		}

		host := ""
		min := math.MaxInt16
		for key := range dict {
			if dict[key] < min {
				host = key
			}
		}
		s.Result.Host = host

		s.Result.Error = s.db.Update(func(tx *bolt.Tx) error {
			bucket.Put([]byte(folderPath), []byte(host))
			return nil
		})
		return nil
	})
	return s
}

func (s *Store) CreateDefaultFileServerRecord() {
	s.db.Update(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.serversBucket)
		bucket.Put([]byte("/"), []byte("localhost:3003"))
		return nil
	})
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
