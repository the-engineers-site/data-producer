package store

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var connection, _ = net.Dial("tcp", os.Getenv("DESTINATION_IP"))

// Double the number of CPUs
var doubledCPUs, err = strconv.Atoi(os.Getenv("EPS"))

func init() {
	if err != nil {
		doubledCPUs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(doubledCPUs)
}

func Send(message string) {
	// Create a wait group to wait for all Go routines to finish
	var wg sync.WaitGroup
	// Launch Go routines
	for i := 0; i < doubledCPUs; i++ {
		wg.Add(1)
		go func(message string) {
			defer wg.Done()
			sendLineAsync(message)
		}(message)
	}

	// Wait for all Go routines to finish
	wg.Wait()
}

func sendLineAsync(message string) {
	_, err := fmt.Fprintln(connection, message)
	if err != nil {
		log.Println("Error while publishing ", err)
	}
}
