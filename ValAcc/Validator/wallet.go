package factoid

import (
	"crypto/ed25519"
	"crypto/sha256"

	"github.com/AccumulusNetwork/ValidatorAccumulator/ValAcc/types"
)

type walletEntry struct {
	PrivateKey types.PrivateKey
	Address    types.Address
	Amount     uint64
}

type Wallet struct {
	CurrentSeed [32]byte
	MyKeys      map[types.Address]walletEntry
}

func NewWallet(seed []byte) *Wallet {
	w := new(Wallet)
	w.MyKeys = make(map[types.Address]walletEntry, 100)
	return w
}

func (w *Wallet) NewAddress() types.Address {
	w.CurrentSeed = sha256.Sum256(w.CurrentSeed[:])

	we := new(walletEntry)
	privateKey := ed25519.NewKeyFromSeed(w.CurrentSeed[:])
	if len(privateKey) != 64 {
		panic("private key is the wrong length")
	}
	copy(we.PrivateKey[:], privateKey)
	we.Address = sha256.Sum256(we.PrivateKey.GetPublicKey())
	w.MyKeys[we.Address] = *we
	return we.Address
}

func (w *Wallet) GetPrivateKey(address types.Address) types.PrivateKey {
	we := w.MyKeys[address]
	return we.PrivateKey
}
