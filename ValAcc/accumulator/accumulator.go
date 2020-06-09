package accumulator

import (
	"time"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/entry"
	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

// Accumulator
// The accumulator takes a feed of EntryHash objects to construct the cryptographic structure proving the order
// and content of the entries submitted to the Validators.  Validators validate the data, and store the data into
// key/value stores, and send streams of hashes to the Accumulators.  Validators are assumed to be knowledgeable
// of the actual use case of the system, and able to validate the data prior to submission to the accumulator.
// Of course, the Accumulator does secure and order the data, so it is reasonable that a validator may optimistically
// record entries that might be invalidated by applications after recording.
type Accumulator struct {
	height    int                  // Height of the current block
	chains    map[types.Hash]Chain // Chains with new entries in this block
	entryFeed chan entry.EntryHash // Stream of entries to be placed into chains
	control   chan time.Time       // Time when this block should be completed, and a new block started
}

func (a *Accumulator) Init() (EntryFeed chan entry.EntryHash) {
	a.chains = make(map[types.Hash]Chain, 1000)
	a.entryFeed = make(chan entry.EntryHash, 10000)
	return a.entryFeed
}
