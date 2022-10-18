package main

import (
	"os"
	"fmt"
	"applyIpsPatch/ips_parser"
)

func main() {
	argsWithoutProg := os.Args[1:]
	ipsFile := argsWithoutProg[0]
	baseRomFile := argsWithoutProg[1]
	fmt.Printf("%s rom file\n", baseRomFile)
	chunks := ips_parser.ReadIpsFile(ipsFile)
	fmt.Printf("Number of chunks %d", len(chunks))
}