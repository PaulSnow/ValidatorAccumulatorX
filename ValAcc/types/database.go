package types

// ======================= Database Support =======================================
// Bucket Names used by the accumulator and validator
const (
	Version              = VersionField(0)          // Version of ValAcc
	NodeFirst            = "first node"             // Key: node.ChainID      Value:  First node hash with this chainID
	NodeNext             = "next node"              // Key: node.GetHash()    Value:  next node in sequence with this chainID
	NodeHead             = "node head"              // Key: node.ChainID      Value:  last node hash for this chainID
	Entry                = "entry"                  // Key: entry.GetHash()   Value:  Entry
	EntryNode            = "entry Node"             // Key: entry.GetHash()   Value:  node where this entry is recorded
	DirectoryBlockHeight = "directory block height" // Key: node.BHeight      Value:  Directory Block node
	Node                 = "node"                   // Key: node.GetHash()    Value:  nodeHash
)
