package store

import (
	"fmt"
	"github.com/boltdb/bolt"
	"math"
	"os"
	"strings"
)

type Result struct {
	Hosts []string
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
			Result:      Result{Host: "", Error: nil},
		}
	}
	return defaultStore
}

func (s *Store) HostsForPath(path string) *Store {
	s.Result.Error = s.db.View(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.filesBucket)
		val := bucket.Get([]byte(path))
		s.Result.Hosts = strings.Split(string(val), ",")
		return nil
	})
	return s
}

func (s *Store) EnsureHostExistsForPath(path string) *Store {
	s.db.View(func(tx *bolt.Tx) error {
		s.HostForPath(path)
		if s.Result.Host != "" {
			return nil
		}
		// Host does not exist for this path, so create it.
		// Now find the host with the least number of files.
		bucket := s.findBucket(tx, s.filesBucket)
		dict := make(map[string]int)
		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			hosts := strings.Split(string(val), ",")
			for host := range hosts {
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
		rootBucket, _ := tx.CreateBucketIfNotExists(s.rootBucket)
		bucket, _ := rootBucket.CreateBucketIfNotExists(s.filesBucket)
		hosts := []string{"http://localhost:3002"}
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
