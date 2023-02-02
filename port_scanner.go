package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func parsePorts(ports string) []int {
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

//parse the ip addresses and check for cidr notation
func parse_ips(ips []string) []string {
	var ips_list []string;
	for _, ip := range ips {
		if strings.Contains(ip, "/") {
			var ip_cidr []string = strings.Split(ip, "/");
			var ip_start = ip_cidr[0];
			ip_chunks := strings.Split(ip_start, ".");
			var ip_mask, err = strconv.Atoi(ip_cidr[1]);
			if err != nil {
				fmt.Println("Error parsing ip mask")
				os.Exit(1);
			}
			num_hosts := math.Pow(2, float64(32 - ip_mask))

			for i := 1; i < int(num_hosts) - 1; i++ {
				var ip string = ip_chunks[0] + "." + ip_chunks[1] + "." + ip_chunks[2] + "." + strconv.Itoa(i);
				ips_list = append(ips_list, ip);
			}
		} else {
			ips_list = append(ips_list, ip);
		}
	}
	fmt.Println(ips_list)
	return ips_list;
}

func scanner(host string, port int, semaphore chan struct{}, wg *sync.WaitGroup){
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
	if err != nil {
		fmt.Println(host, ":", port, "closed")
		return
	}
	conn.Close()
	fmt.Println(host, ":", port, " open")
}


func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run port_scanner.g <flag> <host> <ports>")
		os.Exit(1)
	}
	ports := parsePorts(os.Args[2])
	hosts := parse_ips(os.Args[3:])
	semaphore := make(chan struct{}, 100)
	var wg sync.WaitGroup
	for _, host := range hosts {
		for _, port := range ports {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(host string, port int) {
				defer func() { 
					<-semaphore
					wg.Done()
				}()
				scanner(host, port, semaphore, &wg)
			}(host,port)
		}
	}
	wg.Wait()
}
