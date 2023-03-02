// port scanner
// scans the first 1024 ports of scanme.nmap.org

package main

import (
	"fmt"
	"net"
	"sort"
)

const ADDRESS = "scanme.nmap.org"

// worker
// takes a channel of ports and a channel of results
// if port is closed return a 0
// if port is open return the port number
func worker(ports, results chan int) {

	for p := range ports {
		address := fmt.Sprintf("%s:%d", ADDRESS, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {

	// create channels
	ports := make(chan int, 100)
	results := make(chan int)

	// slice to store open ports
	var openports []int

	// create 100 workers
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// send 1024 ports to the channel
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// receive results
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// close channels
	close(ports)
	close(results)
	sort.Ints(openports)

	// print open ports
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

}
