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

func scanner(host string, port int, semaphore chan struct{}, wg *sync.WaitGroup){
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
	if err != nil {
		fmt.Println(host, ":", port, "closed")
		return
	}
	conn.Close()
	fmt.Println(host, ":", port, "open")
}


func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run port_scanner.g <flag> <host> <ports>")
		os.Exit(1)
	}
	ports := parsePorts(os.Args[2])
	hosts := strings.Split(os.Args[3], " ")
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
