package utm

import (
	"fmt"
	"strconv"
	"unicode"
)

// Zone specifies the zone number and hemisphere
type Zone struct {
	N     int // 1-60
	North bool
}

// String returns a text representation of the zone
func (z Zone) String() string {
	if z.North {
		return fmt.Sprintf("%d (north)", z.N)
	}
	return fmt.Sprintf("%d (south)", z.N)
}

// SRID returns the zone EPSG/SRID code
func (z Zone) SRID() int {
	if z.North {
		return z.N + 32600
	}
	return z.N + 32700
}

// LookupSRID returns a Zone by its EPSG/SRID code
func LookupSRID(srid int) (Zone, bool) {
	if 32601 <= srid && srid <= 32660 {
		return Zone{
			N:     srid - 32600,
			North: true,
		}, true
	}
	if 32701 <= srid && srid <= 32760 {
		return Zone{
			N: srid - 32700,
		}, true
	}
	return Zone{}, false
}

// CentralMeridian returns the zone's center longitude
func (z Zone) CentralMeridian() float64 {
	return float64((z.N-1)*6 - 180 + 3)
}

// LatLonZone returns the Zone for the provided coordinates
func LatLonZone(latitude float64, longitude float64) Zone {
	north := latitude >= 0
	if 56 <= latitude && latitude <= 64 && 3 <= longitude && longitude <= 12 {
		return Zone{N: 32, North: north}
	}
	if 72 <= latitude && latitude <= 84 && longitude >= 0 {
		if longitude <= 9 {
			return Zone{N: 31, North: north}
		} else if longitude <= 21 {
			return Zone{N: 33, North: north}
		} else if longitude <= 33 {
			return Zone{N: 35, North: north}
		} else if longitude <= 42 {
			return Zone{N: 37, North: north}
		}
	}
	return Zone{
		N:     int((longitude+180)/6) + 1,
		North: north,
	}
}

// ParseZone parses a zone number followed by a zone letter
func ParseZone(s string) (Zone, bool) {
	if len(s) < 2 {
		return Zone{}, false
	}
	last := len(s) - 1
	n, err := strconv.Atoi(s[:last])
	if err != nil {
		return Zone{}, false
	}
	var north bool
	switch unicode.ToUpper(rune(s[last])) {
	case 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X':
		north = true
	case 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M':
	default:
		return Zone{}, false
	}
	return Zone{N: n, North: north}, true
}
