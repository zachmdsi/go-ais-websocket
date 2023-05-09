package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP Server is listening on port 8080...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		fmt.Printf("Received signal %v, shutting down...", sig)
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, "data/ais-sample-data.csv", 1000)
	}
}

func handleConnection(conn net.Conn, filename string, batchSize int) {
    defer conn.Close()

    // Open the CSV file for reading
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()

    reader := csv.NewReader(file)
    _, err = reader.Read() // Skip the header
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        return
    }

    lineNumber := 2 // Start at 2 because we skipped the header
    records := make([]string, 0, batchSize)
    for {
        // Read a record from the CSV file
        record, err := reader.Read()
        if err == io.EOF {
            fmt.Println("End of file, closing connection")
            return
        }
        if err != nil {
            fmt.Printf("Error reading line %d: %v\n", lineNumber, err)
            return
        }

        // Append the record to the batch
        data := strings.Join(record, ",") + "\n"
        records = append(records, data)

        // If the batch is full, write it to the client and reset the batch
        if len(records) == batchSize {
            _, err = conn.Write([]byte(strings.Join(records, "")))
            if err != nil {
                fmt.Printf("Error writing data: %v\n", err)
                return
            }
            records = make([]string, 0, batchSize)
        }

        lineNumber++
    }
}
