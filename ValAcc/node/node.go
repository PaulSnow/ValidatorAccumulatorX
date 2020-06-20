package node

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/database"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

// Node
// Data is entered into system by the Accumulator as a series of entries organized by chainIDs
// Unlike Factom, we will attempt to create a chain if the ChainID provided is nil.  We provide
// a function to compute the ChainID from the first entry in a chain, for use by applications.
// If a chain already exists, creating the chain will be ignored.

type Node struct {
	Version     types.VersionField // Version of this data structure
	BHeight     types.BlockHeight  // Block Height
	SequenceNum types.Sequence     // Sequence Number for this chain of nodes
	TimeStamp   types.TimeStamp    // TimeStamp by Accumulator when the structure was built
	ChainID     types.Hash         // The ChainID (Directory Block - zeros, SubNodes - 1st ChainID, Entries - ChainID)
	SubChainIDs []types.Hash       // SubChainIDs to build ChainID, Directory Block - zeros
	Previous    types.Hash         // Hash of previous Block Header
	IsNode      bool               // IsNode is true for a node, is false for an entry
	ListMDRoot  types.Hash         // Merkle DAG of the entries of the List (only the hashes)
	List        []NEList           // List of ChainIDs/MDRoots for Directory block or sub block nodes
	EntryList   []types.Hash       // List of Entry Hashes for an Entry node

	MarshalCache []byte // Cache of the marshaled form of the node.  Do NOT marshal a node unless
	//   the node is completely formed!
}

// NEList
// Node List (NEList) is a struct of a ChainID and a Node Hash
type NEList struct {
	ChainID types.Hash // Chain or SubChain ID that leads to a node, or a ChainID that leads to an ANode
	MDRoot  types.Hash // Merkle Dag of either sub nodes or entries
}

// Put
// Put this node into the database.  There is a little special treatment for the Directory Blocks.
// In that case, the ChainID is the DID for the root Accumulator, and there are no SubChainIDs.
func (n Node) Put(db *database.DB) error {
	nHash := n.GetHash()[:]

	// So first do some indexing around the chain of nodes for this ChainID.  Set nodeFirst, nodeNext, nodeHead

	// Get the last node recorded for this ChainID (that's the head hash)
	headHash := db.Get(types.NodeHead, n.ChainID[:])
	if headHash == nil && n.SequenceNum != 0 { // If that's nil, and our sequence number isn't zero, bad stuff is about!
		return errors.New(fmt.Sprintf("chainID %x not found in DB, with sequence number %d", n.ChainID, n.SequenceNum))
	} else if headHash == nil { // If we have no previous hash and our sequence number is zero, this is our first!
		db.Put(types.NodeFirst, n.ChainID[:], nHash)
	} else { // Otherwise if I have a previous hash, then create an index from it to this node
		db.Put(types.NodeNext, headHash, nHash)
	}
	db.Put(types.NodeHead, n.ChainID.Bytes(), nHash)

	// If a node does not have any SubChains to define its ChainID, then its ChainID is really
	// the DID for the root accumulator, and this is a Directory Block.  So we will index it
	// against the block height.  Other nodes are not indexed by block height.
	if len(n.SubChainIDs) == 0 {
		marshal := n.Marshal()
		if marshal == nil {
			return errors.New("failed to marshal node")
		}
		db.PutInt32(types.DirectoryBlockHeight, int(n.BHeight), nHash)
	}

	db.Put(types.Node, nHash, n.Marshal()) // And of course, store the actual content.  Only in one place in the DB

	return nil
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
	if n.ListMDRoot != n2.ListMDRoot {
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
	if len(n.EntryList) != len(n2.EntryList) {
		return false
	}
	for i, list := range n.EntryList {
		if list != n2.EntryList[i] {
			return false
		}
	}
	return true
}

// Marshal
// Convert the given entry into a byte slice. Add to that the SubChainIDs of the ChainID. Returns nil
// if anything goes wrong while marshaling
func (n Node) Marshal() (bytes []byte) {

	if n.MarshalCache != nil {
		return n.MarshalCache
	}

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
	bytes = append(bytes, types.Uint16Bytes(uint16(len(n.SubChainIDs)))...) // Put the number of SubChains
	for _, subChain := range n.SubChainIDs {                                // For each SubChain
		bytes = append(bytes, subChain.Bytes()...) // Put the ExtID's data in the slice
	}
	bytes = append(bytes, n.Previous.Bytes()...)
	bytes = append(bytes, types.BoolBytes(n.IsNode)...)
	bytes = append(bytes, n.ListMDRoot.Bytes()...)
	bytes = append(bytes, types.Uint32Bytes(uint32(len(n.List)))...) // Put the number of List Entries
	for _, list := range n.List {                                    // For each SubChain
		bytes = append(bytes, list.ChainID.Bytes()...) // Chain/SubChain ID
		bytes = append(bytes, list.MDRoot.Bytes()...)  // MD of the sub node or entry
	}
	bytes = append(bytes, types.Uint32Bytes(uint32(len(n.EntryList)))...) // Put the number of List Entries
	for _, list := range n.EntryList {                                    // For each Entry
		bytes = append(bytes, list.Bytes()...) // MD of the sub node or entry
	}
	n.MarshalCache = bytes
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

// GetMDRoot
// The MDRoot is the combination of the header on the left, and the listMDRoot on the right.
// Returns a nil if the MDRoot doesn't exist do to missing data
func (n Node) GetMDRoot() (mdRoot *types.Hash) {
	headerHash := n.GetHash()
	if headerHash == nil {
		return nil
	}
	mdr := types.Hash{}
	th := sha256.Sum256(append(headerHash.Bytes(), n.ListMDRoot.Bytes()...))
	mdr.Extract(th[:])
	return &mdr
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
	data = n.ListMDRoot.Extract(data)
	// Pull out all the List entries
	var listLen uint32
	listLen, data = types.BytesUint32(data)
	for i := uint32(0); i < listLen; i++ {
		ne := new(NEList)
		data = ne.ChainID.Extract(data)
		data = ne.MDRoot.Extract(data)
		n.List = append(n.List, *ne)
	}
	var eListLen uint32
	listLen, data = types.BytesUint32(data)
	for i := uint32(0); i < eListLen; i++ {
		var eHash types.Hash
		eHash.Extract(data)
		n.EntryList = append(n.EntryList, eHash)
	}

	return len(d) - len(data), nil // Return the bytes consumed and a nil that all is well for an error

}
