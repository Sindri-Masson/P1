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

func parse_ports(ports string) []int {
	//parses the ports string and returns a list of the ports as ints
	ports_list := strings.Split(ports, ","); //split the ports string into a list of ports whenever there is a comma
	var ports_list_int []int;
	for _, portString := range ports_list { //for loop to iterate through the ports list
		if strings.Contains(portString, "-") { //if the port is a range of ports
			var start_end []string = strings.Split(portString, "-");
			var start, start_err = strconv.Atoi(start_end[0]); //start of the range
			var end, end_err = strconv.Atoi(start_end[1]); //end of the range
			//ports_list = append(ports_list[:i], ports_list[i+1:]...);
			if start_err != nil || end_err != nil {
				fmt.Println("Error parsing ports")
				os.Exit(1);
			}
			for j := start; j <= end; j++ {
				ports_list_int = append(ports_list_int, j); //add the ports to the list in int type
			}
		} else { //if the port is a single port just add it to the list
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

func parse_ips(ips []string) []string {
	var ips_list []string;
	for _, ip := range ips { //for loop to iterate through the ips list
		if strings.Contains(ip, "/") { //if the ip is a range of ips (CIDR)
			var ip_cidr []string = strings.Split(ip, "/"); //split the ip into the ip and the mask
			var ip_start = ip_cidr[0];
			ip_chunks := strings.Split(ip_start, "."); //split the ip into the 4 chunks
			var ip_mask, err = strconv.Atoi(ip_cidr[1]);  //get the mask
			if err != nil {
				fmt.Println("Error parsing ip mask")
				os.Exit(1);
			}
			num_hosts := math.Pow(2, float64(32 - ip_mask)) 		//calculate the number of hosts

			for i := 1; i < int(num_hosts) - 1; i++ {
				var ip string = ip_chunks[0] + "." + ip_chunks[1] + "." + ip_chunks[2] + "." + strconv.Itoa(i); //add the ip to the list
				ips_list = append(ips_list, ip); 
			}
		} else {
			ips_list = append(ips_list, ip); //add the ip to the list
		}
	}
	fmt.Println(ips_list)
	return ips_list;
}

func scanner(host string, port int, semaphore chan struct{}, wg *sync.WaitGroup){
	//takes one host and and one port at a time and scans
	//then displays whether the port is open or closed
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 3*time.Second)
	if err != nil {
		fmt.Println(host + ":" + strconv.Itoa(port) + " closed")
		return
	}
	conn.Close()
	fmt.Println(host + ":" + strconv.Itoa(port) + " open")
}

func main() {
	start := time.Now()
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run port_scanner.g <flag> <host> <ports>")
		os.Exit(1)
	}
	ports := parse_ports(os.Args[2])
	hosts := parse_ips(os.Args[3:])
	semaphore := make(chan struct{}, 10000) //create a semaphore with a buffer of 10000
	var wg sync.WaitGroup
	for _, host := range hosts {	// for loop to iterate through hosts and ports
		for _, port := range ports {
			semaphore <- struct{}{} //acquire semaphore
			wg.Add(1) //increment waitgroup
			go func(host string, port int) {
				defer func() {
					<-semaphore //release semaphore
					wg.Done() //decrement waitgroup
				}()
				scanner(host, port, semaphore, &wg)
			}(host,port)
		}
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Time taken: ", elapsed)
}