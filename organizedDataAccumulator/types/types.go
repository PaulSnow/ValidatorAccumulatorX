package types

import (
	"crypto/sha256"
	"os"
	"os/user"
)

type VersionField uint8 // We are typing certain things in the protocol
type Hash [32]byte      // We are currently using sha256 hashes
type DataField []byte   // Typing the Data Fields to allow flexibility in representation in the future

var _ Hash

//================= Helper functions

// Return a copy of the hash
func (h Hash) Copy() *Hash {
	return &h
}

// Hash this hash (the left hash) with the given right hash to produce a new hash
func (h Hash) Combine(right Hash) *Hash {
	sum := sha256.New()
	x := sha256.Sum256(h[:]) // Process the left side, i.e. v from this position in c.MD
	sum.Write(x[:])
	x = sha256.Sum256(right[:]) // Process the right side, i.e. whatever hash combinations we have in hash
	sum.Write(x[:])
	var combinedHash Hash
	copy(combinedHash[:], sum.Sum(nil))
	return &combinedHash
}

func GetHomeDir() string {
	valAccHome := os.Getenv("VALACC")
	if valAccHome != "" {
		return valAccHome
	}

	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}
	return homeDir
}

// Bucket Names used by the accumulator and validator

const (
	Entry                = "entry"
	DirectoryBlockHeight = "dbHeight"
	DirectoryBlockHash   = "dbHash"
	Node                 = "node"
)
