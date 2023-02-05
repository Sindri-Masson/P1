# Programming assignment 1: Concurrent port scanner

### Authors: Magnús Atli Gylfason, Sindri Másson

## Introduction

This is a simple port scanner that can scan multiple ranges of ports on multiple hosts concurrently. It can decode CIDR notation.
To run the program, run the following command:

    $ go run port_scanner.go -p <comma separated ports and ranges> <space separated hosts>

The program consists of 4 functions, main, scanner, parse_ports and parse_ips. The main function is the entry point of the program, it parses the command line arguments and starts the goroutines that scan the ports. The function scanner scans a single port on a single host and then writes whether the port is open or not to the console. The function parse_ports parses the ports and ranges of ports given as command line arguments and returns a integer slice of ports to scan. The function parse_ips parses the hosts and CIDR subnets given as command line arguments and returns a slice of hosts to scan.

## Discussion

This program runs a maximum of 10000 goroutines at any time. It starts a new goroutine for each host and port combination. It only starts a new goroutine if there are less than 10000 goroutines running at the time and if there are still port/host combinations to scan(goroutine started in line 94). To make sure all goroutines have finished before the program exits, a waitgroup is used. Each goroutine increments the waitgroup wg by 1 when it is started and decrements it by 1 when it finishes (done with wg.Add(1) in line 93 and wg.Done() in line 97 respectively). The program waits for the waitgroup to reach 0 before exiting.
Each goroutine runs the function scanner which scans a single port on a single host, if the port is open, it prints "host:port open", but if its closed it prints "host:port closed". Since the program finishes when all goroutines have finished, we can assume that all goroutines do finish whether the port it is scanning is open, closed or unreachable.
To avoid deadlocks, the function scanner uses the function net.DialTimeout instead of net.Dial to scan the ports, this makes sure that no goroutine is waiting for a response from a port that is closed. This also makes sure that the program does not hang if a port is open but does not respond to the connection request. Since there are no other situations where a goroutine is blocked, the solution is deadlock free.
To keep the program running at 10000 goroutines at any time, the buffered channel semaphore is used. The channel is initialized with 10000 empty buffer. When a goroutine is started, it takes up a buffer in the channel and when it finishes, it releases the buffer. The channel makes a new goroutine that is trying to get a buffer wait until a buffer is available. This makes sure that there are never more than 10000 goroutines running at any time and also makes sure that each time a goroutine finishes, there is another one to take its place.
The number of goroutines was determined by testing different values. We started with a 100 available buffers in the channel to design the program, we knew however that that was way too few buffers. We then tested different values by starting at 100000 and decreasing the number by 10000 each time, until we recorded no resource exhaustion. That first happened at 10000, we then tested 12000 buffers and saw small amounts of resource exhaustion, so we decided to go with 10000.
