package ais

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ReadData(reader *bufio.Reader) ([]VesselData, error) {
	_, err := reader.ReadString('\n') // Skip the header
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	var data []VesselData
	lineNumber := 2 // Start at 2 because we skipped the header
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read line %d: %v", lineNumber, err)
		}
		fields := strings.Split(strings.TrimSpace(line), ",")
		if len(fields) != 17 {
			return nil, fmt.Errorf("invalid number of fields in line %d: %d", lineNumber, len(fields))
		}

		lat, err := parseFloat(fields[2], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing latitude on line %d: %v", lineNumber, err)
		}

		lon, err := parseFloat(fields[3], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing longitude on line %d: %v", lineNumber, err)
		}

		sog, err := parseFloat(fields[4], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing speed over ground on line %d: %v", lineNumber, err)
		}

		cog, err := parseFloat(fields[5], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing course over ground on line %d: %v", lineNumber, err)
		}

		length, err := parseFloat(fields[12], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing length on line %d: %v", lineNumber, err)
		}

		width, err := parseFloat(fields[13], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing width on line %d: %v", lineNumber, err)
		}

		draft, err := parseFloat(fields[14], 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing draft on line %d: %v", lineNumber, err)
		}

		data = append(data, VesselData{
			MMSI:             fields[0],
			BaseDateTime:     fields[1],
			LAT:              lat,
			LON:              lon,
			SOG:              sog,
			COG:              cog,
			Heading:          fields[6],
			VesselName:       fields[7],
			IMO:              fields[8],
			CallSign:         fields[9],
			VesselType:       fields[10],
			Status:           fields[11],
			Length:           length,
			Width:            width,
			Draft:            draft,
			Cargo:            fields[15],
			TransceiverClass: fields[16],
		})

		lineNumber++
	}

	return data, nil
}
