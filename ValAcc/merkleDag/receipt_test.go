package merkleDag

import (
	"crypto/sha256"
	"testing"
	"time"
)

func TestMD(t *testing.T) {

	start := time.Now()
	_ = start
	startHash := sha256.Sum256([]byte("test hash"))

	for limit := 10; limit < 17000; limit += limit<<17 ^ limit>>8%100 {
		md := new(MD)
		if limit%100 == 0 {
			println("limit ", limit)
		}
		for i := 0; i < limit; i++ {
			md.AddToChain(startHash)
			startHash = sha256.Sum256(startHash[:])
		}
		mdroot := md.GetMDRoot()
		//fmt.Println(" Took ", time.Now().Sub(start))
		//fmt.Printf(" Merkle DAG Root %x\n", *mdroot)

		MDR := new(MDReceipt)
		MDR.BuildMDReceipt(*md, md.HashList[0])
		if !MDR.Validate() {
			t.Errorf("Receipt fails to validate (%d)", limit)
		}
		//fmt.Printf(" Merkle DAG Root %x\n", MDR.MDRoot)
		if MDR.MDRoot != *mdroot {
			t.Errorf("%d Merkle Roots not equal %x %x", limit, MDR.MDRoot, *mdroot)
		}
	}
}
