package types

// We are typing varies fields used in the AccVal structures.  In particular, a Hash is a
// reference to a 32 byte array.  From that one can easily have a 32 byte slice, or a
// pointer to a 32 byte array as needed.  And from the type Hash, one can easily have a copy
// of the Hash, and it is easy to compare to 32 byte arrays (thus compare to hashes).
//
// This file also defines the various buckets we use in our database to index the data
// collected by accumulators.
//
// Accumulators make no demands on Validators, and in fact make no demands on the implementation
// of sub Accumulators other than this implementation requires sub accumulators to represent
// the first sub field of a ChainID.

import (
	"os"
	"os/user"
)

// ======================= Helper Functions =================================

// GetHomeDir
// Used to find the Home Directory from which the configuration directory for the ValAcc application to
// use for its database.  This is not a terribly refined way of configuring the ValAcc and may be
// refined in the future.
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

// BoolBytes
// Marshal a Bool
func BoolBytes(b bool) []byte {
	if b {
		return append([]byte{}, 1)
	}
	return append([]byte{}, 0)
}

// BytesBool
// Unmarshal a Uint8
func BytesBool(data []byte) (f bool, newData []byte) {
	if data[0] != 0 {
		f = true
	}
	return f, data[1:]
}

// Uint16Bytes
// Marshal a int32 (big endian)
func UInt16Bytes(i uint16) []byte {
	return append([]byte{}, byte(i>>8), byte(i))
}

// BytesUint16
// Unmarshal a uint32 (big endian)
func BytesUint16(data []byte) (uint16, []byte) {
	return uint16(data[0])<<8 + uint16(data[1]), data[2:]
}

// Uint32Bytes
// Marshal a int32 (big endian)
func UInt32Bytes(i uint32) []byte {
	return append([]byte{}, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

// BytesUint32
// Unmarshal a uint32 (big endian)
func BytesUInt32(data []byte) (uint32, []byte) {
	return uint32(data[0])<<24 + uint32(data[1]<<16) + uint32(data[2])<<8 + uint32(data[3]), data[4:]
}
