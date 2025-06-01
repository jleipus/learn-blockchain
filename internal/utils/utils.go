package utils

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/exp/constraints"
)

// IntToHex converts an integer to a byte slice in big-endian format.
func IntToHex[T constraints.Integer](n T) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// ReverseBytes reverses a byte array.
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
