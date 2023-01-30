package main

import (
	"fmt"
	"os"
)

func main() {
	// Get the IP address and port range from the command line
	// if len(os.Args) <= 4 {
	//	fmt.Printf("Usage: go run port_scanner.go <ports> <ip1> <ip2> ... <ipn>")
	//	os.Exit(1)
	//}
	//print out the arguments
	fmt.Println(os.Args[2])
}