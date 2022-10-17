package main

import (
	"os"
	//"log"
	"fmt"
	//"bytes"
	"applyIpsPatch/ips_parser"
)

// const EOF_MARKER = 0x454f46

// func readIpsFile(ipsFile string) {
// 	startingBytes := []byte{0x50, 0x41, 0x54, 0x43, 0x48}
// 	f, err := os.Open(ipsFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	header := make([]byte, 5)
// 	n1, err := f.Read(header)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("%d bytes: %s\n", n1, string(header[:n1]))

// 	res := bytes.Compare(header, startingBytes)
// 	if (res != 0) {
// 		log.Fatal("The header bytes are incorrect")
// 	}

// 	f.Close();
// }

func main() {
	argsWithoutProg := os.Args[1:]
	ipsFile := argsWithoutProg[0]
	baseRomFile := argsWithoutProg[1]
	fmt.Printf("%s rom file\n", baseRomFile)
	ips_parser.ReadIpsFile(ipsFile)
}