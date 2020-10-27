// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid

import (
	"crypto/ed25519"
	"crypto/sha256"
	"runtime/debug"

	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/types"
)

var _ = debug.PrintStack

type Transaction struct {
	amount    uint64          // Amount of a factoid tx
	input     types.Address   // Address of input
	output    types.Address   // Address of output
	signature types.Signature // Signature
}

func (t *Transaction) SigBytes() (data []byte) {
	data = types.Uint64Bytes(t.amount)
	data = append(data, t.input.Bytes()...)
	data = append(data, t.output.Bytes()...)
	return data
}

func (t *Transaction) Bytes() (data []byte) {
	data = t.SigBytes()
	data = append(data, t.signature.Bytes()...)
	return data
}

func (t *Transaction) Extract(data []byte) []byte {
	t.amount, data = types.BytesUint64(data)
	data = t.input.Extract(data)
	data = t.output.Extract(data)
	data = t.signature.Extract(data)
	return data
}

func (w *Transaction) New() *Transaction {
	return new(Transaction)
}

func (t *Transaction) GetHash() (h types.Hash) {
	h = sha256.Sum256(t.Bytes())
	return h
}

func (t *Transaction) Sign(p types.PrivateKey) {
	sb := t.SigBytes()
	sig := ed25519.Sign(p[:], sb)
	copy(t.signature.Signature[:], sig)
	t.signature.PublicKey = append(t.signature.PublicKey[:0], p.GetPublicKey()...)
}

func (t *Transaction) GetSignature() types.Signature {
	return t.signature
}

func (t *Transaction) ValidateSig() bool {
	sb := t.SigBytes()
	// A signature's public key hashed must match the address
	if t.input != sha256.Sum256(t.signature.PublicKey) {
		return false
	}
	return ed25519.Verify(t.signature.PublicKey, sb, t.signature.Signature[:])
}
