package accumulator

import (
	"time"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/database"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/merkleDag"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/node"
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

// ChainAcc
// Tracks the construction of the Merkle DAG and collects the Hash sequence to build the MD
type ChainAcc struct {
	Node node.Node
	MD   *merkleDag.MD // The class for creating the MD and MD Roots
}

func NewChainAcc(DB database.DB, eHash node.EntryHash, bHeight types.BlockHeight) *ChainAcc {
	chainAcc := new(ChainAcc)

	previousHash := DB.Get(types.NodeHead, eHash.ChainID[:])
	if previousHash != nil {
		previousBytes := DB.Get(types.Node, previousHash[:])
		var previous node.Node
		previous.Unmarshal(previousBytes)
		chainAcc.Node.SequenceNum = previous.SequenceNum
		chainAcc.Node.MarshalCache = previousBytes
		chainAcc.Node.Previous = *previous.GetHash()
	}
	chainAcc.Node.Version = types.Version
	chainAcc.Node.SubChainIDs = eHash.SubChains
	chainAcc.Node.ChainID = eHash.ChainID
	chainAcc.Node.TimeStamp = types.TimeStamp(time.Now().UnixNano())
	chainAcc.Node.BHeight = bHeight
	chainAcc.Node.IsNode = false
	chainAcc.MD = new(merkleDag.MD)
	return chainAcc
}
