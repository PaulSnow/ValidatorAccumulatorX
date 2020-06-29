package merkleDag

import (
	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"
)

type ReceiptNode struct {
	Right bool
	Hash  types.Hash // Left Hash of a Merkle Directed Acyclic Graph (MD)
}

type MDReceipt struct {
	EntryHash types.Hash    // Entry Hash of the data subject to the MDReceipt
	Nodes     []ReceiptNode // Path through the data collected by the MerkleDag
	MDRoot    types.Hash    // Merkle DAG root from the Accumulator network.
	// We likely want a struct here provided by the underlying blockchain where we are recording
	// the MDRoots for the Accumulator
}

// BuildMDReceipt
// Building a receipt is a bit more complex than just computing the Merkle DAG.  We look at the stream of
// hashes as we build the MD, and detect when the Hash for which we need a receipt shows up.  Then we
// track it through the intermediate hashes, tracking if it is combined from the right or from the left.
// Then when we have to calculate the Merkle DAG root, we do one more pass through the combining of the
// trailing hashes to finish off the receipt.
func (mdr *MDReceipt) BuildMDReceipt(MerkleDag MD, data types.Hash) {
	mdr.Nodes = mdr.Nodes[:0] // Throw away any old paths
	mdr.EntryHash = data      // The Data for which this is a Receipt
	md := []*types.Hash{nil}  // The intermediate hashes used to compute the Merkle DAG root
	right := true             // We assume we will be combining from the right
	idx := -1                 // idx of -1 means not yet found the hash for which we want a receipt in the hash stream

DataLoop: // Loop through the data behind the Merkle DAG
	for _, h := range MerkleDag.HashList {
		// Always make sure md ends with a nil; limits corner cases
		if md[len(md)-1] != nil { // Make sure md still ends in a nil
			md = append(md, nil)
		}
		// Look for the data that we are computing a receipt for
		if idx < 0 && h == data {
			idx = 0
		}
		// Then add this data to the Merkle DAG we are creating
		for i, v := range md {
			if v == nil {
				if i == idx { // If we move from right to left, set right to false
					right = false
				}
				md[i] = h.Copy()
				continue DataLoop
			}
			if i == idx { // If we are on the path, and we are going to combine, then record
				rn := new(ReceiptNode)
				mdr.Nodes = append(mdr.Nodes, *rn)
				rn.Right = right
				if right { // If our hash is on the right, we need to record the hash on the left
					rn.Hash = *v.Copy()
				} else { // If our hash is on the left, we need to record the hash on the right
					rn.Hash = h
				}
				right = true // Regardless, "our" hash is left on the right
				idx++        // And our hash progresses down the table
			}
			h = *v.Combine(h) // Combine v (left) and hash (right) to get a new combined hash to use forward
			md[i] = nil       // Now that we have combined v and hash, this spot is now empty, so clear it.
		}
	}

	if idx == -1 {
		mdr.Nodes = mdr.Nodes[:0]
		return
	}
	var mdroot *types.Hash
	right = false

	// We close the Merkle DAG
	for i, v := range md {
		if i == idx && mdroot == nil {
			right = true
		}
		if mdroot == nil { // We will pick up the first hash in m.MD no matter what.
			mdroot = v // If we assign a nil over a nil, no harm no foul.  Fewer cases to test this way.
		} else if v != nil { // If MDRoot isn't nil and v isn't nil, we combine them.

			if i >= idx {
				rn := new(ReceiptNode)
				mdr.Nodes = append(mdr.Nodes, *rn)
				rn.Right = right
				if right {
					rn.Hash = *v.Copy()
					rn.Right = false
				} else {
					rn.Hash = *mdroot.Copy()
					right = false
				}
				idx++
			}
			mdroot = v.Combine(*mdroot) // v is on the left, MDRoot candidate is on the right, for a new MDRoot
			mdr.MDRoot = *mdroot.Copy() // The last one is the one we want
		}
	}
	return
}

// Validate
// Run down the Merkle DAG and prove that this receipt self validates
func (mdr MDReceipt) Validate() bool {
	sum := mdr.EntryHash
	for _, n := range mdr.Nodes {
		hash := n.Hash.Copy()
		if n.Right {
			sum = *hash.Combine(sum)
		} else {
			sum = *sum.Combine(*hash)
		}
	}
	return sum == mdr.MDRoot
}
