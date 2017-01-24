package store

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/paterson/binder/utils/constants"
	"os"
	"sort"
	"strings"
)

type Result struct {
	Hosts        string
	Error        error
	Locked       bool
	ValidLockKey bool
}

type Store struct {
	db          *bolt.DB
	rootBucket  []byte
	filesBucket []byte
	locksBucket []byte
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
			locksBucket: []byte("locks"),
			Result:      Result{Hosts: "", Error: nil, Locked: false, ValidLockKey: false},
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
	s.db.Update(func(tx *bolt.Tx) error {
		s.HostsForPath(path)
		if s.Result.Hosts != "" {
			return nil
		}
		// Host does not exist for this path, so create some.
		// Use the hosts with the least number of files in them.
		bucket := s.findBucket(tx, s.filesBucket)
		dict := make(map[string]int)
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			hosts := strings.Split(string(v), ",")
			for _, host := range hosts {
				dict[host] = dict[host] + 1
			}
		}
		sortedKeys := sortMapKeysByValue(dict)
		n := min(len(sortedKeys), 4) // Use at most 4 hosts for replication
		hosts := strings.Join(sortedKeys[0:n], ",")
		s.Result.Hosts = hosts
		bucket.Put([]byte(path), []byte(hosts))
		return nil
	})
	return s
}

// Valid Lock key if value for filepath == key
func (s *Store) IsValidLockKeyForPath(key, filepath string) *Store {
	s.db.View(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.locksBucket)
		val := bucket.Get([]byte(filepath))
		s.Result.ValidLockKey = string(val) == key && key != ""
		return nil
	})
	return s
}

// Locked if value for filepath is not empty
func (s *Store) GetLockStatusForPath(filepath string) *Store {
	s.db.View(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.locksBucket)
		val := bucket.Get([]byte(filepath))
		s.Result.Locked = string(val) != ""
		return nil
	})
	return s
}

func (s *Store) LockPathWithKey(filepath, key string) *Store {
	s.db.Update(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.locksBucket)
		bucket.Put([]byte(filepath), []byte(key))
		return nil
	})
	return s
}

func (s *Store) UnlockPathWithKey(filepath, key string) *Store {
	s.db.Update(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.locksBucket)
		val := bucket.Get([]byte(filepath))
		if string(val) == key && key != "" {
			bucket.Put([]byte(filepath), []byte(""))
			s.Result.Locked = false
		}
		return nil
	})
	return s
}

func (s *Store) Seed() {
	s.db.Update(func(tx *bolt.Tx) error {
		rootBucket, _ := tx.CreateBucketIfNotExists(s.rootBucket)
		bucket, _ := rootBucket.CreateBucketIfNotExists(s.filesBucket)
		rootBucket.CreateBucketIfNotExists(s.locksBucket)
		bucket.Put([]byte("/"), []byte(""))
		return nil
	})
}

func (s *Store) AddFileserver(fileserver string) *Store {
	s.db.Update(func(tx *bolt.Tx) error {
		bucket := s.findBucket(tx, s.locksBucket)
		val := bucket.Get([]byte("/"))
		hosts := strings.Split(string(val), ",")
		hosts = append(hosts, fileserver)
		str := strings.Join(hosts, ",")
		bucket.Put([]byte("/"), []byte(str))
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

// Takes a map of strings -> int, and returns back the strings ordered by their corresponding int value
func sortMapKeysByValue(m map[string]int) (res []string) {
	var a []int
	n := map[int][]string{}

	for k, v := range m {
		n[v] = append(n[v], k)
	}

	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.IntSlice(a))
	for _, k := range a {
		for _, s := range n[k] {
			res = append(res, s)
		}
	}
	return res
}

// math.Min is for float64, and returns a float64, so creating a faster simpler method for min of two ints.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
