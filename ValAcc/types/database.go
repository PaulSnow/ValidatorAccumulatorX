package types

// ======================= Database Support =======================================

// Bucket Names used by the accumulator and validator
const (
	NodeHead             = "node head" // Key: ChainID        Value:  last node hash for this chainID
	Entry                = "entry"     // Key: EntryHash      Value:  Entry
	DirectoryBlockHeight = "dbHeight"  // Key: Height uint16  Value:  Directory Block node
	DirectoryBlockHash   = "dbHash"    // Key: DBHash         Value:  Directory Block node
	Node                 = "node"      // Key: nodeHash       Value:  nodeHash
)
