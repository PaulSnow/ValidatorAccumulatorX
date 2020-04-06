// Router keeps track of all the chains in flight.  It routes transactions to the right chain.
//
package accumulator

import (
	"sync"
)

type ChainAcc struct {
	acc Chain
	txs chan []byte
}

type Router struct {
	chainAccs map[[32]byte]*ChainAcc
	Mux       sync.Mutex
}

// We create an initial Map to track the chains we are writing to.
func (r *Router) Init() {
	r.Mux.Lock()
	r.chainAccs = make(map[[32]byte]*ChainAcc, 10)
	r.Mux.Unlock()
}

// AddTx
// Add a TX to its chain.  Allocates a chain if one does not exist.
func (r *Router) AddTx(chainId [32]byte, hash []byte) {
	//Check if we have a chain.
	if r.chainAccs[chainId] == nil {
		// Lock the router
		r.Mux.Lock()
		// Make sure we still don't have a chain.
		if r.chainAccs[chainId] == nil {
			// Allocate a new accumulator for a new chain
			chainAcc := new(ChainAcc)
			Txs := make(HashStream, 10000)
			r.chainAccs[chainId] = chainAcc
			r.chainAccs[chainId].txs = Txs
			// set the accumulator in motion
			go chainAcc.acc.Run(Txs)
		}
		// Unlock the router
		r.Mux.Unlock()
	}
	// Send the transaction hash to the accumulator for its chain
	r.chainAccs[chainId].txs <- hash

}

// CloseAll
// Close all the accumulators and sum up the results, and pass results back to caller.
func (r *Router) CloseAll() (chains int64, count int64, pending int64) {
	// Lock the router
	r.Mux.Lock()
	chains = int64(len(r.chainAccs))
	var keys [][32]byte
	// Get all the keys.  Was getting map conflicts.  This should not be necessary
	for k, _ := range r.chainAccs {
		var key [32]byte
		copy(key[:], k[:])
		keys = append(keys, k)
	}
	// Go through all the chains, and close them.  Collect their state for reporting.
	for _, k := range keys {
		c := r.chainAccs[k]
		c.acc.Mux.Lock()
		count += c.acc.Count
		pending += int64(len(c.txs))
		c.acc.CloseMR()
	}
	r.Mux.Unlock()
	return
}
