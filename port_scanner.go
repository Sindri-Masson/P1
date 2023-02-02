package main

import (
	"fmt"
	"os"
	"strconv"
	"net"
	"sync"
	"strings"
	"time"
)

func parse_ports(ports string) []int {
	ports_list := strings.Split(ports, ",");
	var ports_list_int []int;
	for _, portString := range ports_list {
		if strings.Contains(portString, "-") {
			var start_end []string = strings.Split(portString, "-");
			var strStart string = start_end[0];
			var start, start_err = strconv.Atoi(strStart);
			var end, end_err = strconv.Atoi(start_end[1]);
			//ports_list = append(ports_list[:i], ports_list[i+1:]...);
			if start_err != nil || end_err != nil {
				fmt.Println("Error parsing ports")
				os.Exit(1);
			}
			for j := start; j <= end; j++ {
				ports_list_int = append(ports_list_int, j);
			}
		} else {
			var port, err = strconv.Atoi(portString);
			fmt.Println(port);
			if err != nil {
				fmt.Println("Error parsing ports")
				os.Exit(1);
			}
			ports_list_int = append(ports_list_int, port);
		}
	}
	return ports_list_int
}

func scan(host string, ports []int, semaphore chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done();
	for _, port := range ports {
		semaphore <- struct{}{}
		go func(host string, port int) {
			defer func() { 
				<-semaphore 
			}()
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
		if err != nil {
			fmt.Println(host, ":", port, "closed")
			return
		}
		conn.Close()
		fmt.Println(host, ":", port, "open")
		}(host,port)
		wg.Add(1)
	}
}



func main() {
	// Get the IP address and port range from the command line
	if len(os.Args) < 3 {
		fmt.Printf("Usage: go run port_scanner.go <ports> <ip1> <ip2> ... <ipn>")
		os.Exit(1)
	}
	//print out the arguments
	
	var portsS string = os.Args[2];
	var hostsS string = os.Args[3];
	ports := parse_ports(portsS);
	hosts := strings.Split(hostsS, " ");
	//var ips []string = os.Args[2:len(os.Args)];
	semaphore := make(chan struct{}, 20);
	var wg sync.WaitGroup
	wg.Add(len(hosts))
	for _, host := range hosts {
		go func(host string) {
			defer wg.Done()
			var wg sync.WaitGroup
			wg.Add(len(ports))
			scan(host, ports, semaphore, &wg)
			wg.Wait()
		}(host)
	}
	wg.Wait()
	

	fmt.Println(portsS)
	fmt.Println(ports)
	fmt.Println(hosts)
	fmt.Println(hostsS)
}