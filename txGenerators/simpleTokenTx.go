package txGenerators

import (
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/PaulSnow/LoadTest/accumulator"
)

func GenSimpleTokenTxs(router accumulator.Router) {
	// Addresses transacting
	// An initial balance
	addresses := append([]float64{}, 1000000)

	for {
		// pick two addresses
		SAdr := rand.Int() % len(addresses)
		var BSadr []byte
		v := SAdr
		for v > 0 {
			BSadr = append(BSadr, byte(v))
			v >>= 8
		}

		HSAdr := sha256.Sum256(BSadr)
		DAdr := rand.Int() % len(addresses)
		if (len(addresses) < 1000 || rand.Float32() < .1) && len(addresses) < 100000 {
			SAdr = 0
			DAdr = len(addresses)
			addresses = append(addresses, 0)
		}
		// Make transfers to different addresses
		if SAdr == DAdr {
			continue
		}

		amt := addresses[SAdr] * (rand.Float64() + 1) / 10

		tx := fmt.Sprintf("BSadr: %d  dest: %d  amt: %f ", SAdr, DAdr, amt)

		addresses[SAdr] -= amt
		addresses[DAdr] += amt

		h := sha256.Sum256([]byte(tx))
		router.AddTx(HSAdr, h[:])
	}

}
