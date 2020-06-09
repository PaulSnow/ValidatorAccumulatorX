package entry

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

// Entry
// Data is entered into system by the Accumulator as a series of entries organized by chainIDs
// Unlike Factom, we will attempt to create a chain if the ChainID provided is nil.  We provide
// a function to compute the ChainID from the first entry in a chain, for use by applications.
// If a chain already exists, creating the chain will be ignored.
type Entry struct {
	Version types.VersionField // Version of this data structure
	ChainID types.Hash         // The chain id associated with this entry
	ExtIDs  types.DataField    // External ids used to create the chain id above ( see ExternalIDsToChainID() )
	Content types.DataField    // BytesSlice for holding generic data for this entry
	hash    interfaces.IHash
}

// The accumulator assumes the Entry has already been written to the (a) database.  It is only dealing with EntryHashes
type EntryHash struct {
	NewChain  bool // Create a chain.  If the chain already exists, entry is ignored
	ChainID   types.Hash
	EntryHash types.Hash
}
