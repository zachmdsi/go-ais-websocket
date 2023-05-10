package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)


func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP Server is listening on port 8080...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var connections []net.Conn
	var mutex sync.Mutex
	var shuttingDown bool

	go func() {
		sig := <-signals
		fmt.Printf("Received signal %v, shutting down...\n", sig)
		shuttingDown = true
		ln.Close()
		mutex.Lock()
		for _, conn := range connections {
			conn.Close()
		}
		mutex.Unlock()
	}()

	for {
		if shuttingDown {
			fmt.Println("Server is shutting down...")
			break
		}
		conn, err := ln.Accept()
		if err != nil {
			if !shuttingDown {
				fmt.Printf("Error accepting connection: %v\n", err)
			}
			continue
		}
		mutex.Lock()
		connections = append(connections, conn)
		mutex.Unlock()
		go handleConnection(conn, "data/ais-sample-data.csv", 100, &connections, &mutex)
	}
}

func handleConnection(conn net.Conn, filename string, batchSize int, connections *[]net.Conn, mutex *sync.Mutex) {
	defer func() {
		conn.Close()
		mutex.Lock()
		for i, c := range *connections {
			if c == conn {
				*connections = append((*connections)[:i], (*connections)[i+1:]...)
				break
			}
		}
		mutex.Unlock()
	}()

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
        if err != nil {
            if err == io.EOF {
                fmt.Println("End of file, closing connection")
            } else {
                fmt.Printf("Error reading line %d: %v\n", lineNumber, err)
            }

            // Send any remaining records
            if len(records) > 0 {
                _, err = conn.Write([]byte(strings.Join(records, "")))
                if err != nil {
                    fmt.Printf("Error writing data: %v\n", err)
                }
            }
            return
        } 

        // Append the record to the batch
        data := strings.Join(record, ",") + "\n"
        records = append(records, data)

        // If the batch is full, write it to the client, reset the batch, and pause for 1 second
        if len(records) == batchSize {
            _, err = conn.Write([]byte(strings.Join(records, "")))
            if err != nil {
                fmt.Printf("Error writing data: %v\n", err)
                return
            }
            records = make([]string, 0, batchSize)
            time.Sleep(1 * time.Second)
        }

        lineNumber++
    }
}
