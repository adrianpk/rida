package client

import (
	"math"

	"github.com/adrianpk/rida/internal/telemetry"
)

// Number of meters per degree of latitude (approximately constant)
const mPerDegLat = 111320.0

// deltaLat returns the latitude offset in degrees for a given distance in
// meters.
func deltaLat(D float64) float64 {
	return D / mPerDegLat
}

// deltaLng returns the longitude offset in degrees for a given distance in
// meters at a specific latitude.
func deltaLng(D, lat float64) float64 {
	phi := lat * math.Pi / 180
	return D / (mPerDegLat * math.Cos(phi))
}

// BoundingBox returns a bounding box (telemetry.Area)
// around the given lat/lng for a given distance in meters.
func BoundingBox(lat, lng, distanceMeters float64) telemetry.Area {
	dLat := deltaLat(distanceMeters)
	dLng := deltaLng(distanceMeters, lat)
	return telemetry.Area{
		MinLat: lat - dLat,
		MaxLat: lat + dLat,
		MinLng: lng - dLng,
		MaxLng: lng + dLng,
	}
}
