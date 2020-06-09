package accumulator

import "github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"

// NodeBlock
// A Node Block collects a set of Merkle DAGs and creates a MD from that.  At the root, the NodeBlock is
// called the Directory Block in Factom.  But as this is a much higher volume architecture, we don't limit ourselves
// to Directory Blocks and Entry Blocks, but use a more general structure of nodes.  At the leaf level, a node
// is pretty much an entry block, and the previous Hash links back to the previous entry block for a chain.  Sequence
// numbers count from 1 (the first entry block in a chain) and up.
//
// If the SequenceNumber is zero and the Previous Hash is nil, this is an intermediate node, and all chains are
// equal to the given ChainID, and less than the next NodeBlock's ChainID.
// covering particular ranges of chains
type NodeBlock struct {
	ChainID        *types.Hash
	SequenceNumber int32
	Height         uint32
	Timestamp      uint64
	MD             types.Hash
	Previous       types.Hash
}
