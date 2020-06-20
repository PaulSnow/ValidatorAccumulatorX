package node

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

func TestEntry(t *testing.T) {
	e := new(Entry)
	e.Version = 1
	e.TimeStamp = types.TimeStamp(time.Now().Unix())

	AccDID := sha256.Sum256([]byte("TestAcc"))

	subChainID1 := sha256.Sum256([]byte("a hash 1"))
	subChainID2 := sha256.Sum256([]byte("a hash 2"))
	subChainID3 := sha256.Sum256([]byte("a hash 3"))
	subChainID4 := sha256.Sum256([]byte("a hash 4"))
	e.SubChainIDs = append(e.SubChainIDs, subChainID1)
	e.SubChainIDs = append(e.SubChainIDs, subChainID2)
	e.SubChainIDs = append(e.SubChainIDs, subChainID3)
	e.SubChainIDs = append(e.SubChainIDs, subChainID4)

	expected := "9e4961b2d1d600a59494830888c4b2085467778610d621ac008097d5ba05b866"
	chainID := types.GetChainID(AccDID, e.SubChainIDs)

	var expectedChainID [32]byte
	_, err := hex.Decode(expectedChainID[:], []byte(expected))
	if err != nil || !bytes.Equal(expectedChainID[:], chainID[:]) {
		t.Errorf("Didn't get the expected ChainID. Got %x Expected %x", chainID, expectedChainID)
	}

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
	expectedLen := 220
	if len != 220 {
		t.Errorf("Length of data consumed (%d) not as expected (%d)", len, expectedLen)
	}
}
