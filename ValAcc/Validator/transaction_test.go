// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"
)

func TestCreateKey(t *testing.T) {
	cnt := 1000000
	start := time.Now()
	val := sha256.Sum256([]byte("start"))
	for i := 0; i < cnt; i++ {
		ed25519.NewKeyFromSeed(val[:])
		val = sha256.Sum256(val[:])
	}
	ts := time.Now().Sub(start).Seconds()
	fmt.Println("Took ", ts, " seconds to create ", cnt, " new keys, or ", cnt/int(ts), "per second")
}

func TestSignData(t *testing.T) {
	cnt := 1000000
	start := time.Now()
	val := sha256.Sum256([]byte("start"))
	key := ed25519.NewKeyFromSeed(val[:])
	for i := 0; i < cnt; i++ {
		ed25519.Sign(key, val[:])
		val = sha256.Sum256(val[:])
	}
	ts := time.Now().Sub(start).Seconds()
	fmt.Println("Took ", ts, " seconds to sign ", cnt, " data, or ", cnt/int(ts), "per second")
}

func TestValidatesig(t *testing.T) {
	cnt := 1000000
	p := 8
	start := time.Now()
	val := sha256.Sum256([]byte("start"))
	pub, pri, _ := ed25519.GenerateKey(rand.Reader)
	sig := ed25519.Sign(pri, val[:])
	cs := make(chan int, p)
	for i := 0; i < p; i++ {
		go func() {
			for i := 0; i < cnt/p; i++ {
				ed25519.Verify(pub, val[:], sig)
			}
			cs <- 1
		}()
	}
	for i := 0; i < p; i++ {
		<-cs
	}
	ts := time.Now().Sub(start).Seconds()
	fmt.Println("Took ", ts, " seconds to validate ", cnt, " sigs, or ", cnt/int(ts), "per second")
}
