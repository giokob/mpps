package utils

func CalculateSum(data []byte) (byte, int) {
	maxIndex := len(data) - 1
	c := byte(0)
	for i := 0; i < maxIndex; i++ {
		c = c ^ data[i]
	}
	return c, maxIndex
}

func CheckSum(data []byte) bool {
	c, maxIndex := CalculateSum(data)
	return (c == data[maxIndex])
}

func AddCheckSum(data []byte) {
	c, maxIndex := CalculateSum(data)
	data[maxIndex] = c
}

func AddSequence(data []byte, s int) {
	data[0] = byte((s >> 24) & 255)
	data[1] = byte((s >> 16) & 255)
	data[2] = byte((s >> 8) & 255)
	data[3] = byte(s & 255)
}

func ReadSequence(data []byte) int {
	var seq int = 0
	seq = int(data[0]) << 24
	seq |= int(data[1]) << 16
	seq |= int(data[2]) << 8
	seq |= int(data[3])
	return seq
}
