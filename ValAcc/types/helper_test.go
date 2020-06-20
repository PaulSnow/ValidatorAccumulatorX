package types

import (
	"fmt"
	"testing"
)

func TestHelperFunc(t *testing.T) {
	{
		v1 := uint32(0x11223344)
		b := Uint32Bytes(v1)
		v2 := uint32(0x55667788)
		b = append(b, Uint32Bytes(v2)...)
		fmt.Printf("b:  %x\n", b)
		r1, b := BytesUint32(b)
		fmt.Printf("r1: %x\n", r1)
		fmt.Printf("b:  %x\n", b)
		r2, b := BytesUint32(b)
		fmt.Printf("r2: %x\n", r2)
		fmt.Printf("b:  %x\n", b)
		if v1 != r1 {
			t.Error("v1, r1 didn't match")
		}
		if v2 != r2 {
			t.Error("v2 r2 didn't match")
		}
	}
	{
		v1 := uint16(0x1122)
		b := Uint16Bytes(v1)
		v2 := uint16(0x5566)
		b = append(b, Uint16Bytes(v2)...)
		fmt.Printf("b:  %x\n", b)
		r1, b := BytesUint16(b)
		fmt.Printf("r1: %x\n", r1)
		fmt.Printf("b:  %x\n", b)
		r2, b := BytesUint16(b)
		fmt.Printf("r2: %x\n", r2)
		fmt.Printf("b:  %x\n", b)
		if v1 != r1 {
			t.Error("v1, r1 didn't match")
		}
		if v2 != r2 {
			t.Error("v2 r2 didn't match")
		}
	}
}
