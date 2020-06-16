package accumulator

import (
	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/merkleDag"
	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

// ChainAcc
// Tracks the construction of the Merkle DAG and collects the Hash sequence to build the MD
type ChainAcc struct {
	ChainID  types.Hash    // ChainID for the entries that contribute to MD
	Previous types.Hash    // The previous ChainAcc for this ChainID
	MD       *merkleDag.MD // The class for creating the MD and MD Roots
}

func NewChainAcc() *ChainAcc {
	chainAcc := new(ChainAcc)
	chainAcc.MD = new(merkleDag.MD)
	return chainAcc
}
