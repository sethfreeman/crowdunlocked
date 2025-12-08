package domain

import (
	"math"
	"testing"
)

func TestEncodeGeohash(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		lng       float64
		precision int
		want      string
	}{
		{
			name:      "San Francisco",
			lat:       37.7749,
			lng:       -122.4194,
			precision: 6,
			want:      "9q8yyk",
		},
		{
			name:      "New York",
			lat:       40.7128,
			lng:       -74.0060,
			precision: 6,
			want:      "dr5reg",
		},
		{
			name:      "London",
			lat:       51.5074,
			lng:       -0.1278,
			precision: 6,
			want:      "gcpvj0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeGeohash(tt.lat, tt.lng, tt.precision)
			if got != tt.want {
				t.Errorf("EncodeGeohash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeGeohash(t *testing.T) {
	tests := []struct {
		name    string
		geohash string
		wantLat float64
		wantLng float64
		delta   float64
	}{
		{
			name:    "San Francisco",
			geohash: "9q8yyk",
			wantLat: 37.7749,
			wantLng: -122.4194,
			delta:   0.01, // ~1km tolerance
		},
		{
			name:    "New York",
			geohash: "dr5reg",
			wantLat: 40.7128,
			wantLng: -74.0060,
			delta:   0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLat, gotLng := DecodeGeohash(tt.geohash)
			if math.Abs(gotLat-tt.wantLat) > tt.delta {
				t.Errorf("DecodeGeohash() lat = %v, want %v (±%v)", gotLat, tt.wantLat, tt.delta)
			}
			if math.Abs(gotLng-tt.wantLng) > tt.delta {
				t.Errorf("DecodeGeohash() lng = %v, want %v (±%v)", gotLng, tt.wantLng, tt.delta)
			}
		})
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		lng       float64
		precision int
		delta     float64
	}{
		{
			name:      "precision 6",
			lat:       37.7749,
			lng:       -122.4194,
			precision: 6,
			delta:     0.01,
		},
		{
			name:      "precision 8",
			lat:       40.7128,
			lng:       -74.0060,
			precision: 8,
			delta:     0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geohash := EncodeGeohash(tt.lat, tt.lng, tt.precision)
			gotLat, gotLng := DecodeGeohash(geohash)
			
			if math.Abs(gotLat-tt.lat) > tt.delta {
				t.Errorf("Roundtrip lat = %v, want %v (±%v)", gotLat, tt.lat, tt.delta)
			}
			if math.Abs(gotLng-tt.lng) > tt.delta {
				t.Errorf("Roundtrip lng = %v, want %v (±%v)", gotLng, tt.lng, tt.delta)
			}
		})
	}
}

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lng1     float64
		lat2     float64
		lng2     float64
		wantKm   float64
		deltaKm  float64
	}{
		{
			name:    "SF to LA",
			lat1:    37.7749,
			lng1:    -122.4194,
			lat2:    34.0522,
			lng2:    -118.2437,
			wantKm:  559.0,
			deltaKm: 10.0,
		},
		{
			name:    "NYC to Boston",
			lat1:    40.7128,
			lng1:    -74.0060,
			lat2:    42.3601,
			lng2:    -71.0589,
			wantKm:  306.0,
			deltaKm: 10.0,
		},
		{
			name:    "Same location",
			lat1:    37.7749,
			lng1:    -122.4194,
			lat2:    37.7749,
			lng2:    -122.4194,
			wantKm:  0.0,
			deltaKm: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDistance(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			if math.Abs(got-tt.wantKm) > tt.deltaKm {
				t.Errorf("CalculateDistance() = %v km, want %v km (±%v)", got, tt.wantKm, tt.deltaKm)
			}
		})
	}
}

func TestGetGeohashPrefixes(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lng      float64
		radiusKm float64
		minCount int
	}{
		{
			name:     "Small radius",
			lat:      37.7749,
			lng:      -122.4194,
			radiusKm: 1.0,
			minCount: 1, // At least center geohash
		},
		{
			name:     "Medium radius",
			lat:      37.7749,
			lng:      -122.4194,
			radiusKm: 10.0,
			minCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGeohashPrefixes(tt.lat, tt.lng, tt.radiusKm)
			if len(got) < tt.minCount {
				t.Errorf("GetGeohashPrefixes() returned %v prefixes, want at least %v", len(got), tt.minCount)
			}
			
			// Verify all prefixes are valid geohashes
			for _, prefix := range got {
				if len(prefix) == 0 {
					t.Errorf("GetGeohashPrefixes() returned empty prefix")
				}
			}
		})
	}
}
