package node

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

// ANode
// Data is entered into system by the Accumulator as a series of entries organized by chainIDs
// Unlike Factom, we will attempt to create a chain if the ChainID provided is nil.  We provide
// a function to compute the ChainID from the first entry in a chain, for use by applications.
// If a chain already exists, creating the chain will be ignored.
//
// Node
//      Version					uint8
//      Block Height			uint32
//      SequenceNum             uint32
//      TimeStamp				int64
//      ChainID					[32]byte
//		Previous Node Header	[32]byte
//		SubNode/Entries flag	uint8
//		List MDRoot             [32]byte
//		Len(List)				uint32
//		List (SubNodes/Entries)
//			Chain/SubChain ID	[32]byte     Sorted by Chain/SubChain ID
//			Node/ANode Hash     [32]byte

type Node struct {
	Version     types.VersionField // Version of this data structure
	BHeight     types.BlockHeight  // Block Height
	SequenceNum types.Sequence     // Sequence Number for this chain of nodes
	TimeStamp   types.TimeStamp    // TimeStamp by Accumulator when the structure was built
	ChainID     types.Hash         // The ChainID (Directory Block - zeros, SubNodes - 1st ChainID, Entries - ChainID)
	SubChainIDs []types.Hash       // SubChainIDs to build ChainID, Directory Block - zeros
	Previous    types.Hash         // Hash of previous Block Header
	IsNode      bool               // IsNode is true for a node, is false for an entry
	MDRoot      types.Hash         // Merkle DAG of the entries of the List (only the hashes)
	List        []NEList           // List of ChainIDs/MDRoots for nodes/entries
}

// NEList
// Node ANode List (NEList) is a struct of a ChainID and a Node or ANode Hash
type NEList struct {
	ChainID types.Hash // Chain or SubChain ID that leads to a node, or a ChainID that leads to an ANode
	MDRoot  types.Hash // Merkle Dag of either sub nodes or entries
}

// SameAs
// Compares two entries, mostly used for testing
func (n Node) SameAs(n2 Node) bool {
	if n.Version != n2.Version {
		return false
	}
	if n.BHeight != n2.BHeight {
		return false
	}
	if n.SequenceNum != n2.SequenceNum {
		return false
	}
	if n.TimeStamp != n2.TimeStamp {
		return false
	}
	if n.ChainID != n2.ChainID {
		return false
	}
	if len(n.SubChainIDs) != len(n2.SubChainIDs) {
		return false
	}
	for i := range n.SubChainIDs {
		if n.SubChainIDs[i] != n2.SubChainIDs[i] {
			return false
		}
	}
	if n.Previous != n2.Previous {
		return false
	}
	if n.IsNode != n2.IsNode {
		return false
	}
	if n.MDRoot != n2.MDRoot {
		return false
	}
	if len(n.List) != len(n2.List) {
		return false
	}
	for i, list := range n.List {
		if list.ChainID != n2.List[i].ChainID {
			return false
		}
		if list.MDRoot != n2.List[i].MDRoot {
			return false
		}
	}
	return true
}

// Marshal
// Convert the given entry into a byte slice. Add to that the SubChainIDs of the ChainID. Returns nil
// if anything goes wrong while marshaling
func (n Node) Marshal() (bytes []byte) {

	// On any error, return a nil for the byte representation of the ANode
	defer func() {
		if r := recover(); r != nil {
			bytes = nil
			return
		}
	}()

	bytes = append(bytes, n.Version.Bytes()...) // Put the version into the slice
	bytes = append(bytes, n.BHeight.Bytes()...)
	bytes = append(bytes, n.SequenceNum.Bytes()...)
	bytes = append(bytes, n.TimeStamp.Bytes()...)
	bytes = append(bytes, n.ChainID.Bytes()...)                             // Put the ChainID into the slice
	bytes = append(bytes, types.UInt16Bytes(uint16(len(n.SubChainIDs)))...) // Put the number of SubChains
	for _, subChain := range n.SubChainIDs {                                // For each SubChain
		bytes = append(bytes, subChain.Bytes()...) // Put the ExtID's data in the slice
	}
	bytes = append(bytes, n.Previous.Bytes()...)
	if n.IsNode {
		bytes = append(bytes, 1)
	} else {
		bytes = append(bytes, 0)
	}
	bytes = append(bytes, n.MDRoot.Bytes()...)
	bytes = append(bytes, types.UInt32Bytes(uint32(len(n.List)))...) // Put the number of SubChains
	for _, list := range n.List {                                    // For each SubChain
		bytes = append(bytes, list.ChainID.Bytes()...) // Chain/SubChain ID
		bytes = append(bytes, list.MDRoot.Bytes()...)  // MD of the sub node or entry
	}
	return bytes // Return the slice
}

// GetHash
// Return the hash for this sub node or entry node
func (n Node) GetHash() (hash *types.Hash) {
	h := n.Marshal() // Get the bytes behind the EntryHash
	if h == nil {    // A nil would mean the ANode didn't marshal
		return nil
	}
	hash = new(types.Hash) // Get the Hash object to return
	hs := sha256.Sum256(h) // Get the array holding the hash (so we can create a slice)
	hash.Extract(hs[:])    // Populate the Hash object
	return hash            // Return a pointer to the Hash object.
}

// Unmarshal
// Extract an entry from a byte slice.  Returns an error if the unmarshal fails, or the length of the
// data consumed and a nil.
func (n *Node) Unmarshal(data []byte) (dataConsumed int, err error) {

	// On any error, no data is consumed and return an error as to why unmarshal fails
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ANode Failed to unmarshal %v", r))
		}
	}()
	d := data // d keeps the original slice

	data = n.Version.Extract(data)     // Extract the version
	data = n.BHeight.Extract(data)     // Extract the BlockHeight
	data = n.SequenceNum.Extract(data) // Extract the BlockHeight
	data = n.TimeStamp.Extract(data)   // Extract the TimeStamp
	data = n.ChainID.Extract(data)     // Extract the ChainID
	n.SubChainIDs = n.SubChainIDs[0:0] // Clear any ExtIDs that might already be in this ANode
	// Pull out all the subChain IDs
	var numSubChains uint16
	numSubChains, data = types.BytesUint16(data) // Get the number of SubChainIDs we should have
	for i := uint16(0); i < numSubChains; i++ {  // Pull each of them out of the data slice
		sc := types.Hash{}                        // Get a Hash to put the SubChainID in
		data = sc.Extract(data)                   // Extract the ExtID
		n.SubChainIDs = append(n.SubChainIDs, sc) // Put it in the ExtID list
	}
	data = n.Previous.Extract(data)
	n.IsNode, data = types.BytesBool(data) // Extract the node/entries flag
	data = n.MDRoot.Extract(data)
	// Pull out all the List entries
	var listLen uint32
	listLen, data = types.BytesUInt32(data)
	for i := uint32(0); i < listLen; i++ {
		ne := new(NEList)
		data = ne.ChainID.Extract(data)
		data = ne.MDRoot.Extract(data)
		n.List = append(n.List, *ne)
	}

	return len(d) - len(data), nil // Return the bytes consumed and a nil that all is well for an error

}
