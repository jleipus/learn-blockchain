package utils

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/exp/constraints"
)

func IntToHex[T constraints.Integer](n T) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
