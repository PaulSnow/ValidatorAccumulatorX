package types

type BlockHeight int32 // We are typing certain things in the protocol

func (bh BlockHeight) Bytes() []byte {
	return append([]byte{},
		byte(bh>>56), byte(bh>>48), byte(bh>>40), byte(bh>>32),
		byte(bh>>24), byte(bh>>16), byte(bh>>8), byte(bh))
}

func (bh *BlockHeight) Extract(data []byte) []byte {
	*bh = BlockHeight(((((((int64(data[0])<<8+int64(data[1]))<<8+int64(data[2]))<<8+int64(data[3]))<<8+
		int64(data[4]))<<8+int64(data[5]))<<8+int64(data[6]))<<8 + int64(data[7]))
	return data[8:]
}
