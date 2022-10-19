package main

import (
	"applyIpsPatch/ips_parser"
	"fmt"
	"log"
	"os"
)

func applyIpsPatch(baseFilename string, chunks []interface{}) {
	baseFile, err := os.Open(baseFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer baseFile.Close()

	for _, ch := range chunks {
		switch c := ch.(type) {
		case *ips_parser.NormalChunk:
			fmt.Printf("NormalChunk: offset=%d\n", c.Offset)
			_, err := baseFile.Seek(int64(c.Offset), 0)
			if err != nil {
				log.Fatal(err)
			}
			// baseFile.Write(c.Data)
			// ytesWritten, err := file.Write(byteSlice)

		case *ips_parser.RunLengthEncodingChunk:
			fmt.Printf("RunLengthEncodingChunk: offset=%d\n", c.Offset)
		default:
			fmt.Printf("Encountered unknown chunk type: %T\n", c)
		}
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	ipsFile := argsWithoutProg[0]
	baseRomFile := argsWithoutProg[1]
	fmt.Printf("%s rom file\n", baseRomFile)
	chunks := ips_parser.ReadIpsFile(ipsFile)
	fmt.Printf("Number of chunks %d\n\n", len(chunks))
	applyIpsPatch(baseRomFile, chunks)
}
