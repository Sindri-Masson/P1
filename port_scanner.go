package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parse_ports(ports string) []string {
	var ports_list []string = strings.Split(ports, ",");
	for i := 0; i < len(ports_list); i++ {
		if strings.Contains(ports_list[i], "-") {
			var start_end []string = strings.Split(ports_list[i], "-");
			var strStart string = start_end[0];
			var start, start_err = strconv.Atoi(strStart);
			var end, end_err = strconv.Atoi(start_end[1]);
			ports_list = append(ports_list[:i], ports_list[i+1:]...);
			fmt.Println(start_err, end_err)
			for j := start; j <= end; j++ {
				ports_list = append(ports_list, strconv.Itoa(j));
			}
		}
	}
	return ports_list
}

func main() {
	// Get the IP address and port range from the command line
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run port_scanner.go <ports> <ip1> <ip2> ... <ipn>")
		os.Exit(1)
	}
	//print out the arguments
	var ports string = os.Args[1];
	var ports_list []string = parse_ports(ports);
	var ips []string = os.Args[2:len(os.Args)];

	fmt.Println(ports)
	fmt.Println(ports_list)
	fmt.Println(ips)
}