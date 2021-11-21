package main

import (
	"encoding/binary"
)

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func serializePiperResp(content []byte, respType uint8) []byte {
	var response []byte
	response = append(response, respType)
	lenb := make([]byte, 8)
	binary.LittleEndian.PutUint64(lenb, uint64(len(content)))
	response = append(response, lenb...)
	response = append(response, content...)
	return response
}
