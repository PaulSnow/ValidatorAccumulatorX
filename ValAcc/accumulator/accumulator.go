package accumulator

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/FactomProject/factomd/util/atomic"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/database"

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
	DB            *database.DB             // Database to hold and index the data collected by the Accumulator
	chainID       *types.Hash              // Digital ID of the Accumulator.
	height        types.BlockHeight        // Height of the current block
	chains        map[types.Hash]*ChainAcc // Chains with new entries in this block
	entryFeed     chan node.EntryHash      // Stream of entries to be placed into chains
	control       chan bool                // We are sent a "true" when it is time to end the block
	mdFeed        chan *types.Hash         // Give back the MD Hashes as they are produced
	previous      *node.Node               // Previous Directory Block
	EntryCnt      atomic.AtomicInt64       // Count of entries written
	ChainsInBlock atomic.AtomicInt64       // Count of chains written to
	ChainCnt      atomic.AtomicInt64       // Count of all chains
}

// Allocate the HashMap and Channels for this accumulator
// The ChainID is the Digital Identity of the Accumulator.  We will want to integrate
// useful digital IDs into the accumulator structure to ensure the integrity of the data
// collected.
func (a *Accumulator) Init(db *database.DB, chainID *types.Hash) (
	EntryFeed chan node.EntryHash, // Return the EntryFeed channel to send ANode Hashes to the accumulator
	control chan bool, // The control channel signals End of Block to the accumulator
	mdFeed chan *types.Hash) { // the Merkle DAG Feed (mdFeed) returns block merkle DAG roots

	a.DB = db
	a.chainID = chainID
	headHash := db.Get(types.NodeHead, chainID[:])
	if headHash != nil {
		head := db.Get(types.Node, headHash)
		if head == nil {
			panic("no head found for the directory blocks in the database")
		}
		var headNode node.Node
		_, err := headNode.Unmarshal(head)
		if err != nil {
			panic(fmt.Sprintf("error unmarshaling the head of the directory block.\n%v", err))
		}
		a.previous = &headNode
		a.height = headNode.BHeight + 1
	}
	a.chains = make(map[types.Hash]*ChainAcc, 1000)
	a.entryFeed = make(chan node.EntryHash, 10000)
	a.control = make(chan bool, 1)
	a.mdFeed = make(chan *types.Hash, 1)

	fmt.Sprintf("Starting the Accumulator at height %d\n", a.height)

	return a.entryFeed, a.control, a.mdFeed
}

func (a *Accumulator) GetEntryFeed() chan node.EntryHash {
	return a.entryFeed
}

func (a *Accumulator) Run() {
	var totalEntries int64  // We count the entries and chains as we go, but update the atomic counts
	var ChainsInBlock int64 //  at the end of each block

	var goWrites atomic.AtomicInt

	for {
		// While we are processing a block
	block:
		for {

			// Block processing involves pulling Entries out of the entryFeed and adding
			// it to the Merkle DAG (MD)
			select {
			case ctl := <-a.control: // Have we been asked to end the block?
				if ctl {
					println("Processing EOB ", a.height)
					a.height++
					break block // Break block processing
				}
			default:
				select {
				case entry := <-a.entryFeed: // Get the next ANode
					chain := a.chains[entry.ChainID] // See if we have a chain for it
					totalEntries++
					if chain == nil { // If we don't have a chain for it, then we add one to our tmp state
						ChainsInBlock++
						chain = NewChainAcc(*a.DB, entry, a.height) // Create our collector for this chain
						a.chains[entry.ChainID] = chain             // Add it to our tmp state
						chain.MD.AddToChain(entry.EntryHash)        // Add this entry to our chain state
					} else {
						// This is where we make sure every Entry added to a chain is a non-duplicate to all
						// entries.  This assumes that the chains for an accumulator are unique to that accumulator,
						// which is true by design.  So if the entry isn't in the chain right now, and not in the db,
						// then it is unique.
						if chain.entries[entry.EntryHash] == 0 { // Added this entry to this chain already?
							if a.DB.Get(types.EntryNode, entry.EntryHash.Bytes()) == nil { // Have the entry in the DB already?
								chain.entries[entry.EntryHash] = 1   // No? Then mark it in the chain
								chain.MD.AddToChain(entry.EntryHash) // Add it to the chain
							}
						}
					}
				default:
					time.Sleep(100 * time.Millisecond) // If there is nothing to do, pause a bit
				}
			}
		}

		if goWrites.Load() > 0 {
			fmt.Println("Waiting on", goWrites.Load(), "database updates.")
			for goWrites.Load() > 0 {
				time.Sleep(1 * time.Second)
			}
		}

		var chainEntries []node.NEList
		for _, v := range a.chains {
			v.Node.ListMDRoot = *v.MD.GetMDRoot()
			v.Node.EntryList = v.MD.HashList
			v.Node.IsNode = false

			tNode := v.Node
			go func() {
				goWrites.Add(1)
				tNode.Put(a.DB)
				goWrites.Add(-1)
			}()

			ne := new(node.NEList)
			ne.ChainID = v.Node.ChainID
			ne.MDRoot = v.Node.ListMDRoot
			chainEntries = append(chainEntries, *ne)

		}

		sort.Slice(chainEntries, func(i, j int) bool {
			return bytes.Compare(chainEntries[i].ChainID[:], chainEntries[j].ChainID[:]) < 0
		})

		// Print some statistics
		var sum int
		for _, v := range a.chains {
			sum += len(v.MD.HashList)
		}

		a.EntryCnt.Store(totalEntries)
		a.ChainsInBlock.Store(ChainsInBlock)
		a.ChainCnt.Add(ChainsInBlock)
		ChainsInBlock = 0

		// Calculate the ListMDRoot for all the accumulated MDRoots for all the chains
		MDAcc := new(merkleDag.MD)
		for _, v := range chainEntries {
			MDAcc.AddToChain(v.MDRoot)
		}

		// Populate the directory block with the data collected over the last block period.
		directoryBlock := new(node.Node)
		directoryBlock.Version = types.Version
		directoryBlock.ChainID = *a.chainID
		directoryBlock.BHeight = a.height
		if directoryBlock.SequenceNum > 0 {
			directoryBlock.Previous = *a.previous.GetHash()
		}
		directoryBlock.SequenceNum = types.Sequence(a.height)
		directoryBlock.TimeStamp = types.TimeStamp(time.Now().UnixNano())
		directoryBlock.IsNode = true
		lMDR := MDAcc.GetMDRoot()
		if lMDR != nil {
			directoryBlock.ListMDRoot = *lMDR
		}

		// Write the directory
		directoryBlock.Put(a.DB)

		a.mdFeed <- directoryBlock.GetMDRoot()

		// Clear out all the chain heads, to start another round of accumulation in the next block
		a.chains = make(map[types.Hash]*ChainAcc, 1000)
	}
}
