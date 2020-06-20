package types

type BlockHeight uint32 // We are typing certain things in the protocol

func (bh BlockHeight) Bytes() []byte {
	return Uint32Bytes(uint32(bh))
}

func (bh *BlockHeight) Extract(data []byte) []byte {
	bhv, newData := BytesUint32(data)
	*bh = BlockHeight(bhv)
	return newData
}
