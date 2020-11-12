package util

import "strconv"

//ParseUint8 parse a uint8 from string
func ParseUint8(s string) uint8 {
	i, _ := strconv.ParseUint(s, 10, 8)
	return uint8(i)
}

//ParseUint16 parse a uint16 from string
func ParseUint16(s string) uint16 {
	i, _ := strconv.ParseUint(s, 10, 16)
	return uint16(i)
}

//ParseUint64 parse a uint64 from string
func ParseUint64(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint64(i)
}
