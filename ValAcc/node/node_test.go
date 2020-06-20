package node

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

/*
type Node struct {
	Version     types.VersionField // Version of this data structure
	BHeight     types.BlockHeight  // Block Height
    SequenceNum
	TimeStamp   types.TimeStamp    // TimeStamp by Accumulator when the structure was built
	ChainID     types.Hash         // The ChainID (Directory Block - zeros, SubNodes - 1st ChainID, Entries - ChainID)
	SubChainIDs []types.Hash       // SubChainIDs to build ChainID, Directory Block - zeros
	Previous    types.Hash         // Hash of previous Block Header
	IsNode      bool               // IsNode is true for a node, is false for an entry
	ListMDRoot      types.Hash         // Merkle DAG of the entries of the List (only the hashes)
	List        []NEList           // List of ChainIDs/MDRoots for nodes/entries
}
*/
func TestNode(t *testing.T) {
	n := new(Node)
	n.Version = 1
	n.BHeight = 232433
	n.SequenceNum = 1392
	n.TimeStamp = types.TimeStamp(time.Now().Unix())

	AccDID := sha256.Sum256([]byte("TestAcc"))

	subChainID1 := sha256.Sum256([]byte("FirstChainID 1"))
	subChainID2 := sha256.Sum256([]byte("FirstChainID 2"))
	subChainID3 := sha256.Sum256([]byte("FirstChainID 3"))
	subChainID4 := sha256.Sum256([]byte("FirstChainID 4"))
	n.SubChainIDs = append(n.SubChainIDs, subChainID1)
	n.SubChainIDs = append(n.SubChainIDs, subChainID2)
	n.SubChainIDs = append(n.SubChainIDs, subChainID3)
	n.SubChainIDs = append(n.SubChainIDs, subChainID4)

	expected := "dc36bf13e36984b08d12e66a3f1d1518d380977fb5f1e1814282a3347de3dc96"
	chainID := types.GetChainID(AccDID, n.SubChainIDs)
	copy(n.ChainID[:], chainID[:])

	var expectedChainID [32]byte
	_, err := hex.Decode(expectedChainID[:], []byte(expected))
	if err != nil || !bytes.Equal(expectedChainID[:], chainID[:]) {
		t.Errorf("Didn't get the expected ChainID. Got %x Expected %x", chainID, expectedChainID)
	}
	n.Previous = sha256.Sum256([]byte("Hash of Previous Node"))
	n.IsNode = true
	newNE := func(i int) (ne NEList) {
		ne.ChainID = sha256.Sum256([]byte(fmt.Sprint("list item", i)))
		ne.MDRoot = sha256.Sum256([]byte(fmt.Sprint("MDRoot ", i)))
		return ne
	}
	n.List = append(n.List, newNE(1))
	n.List = append(n.List, newNE(2))
	n.List = append(n.List, newNE(3))

	nodeSlice := n.Marshal()
	if nodeSlice == nil {
		t.Error("Failed to marshal an ANode")
	}
	var n2 Node
	nodeLen, err := n2.Unmarshal(nodeSlice)
	if err != nil {
		t.Error("Failed to unmarshal an ANode")
	}
	if !n.SameAs(n2) {
		t.Error("Did not unmarshal an ANode as expected")
	}
	expectedLen := 440
	if nodeLen != expectedLen {
		t.Errorf("Length of data consumed (%d) not as expected (%d)", nodeLen, expectedLen)
	}
}
