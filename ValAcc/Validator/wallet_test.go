// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateTransactions(t *testing.T) {
	wallet := NewWallet([]byte("my test"))

	var txs []*Transaction
	start := time.Now()
	numTxs := 100000
	for i := 0; i < numTxs; i++ {
		tx := new(Transaction)
		tx.input = wallet.NewAddress()
		tx.output = tx.input
		tx.amount = 1 * 10e8
		tx.Sign(wallet.MyKeys[tx.input].PrivateKey)
		txs = append(txs, tx)
	}
	_ = txs
	fmt.Printf("Transactions created per second: %5.2f \n", float64(numTxs)/float64(time.Now().Sub(start).Seconds()))

	start = time.Now()
	for _, tx := range txs {
		if !tx.ValidateSig() {
			t.Error("Validate failed")
		}
	}
	fmt.Printf("Transactions verified per second: %5.2f \n", float64(numTxs)/float64(time.Now().Sub(start).Seconds()))
}
