package types

type Sequence uint32 // We are typing certain things in the protocol

func (s Sequence) Bytes() []byte {
	return Uint32Bytes(uint32(s))
}

func (s *Sequence) Extract(data []byte) []byte {
	ss, newData := BytesUint32(data)
	*s = Sequence(ss)
	return newData
}
