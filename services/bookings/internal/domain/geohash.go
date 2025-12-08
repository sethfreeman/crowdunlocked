package domain

import (
	"math"
)

// Geohash utilities for spatial indexing
// Based on standard geohash algorithm

const (
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz"
)

// EncodeGeohash encodes latitude and longitude into a geohash string
func EncodeGeohash(lat, lng float64, precision int) string {
	if precision <= 0 || precision > 12 {
		precision = 6 // Default to ~1.2km precision
	}

	latRange := [2]float64{-90.0, 90.0}
	lngRange := [2]float64{-180.0, 180.0}

	geohash := make([]byte, 0, precision)
	bits := 0
	bit := 0
	ch := 0

	for len(geohash) < precision {
		if bits%2 == 0 {
			// Longitude
			mid := (lngRange[0] + lngRange[1]) / 2
			if lng > mid {
				ch |= (1 << (4 - bit))
				lngRange[0] = mid
			} else {
				lngRange[1] = mid
			}
		} else {
			// Latitude
			mid := (latRange[0] + latRange[1]) / 2
			if lat > mid {
				ch |= (1 << (4 - bit))
				latRange[0] = mid
			} else {
				latRange[1] = mid
			}
		}

		bits++
		bit++

		if bit == 5 {
			geohash = append(geohash, base32[ch])
			bit = 0
			ch = 0
		}
	}

	return string(geohash)
}

// DecodeGeohash decodes a geohash string into latitude and longitude
func DecodeGeohash(geohash string) (lat, lng float64) {
	latRange := [2]float64{-90.0, 90.0}
	lngRange := [2]float64{-180.0, 180.0}

	isEven := true
	for _, c := range geohash {
		idx := indexInBase32(byte(c))
		if idx == -1 {
			continue
		}

		for i := 4; i >= 0; i-- {
			bit := (idx >> i) & 1
			if isEven {
				// Longitude
				mid := (lngRange[0] + lngRange[1]) / 2
				if bit == 1 {
					lngRange[0] = mid
				} else {
					lngRange[1] = mid
				}
			} else {
				// Latitude
				mid := (latRange[0] + latRange[1]) / 2
				if bit == 1 {
					latRange[0] = mid
				} else {
					latRange[1] = mid
				}
			}
			isEven = !isEven
		}
	}

	lat = (latRange[0] + latRange[1]) / 2
	lng = (lngRange[0] + lngRange[1]) / 2
	return
}

// GetGeohashPrefixes returns geohash prefixes for radius search
// Returns multiple prefixes that cover the search area
func GetGeohashPrefixes(lat, lng, radiusKm float64) []string {
	// Calculate precision based on radius
	precision := getGeohashPrecision(radiusKm)
	
	center := EncodeGeohash(lat, lng, precision)
	
	// Get neighboring geohashes to cover the radius
	neighbors := getNeighbors(center)
	
	prefixes := make([]string, 0, len(neighbors)+1)
	prefixes = append(prefixes, center)
	prefixes = append(prefixes, neighbors...)
	
	return prefixes
}

// getGeohashPrecision returns appropriate precision for radius
func getGeohashPrecision(radiusKm float64) int {
	// Geohash precision to approximate area coverage
	// 1: ±2500 km
	// 2: ±630 km
	// 3: ±78 km
	// 4: ±20 km
	// 5: ±2.4 km
	// 6: ±0.61 km
	// 7: ±0.076 km
	
	switch {
	case radiusKm > 100:
		return 3
	case radiusKm > 20:
		return 4
	case radiusKm > 5:
		return 5
	case radiusKm > 1:
		return 6
	default:
		return 7
	}
}

// getNeighbors returns the 8 neighboring geohashes
func getNeighbors(geohash string) []string {
	if len(geohash) == 0 {
		return []string{}
	}

	lat, lng := DecodeGeohash(geohash)
	precision := len(geohash)
	
	// Approximate cell size for this precision
	cellSize := getCellSize(precision)
	
	neighbors := make([]string, 0, 8)
	
	// Generate 8 neighbors (N, NE, E, SE, S, SW, W, NW)
	offsets := [][2]float64{
		{cellSize, 0},      // E
		{cellSize, cellSize},   // NE
		{0, cellSize},      // N
		{-cellSize, cellSize},  // NW
		{-cellSize, 0},     // W
		{-cellSize, -cellSize}, // SW
		{0, -cellSize},     // S
		{cellSize, -cellSize},  // SE
	}
	
	for _, offset := range offsets {
		neighborLat := lat + offset[1]
		neighborLng := lng + offset[0]
		
		// Wrap around if needed
		if neighborLat > 90 {
			neighborLat = 90
		}
		if neighborLat < -90 {
			neighborLat = -90
		}
		if neighborLng > 180 {
			neighborLng -= 360
		}
		if neighborLng < -180 {
			neighborLng += 360
		}
		
		neighbor := EncodeGeohash(neighborLat, neighborLng, precision)
		if neighbor != geohash {
			neighbors = append(neighbors, neighbor)
		}
	}
	
	return neighbors
}

// getCellSize returns approximate cell size in degrees for precision
func getCellSize(precision int) float64 {
	// Approximate cell sizes in degrees
	sizes := map[int]float64{
		1: 45.0,
		2: 5.625,
		3: 1.40625,
		4: 0.17578125,
		5: 0.0439453125,
		6: 0.0054931640625,
		7: 0.001373291015625,
	}
	
	if size, ok := sizes[precision]; ok {
		return size
	}
	return 0.001 // Default
}

// CalculateDistance calculates distance between two points using Haversine formula
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := toRadians(lat2 - lat1)
	dLng := toRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func indexInBase32(c byte) int {
	for i := 0; i < len(base32); i++ {
		if base32[i] == c {
			return i
		}
	}
	return -1
}
