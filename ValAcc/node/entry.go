package node

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/AccumulateNetwork/ValidatorAccumulator/ValAcc/types"
)

// ANode
// Data is entered into system by the Accumulator as a series of entries organized by chainIDs
// Unlike Factom, we will attempt to create a chain if the ChainID provided is nil.  We provide
// a function to compute the ChainID from the first entry in a chain, for use by applications.
// If a chain already exists, creating the chain will be ignored.
//
// ANode
//    Version          uint8
//    TimeStamp        uint64
//    ChainID          [32]byte
//    len(SubChains)   uint64
//    SubChains        []SubChain
//       Len(SubChain)    uint16
//       SubChain[]       []byte
//    #ExtIDs          uint16
//    ExtIDs           []ExtID
//      len(ExtID)        uint16
//      ExtID             []byte
//    len(content)     uint16
//    Content          []byte
type ANode struct {
	Version     types.VersionField // Version of this data structure
	TimeStamp   types.TimeStamp    // Timestamp of the construction of this entry
	ChainID     types.Hash         // The ChainID
	SubChainIDs []types.Hash       // SubChainIDs required to build the ChainID
	ExtIDs      []types.DataField  // External ids used to create the chain id above ( see ExternalIDsToChainID() )
	Content     types.DataField    // BytesSlice for holding generic data for this entry
}

// SameAs
func (e ANode) SameAs(e2 ANode) bool {
	if e.Version != e2.Version {
		return false
	}
	if e.ChainID != e2.ChainID {
		return false
	}
	if len(e.SubChainIDs) != len(e2.SubChainIDs) {
		return false
	}
	for i := range e.SubChainIDs {
		if e.SubChainIDs[i] != e2.SubChainIDs[i] {
			return false
		}
	}
	if len(e.ExtIDs) != len(e2.ExtIDs) {
		return false
	}
	for i := range e.ExtIDs {
		if !bytes.Equal(e.ExtIDs[i], e2.ExtIDs[i]) {
			return false
		}
	}
	if !bytes.Equal(e.Content, e2.Content) {
		return false
	}
	return true
}

// Marshal
// Convert the given entry into a byte slice. Add to that the SubChainIDs of the ChainID. Returns nil
// if anything goes wrong while marshaling
func (e ANode) Marshal() (bytes []byte) {

	// On any error, return a nil for the byte representation of the ANode
	defer func() {
		if r := recover(); r != nil {
			bytes = nil
			return
		}
	}()

	bytes = append(bytes, e.Version.Bytes()...)   // Put the version into the slice
	bytes = append(bytes, e.TimeStamp.Bytes()...) // Add the TimeStamp
	bytes = append(bytes, e.ChainID.Bytes()...)   // Put the ChainID into the slice

	bytes = append(bytes, types.Uint16Bytes(uint16(len(e.SubChainIDs)))...) // Put the number of ExtIDs in the slice
	for _, subChain := range e.SubChainIDs {                                // For each ExtID
		bytes = append(bytes, subChain.Bytes()...) // Put the ExtID's data in the slice
	}

	bytes = append(bytes, types.Uint16Bytes(uint16(len(e.ExtIDs)))...) // Put the number of ExtIDs in the slice
	for _, extID := range e.ExtIDs {                                   // For each ExtID
		bytes = append(bytes, types.Uint16Bytes(uint16(len(extID)))...) // Put its length in the slice
		bytes = append(bytes, extID.Bytes()...)                         // Put the ExtID's data in the slice
	}
	if e.Content != nil { // If we have content
		bytes = append(bytes, types.Uint16Bytes(uint16(len(e.Content)))...) // Put the content length in the slice
		bytes = append(bytes, e.Content.Bytes()...)                         // Put the content in the slice
	} else { // If we don't have content
		bytes = append(bytes, 0, 0) // If no content, put a uint16 0 down for its length
	}

	return bytes // Return the slice
}

// GetHash
// Returns the EntryHash for this entry.  Note the ANode Hash does not include the SubChainIDs
func (e ANode) GetHash() (hash *types.Hash) {
	h := e.Marshal() // Get the bytes behind the EntryHash
	if h == nil {    // A nil would mean the ANode didn't marshal
		return nil
	} // If Marshal Fails, return a nil
	hash = new(types.Hash) // Get the Hash object to return
	hs := sha256.Sum256(h) // Get the array holding the hash (so we can create a slice)
	hash.Extract(hs[:])    // Populate the Hash object
	return hash            // Return a pointer to the Hash object.
}

// Unmarshal
// Extract an entry from a byte slice.  Returns an error if the unmarshal fails, or the length of the
// data consumed and a nil.
func (e *ANode) Unmarshal(data []byte) (dataConsumed int, err error) {

	// On any error, no data is consumed and return an error as to why unmarshal fails
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ANode Failed to unmarshal %v", r))
		}
	}()
	d := data                          // d keeps the original slice
	data = e.Version.Extract(data)     // Extract the version
	data = e.TimeStamp.Extract(data)   // Extract the TimeStamp
	data = e.ChainID.Extract(data)     // Extract the ChainID
	e.SubChainIDs = e.SubChainIDs[0:0] // Clear any ExtIDs that might already be in this ANode

	// Pull out all the subChain IDs
	var numSubChains uint16
	numSubChains, data = types.BytesUint16(data) // Get the number of SubChainIDs we should have
	for i := uint16(0); i < numSubChains; i++ {  // Pull each of them out of the data slice
		sc := types.Hash{}                        // Get a Hash to put the SubChainID in
		data = sc.Extract(data)                   // Extract the ExtID
		e.SubChainIDs = append(e.SubChainIDs, sc) // Put it in the ExtID list
	}

	// Pull out all the Extended IDs
	var numExtIDs uint16
	numExtIDs, data = types.BytesUint16(data) // Get the number of ExtIDs we should have
	for i := uint16(0); i < numExtIDs; i++ {  // Pull each of them out of the data slice
		var lExt uint16
		lExt, data = types.BytesUint16(data) // Each are lead by a length
		ext := types.DataField{}             // Get a DataField to put the ExtID in
		data = ext.Extract(lExt, data)       // Extract the ExtID
		e.ExtIDs = append(e.ExtIDs, ext)     // Put it in the ExtID list
	}

	lContent, data := types.BytesUint16(data) // Get the length of the content
	e.Content.Extract(lContent, data)         //Extract the content

	return len(d) - len(data), nil // Return the bytes consumed and a nil that all is well for an error

}

// EntryHash
// The accumulator assumes the ANode has already been written to the (a) database.  It is only dealing with EntryHashes
type EntryHash struct {
	SubChains []types.Hash // SubChainIDs used to create the ChainID
	ChainID   types.Hash   // The ChainID
	EntryHash types.Hash   // The EntryHash
}
