package ais

import (
	"fmt"
	"strconv"
	"time"
)

func ManageTrackedVessels(data []VesselData) (map[string]TrackedVessel, error) {
	trackedVessels := make(map[string]TrackedVessel)

	for _, record := range data {
		// Check if the vessel is already being tracked
		mmsi := record.MMSI
		trackedVessel, ok := trackedVessels[mmsi]
		if !ok {
			// This vessel is not being tracked yet, so add it to the map
			trackedVessel = TrackedVessel{
				MMSI:      mmsi,
				Timestamp: time.Now(),
			}
		}

		// Update the vessel's data with the latest information
		trackedVessel.LAT = record.LAT
		trackedVessel.LON = record.LON
		trackedVessel.SOG = record.SOG
		trackedVessel.COG = record.COG
		trackedVessel.Timestamp, _ = time.Parse(record.BaseDateTime, "2023-05-09T14:30:00Z")
		trackedVessels[mmsi] = trackedVessel
	}

	// Remove vessels that have not been updated in the last 5 minutes
	for mmsi, trackedVessel := range trackedVessels {
		if time.Since(trackedVessel.Timestamp) > 5*time.Minute {
			delete(trackedVessels, mmsi)
		}
	}

	return trackedVessels, nil
}

func parseFloat(s string, defaultVal float64) (float64, error) {
	if s == "" || s == "N/A" || s == "UNKNOWN" {
		return defaultVal, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal, fmt.Errorf("error parsing float: %v", err)
	}
	return f, nil
}
