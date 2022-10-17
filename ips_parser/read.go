package ips_parser

import (
	"os"
	"log"
	"fmt"
	"bytes"
)

type NormalChunk struct {
	offset        []byte
	dataLength    []byte
	data          []byte
}

type RunLengthEncodingChunk struct {
	offset           []byte
	valueRepeatCount []byte
	value            byte
}

func getNextChunk(f *os.File) []byte {
	chunk := make([]byte, 5)
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

func createNormalChunk(nextBytes []byte) *NormalChunk {
	normalChunk := NormalChunk{offset: nextBytes[0:3], dataLength: nextBytes[3:5]}
	// normalChunk.data = <next-data-Length>
	return &normalChunk
}

// func createRunLengthEncodingChunk(nextBytes []byte) RunLengthEncodingChunk {
	
// }

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

	chunk := getNextChunk(f)
	for !isEofMarker(chunk) {
		chunk = getNextChunk(f)

		if (isRunLengthEncoded(chunk)) {

		} else {

		}
	}

	f.Close();
}