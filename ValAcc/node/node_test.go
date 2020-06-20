package node

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/database"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

// GetTestNode
// Build a node for use in tests
func GetTestNode(t *testing.T) *Node {
	n := new(Node)
	n.Version = types.Version
	n.BHeight = 1
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
	n.ChainID = types.GetChainID(AccDID, n.SubChainIDs)

	var expectedChainID [32]byte
	_, err := hex.Decode(expectedChainID[:], []byte(expected))
	if err != nil || expectedChainID != n.ChainID {
		t.Errorf("Didn't get the expected ChainID. Got %x Expected %x", n.ChainID, expectedChainID)
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
	return n
}

func TestNode(t *testing.T) {
	n := GetTestNode(t)

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
	expectedLen := 444
	if nodeLen != expectedLen {
		t.Errorf("Length of data consumed (%d) not as expected (%d)", nodeLen, expectedLen)
	}
}

// GetTestDB
// Helper function for other tests to get a Test DB for running tests against a physical database
func GetTestDB(t *testing.T) *database.DB {
	dName, e := ioutil.TempDir("", "sampleDir")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(dName)
	db := new(database.DB)
	db.DBHome = dName
	db.Init(0)
	return db
}

func TestNodeDB(t *testing.T) {
	db := GetTestDB(t)
	node := GetTestNode(t)                              // Get a node, but we are going to override a bunch of fields.
	node.SubChainIDs = node.SubChainIDs[:0]             // Clear out the subChainIDs, so this is a Directory Block
	node.ChainID = sha256.Sum256([]byte("TestAcc DID")) // Set the Chain ID to a plausible DID
	node.SequenceNum = 0                                // Gotta be zero for a Directory Block
	node.Put(db)

	hash := (*node.GetHash())[:]

	headHash := db.Get(types.NodeHead, node.ChainID[:]) // Should have a node head
	if !bytes.Equal(headHash, hash) {
		t.Error("could not find the Head node for the directory blocks")
	}

	firstHash := db.Get(types.NodeFirst, node.ChainID[:]) // Should have a first node
	if !bytes.Equal(firstHash, hash) {
		t.Error("could not find the first node for the chainID")
	}

	nextHash := db.Get(types.NodeNext, node.ChainID[:]) // There should be no next node yet
	if nextHash != nil {
		t.Error("should not have a next node for the chainID yet.")
	}

	// Check that the DirectoryBlockHeight has the hash of the node
	nodeHash := db.GetInt32(types.DirectoryBlockHeight, 1)
	if !bytes.Equal((*node.GetHash())[:], nodeHash) {
		fmt.Printf("Node\n%x\n", *node.GetHash())
		fmt.Printf("DB  \n%x\n", nodeHash)
		t.Error("the node written to DB != to node read from DB (DirectoryBlockHeight)")
	}
	nodeBytes1 := db.Get(types.Node, nodeHash)
	nodeBytes2 := db.Get(types.Node, (*node.GetHash())[:])
	nodeBytes3 := node.Marshal()
	if !bytes.Equal(nodeBytes1, nodeBytes2) || !bytes.Equal(nodeBytes2, nodeBytes3) {
		t.Error("bytes for node should be the same, if from database node bucket, " +
			"or using hash from DirectoryBlockHeight, or just marshalling the node")
	}

}
