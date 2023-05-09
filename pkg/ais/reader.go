package ais

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ReadData(filename string) ([]VesselData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open data: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // skip the header
	if err != nil {
		return nil, fmt.Errorf("failed to read: %v", err)
	}

	var data []VesselData
	lineNumber := 2 // Start at 2 because we skipped the header
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read line %d: %v", lineNumber, err)
		}	

		lat, err := parseFloat(record[2], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing latitude: %v", err)
		}

		lon, err := parseFloat(record[3], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing longitude: %v", err)
		}

		sog, err := parseFloat(record[4], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing speed over ground: %v", err)
		}

		cog, err := parseFloat(record[5], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing course over ground: %v", err)
		}

		length, err := parseFloat(record[12], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing length: %v", err)
		}

		width, err := parseFloat(record[13], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing width: %v", err)
		}

		draft, err := parseFloat(record[14], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing draft: %v", err)
		}

		data = append(data, VesselData{
			MMSI:             record[0],
			BaseDateTime:     record[1],
			LAT:              lat,
			LON:              lon,
			SOG:              sog,
			COG:              cog,
			Heading:          record[6],
			VesselName:       record[7],
			IMO:              record[8],
			CallSign:         record[9],
			VesselType:       record[10],
			Status:           record[11],
			Length:           length,
			Width:            width,
			Draft:            draft,
			Cargo:            record[15],
			TransceiverClass: record[16],
		})
	}

	return data, nil
}
