// GenSimpleTokenTxs
// generates simple token transactions without digital signatures and maintains the balances of these
// transactions in memory.  We start with a million transactions and spread them around the network.
// The transactions are sent to the router with the source address and the hash of the transaction.

package txGenerators

import (
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/PaulSnow/ValidatorAccumulator/Simple/accumulator"
)

const maxaddresses = 100000

func GenSimpleTokenTxs(router accumulator.Router) {
	// Addresses transacting
	// An initial balance
	addresses := append([]float64{}, 1000000)

	for {
		// pick two addresses.  Addresses are simply indexes into the array of balances.
		SAdr := rand.Int() % len(addresses)

		// The chain ID to hold the transactions is the hash of the index into the addresses array
		var BSadr []byte
		v := SAdr
		for v > 0 {
			BSadr = append(BSadr, byte(v))
			v >>= 8
		}

		HSAdr := sha256.Sum256(BSadr)

		// Calculate a destination address.
		DAdr := rand.Int() % len(addresses)

		// We move the tokens into the first 1000 transactions, but 10% of the time, create a new address
		// up to a limit.
		if (len(addresses) < 1000 || rand.Float32() < .1) && len(addresses) < maxaddresses {
			SAdr = 0
			DAdr = len(addresses)
			addresses = append(addresses, 0)
		}
		// Make transfers to different addresses
		if SAdr == DAdr {
			continue
		}

		// We move 10% of the source balance
		amt := addresses[SAdr] * (rand.Float64() + 1) / 10

		// The transaction is just the string of the transaction.
		tx := fmt.Sprintf("BSadr: %d  dest: %d  amt: %f ", SAdr, DAdr, amt)

		// update balances
		addresses[SAdr] -= amt
		addresses[DAdr] += amt

		// Get the hash of the transaction
		h := sha256.Sum256([]byte(tx))

		// Write the transaction to the ChainAcc ID in the accumulator
		router.AddTx(HSAdr, h[:])
	}

}
