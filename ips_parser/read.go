package ips_parser

import (
	"os"
	"log"
	"fmt"
	"bytes"
	"encoding/binary"
)

type NormalChunk struct {
	offset        uint64
	dataLength    uint64
	data          uint64
}

type RunLengthEncodingChunk struct {
	offset           uint64
	valueRepeatCount uint64
	value            uint64
}

func getNextBytes(f *os.File, byteCount uint64) []byte {
	chunk := make([]byte, byteCount)
	bytesRead, err := f.Read(chunk)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d bytes: [% x]\n", bytesRead, chunk[:bytesRead])
	return chunk
}

func isRunLengthEncoded(nextBytes []byte) bool {
	rleIndicator := []byte{0x00, 0x00}
	isRle := bytes.Compare(nextBytes[3:5], rleIndicator)
	if (isRle != 0) {
		return true;
	}

	return false;
}

func isEofMarker(nextBytes []byte) bool {
	eofMarker := []byte{0x45, 0x4f, 0x46};
	isEof := bytes.Compare(nextBytes[0:3], eofMarker)
	if (isEof != 0) {
		return true
	}

	return false
}

func convertByteArrayToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

func createNormalChunk(nextBytes []byte) *NormalChunk {
	normalChunk := NormalChunk{offset: convertByteArrayToUint64(nextBytes[0:3]), dataLength: convertByteArrayToUint64(nextBytes[3:5])}
	// normalChunk.data = <next-data-Length>
	return &normalChunk
}

func createRunLengthEncodingChunk(nextBytes []byte) *RunLengthEncodingChunk {
	rleChunk := RunLengthEncodingChunk{offset: convertByteArrayToUint64(nextBytes[0:3])}
	return &rleChunk
}

func ReadIpsFile(ipsFile string) {
	startingBytes := []byte{0x50, 0x41, 0x54, 0x43, 0x48}

	f, err := os.Open(ipsFile)
	if err != nil {
		log.Fatal(err)
	}

	header := make([]byte, 5)
	n1, err := f.Read(header)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d bytes: %s\n", n1, string(header[:n1]))

	res := bytes.Compare(header, startingBytes)
	if (res != 0) {
		log.Fatal("The header bytes are incorrect")
	}

	bytes := getNextBytes(f, 5)
	for !isEofMarker(bytes) {
		if (isRunLengthEncoded(bytes)) {
			createRunLengthEncodingChunk := createRunLengthEncodingChunk(bytes)

			bytes := getNextBytes(f, 2)
			// count := binary.BigEndian.Uint64(repeatedValueBytes)
			createRunLengthEncodingChunk.valueRepeatCount = convertByteArrayToUint64(bytes)

			bytes = getNextBytes(f, 1)
			createRunLengthEncodingChunk.value = convertByteArrayToUint64(bytes)
		} else {
			normalChunk := createNormalChunk(bytes)
			// convertByteArrayToUint64
			byte = getNextBytes(f, normalChunk.dataLength)

		}

		
	}

	f.Close();
}