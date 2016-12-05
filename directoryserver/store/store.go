package store

import (
	"fmt"
	"github.com/boltdb/bolt"
	"math"
	"os"
	"strings"
)

type Result struct {
	Hosts string
	Error error
}

type Store struct {
	db          *bolt.DB
	rootBucket  []byte
	filesBucket []byte
	Result      Result
}

var defaultStore *Store

func DefaultStore() *Store {
	if defaultStore == nil {
		db, err := bolt.Open("directoryservice.db", 0600, nil)
		checkError(err)
		defaultStore = &Store{
			db:          db,
			rootBucket:  []byte("root"),
			filesBucket: []byte("files"),
			Result:      Result{Hosts: "", Error: nil},
		}
	}
	return defaultStore
}

func (s *Store) HostsForPath(path string) *Store {
	s.Result.Error = s.db.View(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.filesBucket)
		val := bucket.Get([]byte(path))
		fmt.Println("Bucket val for path:", string(val))
		s.Result.Hosts = string(val)
		return nil
	})
	return s
}

func (s *Store) EnsureHostExistsForPath(path string) *Store {
	s.db.View(func(tx *bolt.Tx) error {
		s.HostsForPath(path)
		if s.Result.Hosts != "" {
			return nil
		}
		// Host does not exist for this path, so create it.
		// Now find the host with the least number of files.
		bucket := s.findBucket(tx, s.filesBucket)
		dict := make(map[string]int)
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			hosts := strings.Split(string(v), ",")
			for _, host := range hosts {
				dict[host] = dict[host] + 1
			}
		}
		host := ""
		min := math.MaxInt16
		for key := range dict {
			if dict[key] < min {
				host = key
			}
		}
		s.Result.Hosts = host
		s.Result.Error = s.db.Update(func(tx *bolt.Tx) error {
			bucket := s.findOrCreateBucket(tx, s.filesBucket)
			bucket.Put([]byte(path), []byte(host))
			return nil
		})
		fmt.Println("Bucket val for path:", []byte(bucket.Get([]byte(path))))
		return nil
	})
	return s
}

func (s *Store) CreateDefaultFileServerRecord() {
	s.db.Update(func(tx *bolt.Tx) error {
		rootBucket, _ := tx.CreateBucketIfNotExists(s.rootBucket)
		bucket, _ := rootBucket.CreateBucketIfNotExists(s.filesBucket)
		hosts := []string{"http://localhost:3002", "http://localhost:3003"}
		str := strings.Join(hosts, ",")
		bucket.Put([]byte("/"), []byte(str))
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
