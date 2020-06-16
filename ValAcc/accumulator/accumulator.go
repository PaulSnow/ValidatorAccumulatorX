package accumulator

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/merkleDag"
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
	height    int                      // Height of the current block
	chains    map[types.Hash]*ChainAcc // Chains with new entries in this block
	entryFeed chan merkleDag.EntryHash // Stream of entries to be placed into chains
	control   chan bool                // We are sent a "true" when it is time to end the block
	mdFeed    chan *types.Hash         // Give back the MD Hashes as they are produced
}

// Allocate the HashMap and Channels for this accumulator
func (a *Accumulator) Init() (EntryFeed chan merkleDag.EntryHash, control chan bool, mdFeed chan *types.Hash) {
	a.chains = make(map[types.Hash]*ChainAcc, 1000)
	a.entryFeed = make(chan merkleDag.EntryHash, 10000)
	a.control = make(chan bool, 1)
	a.mdFeed = make(chan *types.Hash, 1)
	return a.entryFeed, a.control, a.mdFeed
}

type nodeEntry struct {
	chainID types.Hash
	MD      types.Hash
}

func (a *Accumulator) Run() {
	for {
		// While we are processing a block
	block:
		for {

			// Check to see if the block has ended.  If so, break block processing
			// so we can tie up the block, and start the next one.
			select {
			case <-a.control:
				break block
			default:
			}

			// Block processing involves pulling Entries out of the entryFeed and adding
			// it to the Merkle DAG (MD)
			entry := <-a.entryFeed           // Get the next Entry
			chain := a.chains[entry.ChainID] // See if we have a chain for it
			if chain == nil {                // If we don't have a chain for it, then we add one to our tmp state
				chain = NewChainAcc() // Create our collector for this chain

				a.chains[entry.ChainID] = chain // Add it to our tmp state
			}
			chain.MD.AddToChain(entry.EntryHash) // Add this entry to our chain state
		}

		var chainEntries []*nodeEntry
		for k, v := range a.chains {
			ne := new(nodeEntry)
			ne.chainID = k
			ne.MD = *v.MD.GetMDRoot()
			chainEntries = append(chainEntries, ne)
		}
		sort.Slice(chainEntries, func(i, j int) bool {
			return bytes.Compare(chainEntries[i].chainID[:], chainEntries[j].chainID[:]) < 0
		})

		println("\n===========================\n")
		var sum int
		for _, v := range a.chains {
			sum += len(v.MD.HashList)
		}
		fmt.Printf("%15d Entries\n", sum)
		a.chains = make(map[types.Hash]*ChainAcc, 1000)
		a.mdFeed <- nil
	}
}
