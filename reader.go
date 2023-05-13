package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type VesselData struct {
	MMSI             string
	BaseDateTime     string
	LAT              float64
	LON              float64
	SOG              float64
	COG              float64
	Heading          float64
	VesselName       string
	IMO              string
	CallSign         string
	VesselType       int
	Status           string
	Length           float64
	Width            float64
	Draft            float64
	Cargo            int
	TransceiverClass string
}

func ReadCSV(filename string) ([]VesselData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	data := make([]VesselData, len(lines))
	for i, line := range lines[1:] { // Skipping header line
		data[i] = VesselData{
			MMSI:             line[0],
			BaseDateTime:     line[1],
			LAT:              parseToFloat64(line[2]),
			LON:              parseToFloat64(line[3]),
			SOG:              parseToFloat64(line[4]),
			COG:              parseToFloat64(line[5]),
			Heading:          parseToFloat64(line[6]),
			VesselName:       line[7],
			IMO:              line[8],
			CallSign:         line[9],
			VesselType:       parseToInt(line[10]),
			Status:           line[11],
			Length:           parseToFloat64(line[12]),
			Width:            parseToFloat64(line[13]),
			Draft:            parseToFloat64(line[14]),
			Cargo:            parseToInt(line[15]),
			TransceiverClass: line[16],
		}
	}

	return data, nil
}

func parseToFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("error: %v | input: strconv.ParseFloat(%s, 64)", err, s)
	}
	return v
}

func parseToInt(s string) int {
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("error: %v | input: strconv.Atoi(%s)", err, s)
	}
	return v
}
