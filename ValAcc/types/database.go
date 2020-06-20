package types

// ======================= Database Support =======================================
// Bucket Names used by the accumulator and validator
const (
	NodeHead             = "node head"              // Key: node.ChainID      Value:  last node hash for this chainID
	Entry                = "entry"                  // Key: entry.GetHash()   Value:  Entry
	DirectoryBlockHeight = "directory block height" // Key: node.BHeight      Value:  Directory Block node
	Node                 = "node"                   // Key: node.GetHash()    Value:  nodeHash
)
