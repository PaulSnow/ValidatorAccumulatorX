package types

type DataField []byte // Typing the Data Fields to allow flexibility in representation in the future

func (d *DataField) Bytes() []byte {
	return *d
}

func (d *DataField) Copy() []byte {
	return append([]byte{}, *d...)
}

func (d *DataField) Extract(len uint16, data []byte) []byte {
	*d = append([]byte{}, data[0:int(len)]...)
	return data[len:]
}
