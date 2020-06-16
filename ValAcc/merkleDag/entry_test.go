package merkleDag

import (
	"crypto/sha256"
	"testing"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

func TestEntry(t *testing.T) {
	e := new(Entry)
	e.Version = 1

	subChainID1 := sha256.Sum256([]byte("a hash 1"))
	subChainID2 := sha256.Sum256([]byte("a hash 2"))
	subChainID3 := sha256.Sum256([]byte("a hash 3"))
	subChainID4 := sha256.Sum256([]byte("a hash 4"))
	e.SubChainIDs = append(e.SubChainIDs, subChainID1)
	e.SubChainIDs = append(e.SubChainIDs, subChainID2)
	e.SubChainIDs = append(e.SubChainIDs, subChainID3)
	e.SubChainIDs = append(e.SubChainIDs, subChainID4)

	chainData := append([]byte{}, subChainID1[:]...)
	chainData = append(chainData, subChainID2[:]...)
	chainData = append(chainData, subChainID3[:]...)
	chainData = append(chainData, subChainID4[:]...)
	chainID := sha256.Sum256(chainData)

	copy(e.ChainID[:], chainID[:])
	e.ExtIDs = append(e.ExtIDs, types.DataField([]byte("Field 1")))
	e.ExtIDs = append(e.ExtIDs, types.DataField([]byte("Field 2")))
	e.ExtIDs = append(e.ExtIDs, types.DataField([]byte("Field 3")))
	e.ExtIDs = append(e.ExtIDs, types.DataField([]byte("Field 4")))
	e.ExtIDs = append(e.ExtIDs, types.DataField([]byte("Field 5")))
	e.Content = []byte{}
	entrySlice := e.Marshal()
	if entrySlice == nil {
		t.Error("Failed to marshal an Entry")
	}
	var e2 Entry
	len, err := e2.Unmarshal(entrySlice)
	if err != nil {
		t.Error("Failed to unmarshal an Entry")
	}
	if !e.SameAs(e2) {
		t.Error("Did not unmarshal an Entry as expected")
	}
	if len != 212 {
		t.Errorf("Length of data consumed (%d) not as expected (%d)", len, 82)
	}
}
