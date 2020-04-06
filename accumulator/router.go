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

func (r *Router) Init() {
	r.Mux.Lock()
	r.chainAccs = make(map[[32]byte]*ChainAcc, 10)
	r.Mux.Unlock()
}

func (r *Router) AddTx(chainId [32]byte, hash []byte) {
	if r.chainAccs[chainId] == nil {
		r.Mux.Lock()
		if r.chainAccs[chainId] == nil {
			chainAcc := new(ChainAcc)
			Txs := make(HashStream, 10000)
			r.chainAccs[chainId] = chainAcc
			r.chainAccs[chainId].txs = Txs

			go chainAcc.acc.Run(Txs)
		}
		r.Mux.Unlock()
	}
	r.chainAccs[chainId].txs <- hash

}

func (r *Router) CloseAll() (chains int64, count int64, pending int64) {
	r.Mux.Lock()
	chains = int64(len(r.chainAccs))
	var keys [][32]byte
	for k, _ := range r.chainAccs {
		var key [32]byte
		copy(key[:], k[:])
		keys = append(keys, k)
	}
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
