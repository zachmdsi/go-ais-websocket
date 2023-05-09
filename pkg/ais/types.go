package ais

import "time"

type VesselData struct {
	MMSI             string
	BaseDateTime     string
	LAT              float64
	LON              float64
	SOG              float64
	COG              float64
	Heading          string
	VesselName       string
	IMO              string
	CallSign         string
	VesselType       string
	Status           string
	Length           float64
	Width            float64
	Draft            float64
	Cargo            string
	TransceiverClass string
}

type TrackedVessel struct {
	MMSI      string
	LAT       float64
	LON       float64
	SOG       float64
	COG       float64
	Timestamp time.Time
}
