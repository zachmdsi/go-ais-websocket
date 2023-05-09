package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/zachmdsi/go-ais-vessel-tracking/pkg/ais"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()	

	// Read the initial data from the server
	data, err := ais.ReadData(bufio.NewReader(conn))
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		return
	}
	fmt.Printf("%d vessels are being tracked.\n", len(data))


	// Track the vessels and print the latest information
	_, err = ais.ManageTrackedVessels(data)
	if err != nil {
		fmt.Printf("Error managing tracked vessels: %v", err)
		return
	}
}
