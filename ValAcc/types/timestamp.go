package types

import "time"

type TimeStamp int64

//
func GetCurrentTimeStamp() TimeStamp {
	return TimeStamp(time.Now().Unix())
}

// Extract a timestamp from a byte slice.  Return the updated byte slice.
func (t *TimeStamp) Extract(data []byte) (newData []byte) {
	*t = TimeStamp(((((((int64(data[0])<<8+int64(data[1]))<<8+int64(data[2]))<<8+int64(data[3]))<<8+
		int64(data[4]))<<8+int64(data[5]))<<8+int64(data[6]))<<8 + int64(data[7]))
	return data[8:]
}

func (t TimeStamp) Bytes() []byte {
	return append([]byte{},
		byte(t>>56), byte(t>>48), byte(t>>40), byte(t>>32),
		byte(t>>24), byte(t>>16), byte(t>>8), byte(t))
}
