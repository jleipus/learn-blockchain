package utils

import (
	"bytes"
	"encoding/binary"
	"log"

	"golang.org/x/exp/constraints"
)

func IntToHex[T constraints.Integer](n T) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, n)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
