package types

type VersionField uint8 // We are typing certain things in the protocol

func (v VersionField) Bytes() []byte {
	return append([]byte{}, byte(v))
}

func (v *VersionField) Extract(data []byte) []byte {
	*v = VersionField(data[0])
	return data[1:]
}
