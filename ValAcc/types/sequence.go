package types

type Sequence int32 // We are typing certain things in the protocol

func (s Sequence) Bytes() []byte {
	return append([]byte{},
		byte(s>>56), byte(s>>48), byte(s>>40), byte(s>>32),
		byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

func (s *Sequence) Extract(data []byte) []byte {
	*s = Sequence(((((((int64(data[0])<<8+int64(data[1]))<<8+int64(data[2]))<<8+int64(data[3]))<<8+
		int64(data[4]))<<8+int64(data[5]))<<8+int64(data[6]))<<8 + int64(data[7]))
	return data[8:]
}
