package main

import (
	"crypto/sha256"
)

var HashList [][]byte
var MR       [][]byte

func CreateMR (start int, hash []byte) {
	i:=start
		if len(MR)==i {
			MR = append(MR,hash)
			return
		}
		if MR[i]==nil {
			MR[i]=hash
			return
		}
		h := sha256.New()
		h.Write(MR[i][:])
		h.Write(hash)
		CreateMR(i+1,h.Sum(nil))
}



func main() {
	
}
