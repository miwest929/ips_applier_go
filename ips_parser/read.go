package ips_parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type NormalChunk struct {
	Offset     uint32
	DataLength uint32
	Data       []byte
}

type RunLengthEncodingChunk struct{ Offset, ValueRepeatCount, Value uint32 }

func getNextBytes(f *os.File, byteCount int) []byte {
	bytes := make([]byte, byteCount)
	_, err := f.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func isRunLengthEncoded(nextBytes []byte) bool {
	rleIndicator := []byte{0x00, 0x00}
	isRle := bytes.Compare(nextBytes[3:5], rleIndicator)
	if isRle == 0 {
		return true
	}

	return false
}

func isEofMarker(nextBytes []byte) bool {
	eofMarker := []byte{0x45, 0x4f, 0x46}
	isEof := bytes.Compare(nextBytes[0:3], eofMarker)
	if isEof == 0 {
		fmt.Printf("EOF MARKER bytes=[% x]\n", nextBytes[0:3])
		return true
	}

	return false
}

func convertByteArrayToUint32(bytes []byte) uint32 {
	if len(bytes) < 4 {
		// Data is in BigEndian. Need to align the given bytes
		// so they are 4 bytes long.
		additional := 4 - len(bytes)
		for i := 0; i < additional; i++ {
			// This is method of prepending that I came up with
			bytes = append([]byte{0x00}, bytes...)
		}
	}
	return binary.BigEndian.Uint32(bytes)
}

func createNormalChunk(nextBytes []byte) *NormalChunk {
	normalChunk := NormalChunk{Offset: convertByteArrayToUint32(nextBytes[0:3]), DataLength: convertByteArrayToUint32(nextBytes[3:5])}
	return &normalChunk
}

func createRunLengthEncodingChunk(nextBytes []byte) *RunLengthEncodingChunk {
	rleChunk := RunLengthEncodingChunk{Offset: convertByteArrayToUint32(nextBytes[0:3])}
	return &rleChunk
}

func ReadIpsFile(ipsFile string) []interface{} {
	var chunks []interface{}
	startingBytes := []byte{0x50, 0x41, 0x54, 0x43, 0x48}

	f, err := os.Open(ipsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	header := make([]byte, 5)
	n1, err := f.Read(header)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d bytes: %s\n", n1, string(header[:n1]))

	res := bytes.Compare(header, startingBytes)
	if res != 0 {
		log.Fatal("The header bytes are incorrect")
	}

	bytes := getNextBytes(f, 5)
	for !isEofMarker(bytes) {
		if isRunLengthEncoded(bytes) {
			rleChunk := createRunLengthEncodingChunk(bytes)

			bytes := getNextBytes(f, 2)
			rleChunk.ValueRepeatCount = convertByteArrayToUint32(bytes)

			bytes = getNextBytes(f, 1)
			rleChunk.Value = convertByteArrayToUint32(bytes)

			chunks = append(chunks, rleChunk)
			// fmt.Printf("RunLengthEncodedChunk { offset=%d valueRepeatCount=%d value=%d}\n", rleChunk.Offset, rleChunk.ValueRepeatCount, rleChunk.Value)
		} else {
			normalChunk := createNormalChunk(bytes)
			bytes = getNextBytes(f, int(normalChunk.DataLength))
			normalChunk.Data = bytes

			chunks = append(chunks, normalChunk)
			// fmt.Printf("NormalChunk { offset=%d dataLength=%d data=[% x]}\n", normalChunk.Offset, normalChunk.DataLength, normalChunk.Data)
		}

		bytes = getNextBytes(f, 5)
	}

	return chunks
}
