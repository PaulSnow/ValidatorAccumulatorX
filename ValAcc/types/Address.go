package types

import (
	"crypto/ed25519"
	"crypto/sha256"
)

// Address
// ===========================================================================
type Address [32]byte // An Address is the hash of a public key

// Copy
// Return a copy of the hash
func (a Address) Copy() *Address {
	return &a
}

// Bytes
// Return a []byte for the Address
func (a Address) Bytes() []byte {
	return a[:]
}

// Extract
// Extract the address from the given byte slice
func (a *Address) Extract(data []byte) []byte {
	copy((*a)[:], data[:32])
	return data[32:]
}

// Signature
//=============================================================================
type Signature struct {
	PublicKey []byte
	Signature [64]byte
}

// Copy
// Return a copy of the hash
func (a Signature) Copy() *Signature {
	return &a
}

// Bytes
// Return a []byte for the Address
func (a Signature) Bytes() (b []byte) {
	b = append(b, a.PublicKey...)
	b = append(b, a.Signature[:]...)
	return b
}

// Extract
// Extract the address from the given byte slice
func (a *Signature) Extract(data []byte) []byte {
	copy((a.PublicKey), data[:32])
	data = data[:32]
	copy((a.Signature)[:], data[:64])
	data = data[:64]
	return data
}

// Private Key
// ============================================================================
type PrivateKey [64]byte

// The public key is the last 32 bytes of the PrivateKey
func (p *PrivateKey) GetPublicKey() (publicKey []byte) {
	publicKey = append(publicKey, p[32:]...)
	return publicKey
}

// GetAddress
func (p *PrivateKey) GetAddress() [32]byte {
	return sha256.Sum256(p[:])
}

// Sign
func (p *PrivateKey) Sign(data []byte) []byte {
	return ed25519.Sign(p[:], data)
}

func (p *PrivateKey) Verify(data []byte, signature []byte) bool {
	return ed25519.Verify(p[:], data, signature)
}
