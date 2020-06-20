package accumulator

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/merkleDag"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/node"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

// Accumulator
// The accumulator takes a feed of EntryHash objects to construct the cryptographic structure proving the order
// and content of the entries submitted to the Validators.  Validators validate the data, and store the data into
// key/value stores, and send streams of hashes to the Accumulators.  Validators are assumed to be knowledgeable
// of the actual use case of the system, and able to validate the data prior to submission to the accumulator.
// Of course, the Accumulator does secure and order the data, so it is reasonable that a validator may optimistically
// record entries that might be invalidated by applications after recording.
type Accumulator struct {
	chainID   *types.Hash              // Digital ID of the Accumulator.
	height    int                      // Height of the current block
	chains    map[types.Hash]*ChainAcc // Chains with new entries in this block
	entryFeed chan node.EntryHash      // Stream of entries to be placed into chains
	control   chan bool                // We are sent a "true" when it is time to end the block
	mdFeed    chan *types.Hash         // Give back the MD Hashes as they are produced
}

// Allocate the HashMap and Channels for this accumulator
// The ChainID is the Digital Identity of the Accumulator.  We will want to integrate
// useful digital IDs into the accumulator structure to ensure the integrity of the data
// collected.
func (a *Accumulator) Init(chainID *types.Hash) (
	EntryFeed chan node.EntryHash, // Return the EntryFeed channel to send ANode Hashes to the accumulator
	control chan bool, // The control channel signals End of Block to the accumulator
	mdFeed chan *types.Hash) { // the Merkle DAG Feed (mdFeed) returns block merkle DAG roots

	a.chainID = chainID
	a.chains = make(map[types.Hash]*ChainAcc, 1000)
	a.entryFeed = make(chan node.EntryHash, 10000)
	a.control = make(chan bool, 1)
	a.mdFeed = make(chan *types.Hash, 1)
	return a.entryFeed, a.control, a.mdFeed
}

type nodeEntry struct {
	chainAcc *ChainAcc
	MD       types.Hash
}

func (n nodeEntry) Marshal() (entries []byte) {
	return
}

func (a *Accumulator) Run() {
	for {
		// While we are processing a block
	block:
		for {

			// Block processing involves pulling Entries out of the entryFeed and adding
			// it to the Merkle DAG (MD)
			select {
			case ctl := <-a.control: // Have we been asked to end the block?
				if ctl {
					break block // Break block processing
				}
			case entry := <-a.entryFeed: // Get the next ANode
				chain := a.chains[entry.ChainID] // See if we have a chain for it
				if chain == nil {                // If we don't have a chain for it, then we add one to our tmp state
					chain = NewChainAcc()           // Create our collector for this chain
					a.chains[entry.ChainID] = chain // Add it to our tmp state
				}
				chain.MD.AddToChain(entry.EntryHash) // Add this entry to our chain state
			default:
				time.Sleep(100 * time.Millisecond) // If there is nothing to do, pause a bit
			}
		}

		var chainEntries []*nodeEntry
		for _, v := range a.chains {
			ne := new(nodeEntry)
			ne.chainAcc = v
			ne.MD = *v.MD.GetMDRoot()
			chainEntries = append(chainEntries, ne)
		}

		sort.Slice(chainEntries, func(i, j int) bool {
			return bytes.Compare(chainEntries[i].chainAcc.ChainID[:], chainEntries[j].chainAcc.ChainID[:]) < 0
		})

		println("\n===========================\n")
		var sum int
		for _, v := range a.chains {
			sum += len(v.MD.HashList)
		}
		fmt.Printf("%15d Entries\n", sum)

		// Calculate the ListMDRoot for all the accumulated MDRoots for all the chains
		MDAcc := new(merkleDag.MD)
		for _, v := range chainEntries {
			MDAcc.AddToChain(*v.chainAcc.MD.GetMDRoot())
		}
		a.mdFeed <- MDAcc.GetMDRoot()

		// Clear out all the chain heads, to start another round of accumulation in the next block
		a.chains = make(map[types.Hash]*ChainAcc, 1000)
	}
}
