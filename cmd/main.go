package main

import (
	"fmt"

	"github.com/zachmdsi/go-ais-vessel-tracking/pkg/ais"
)

func main() {
	data, err := ais.ReadData("data/ais-sample-data.csv")
	if err != nil {
		panic(err)
	}

	// Track the vessels and print the latest information
	_, err = ais.ManageTrackedVessels(data)
	if err != nil {
		fmt.Printf("Error managing tracked vessels: %v", err)
		return
	}
}
