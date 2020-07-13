package utm

import (
	"fmt"
	"strconv"
	"unicode"
)

// Zone specifies the zone number and hemisphere
type Zone struct {
	Number int  // Zone number 1 to 60
	Letter rune // Zone letter C to X (omitting O, I)
	North  bool // Zone hemisphere
}

// String returns a text representation of the zone
func (z Zone) String() string {
	if z.Letter == 0 {
		z.Letter = '?'
	}
	if z.North {
		return fmt.Sprintf("%d%c (north)", z.Number, z.Letter)
	}
	return fmt.Sprintf("%d%c (south)", z.Number, z.Letter)
}

// Valid checks if the zone is valid
func (z Zone) Valid() bool {
	switch z.Letter {
	case 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X':
		if !z.North {
			return false
		}
	case 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M':
		if z.North {
			return false
		}
	case 0:
	default:
		return false
	}
	return 1 <= z.Number && z.Number <= 60
}

// SRID returns the zone EPSG/SRID code
func (z Zone) SRID() int {
	if z.North {
		return z.Number + 32600
	}
	return z.Number + 32700
}

// LookupSRID returns a Zone by its EPSG/SRID code.
// Since the srid code only specifies the longitude zone
// number, the zone letter is left unset.
func LookupSRID(srid int) (Zone, bool) {
	if 32601 <= srid && srid <= 32660 {
		return Zone{
			Number: srid - 32600,
			North:  true,
		}, true
	}
	if 32701 <= srid && srid <= 32760 {
		return Zone{
			Number: srid - 32700,
		}, true
	}
	return Zone{}, false
}

// CentralMeridian returns the zone's center longitude
func (z Zone) CentralMeridian() float64 {
	return float64((z.Number-1)*6 - 180 + 3)
}

// LatLonZone returns the Zone for the provided coordinates
func LatLonZone(latitude, longitude float64) Zone {
	const letters = "CDEFGHJKLMNPQRSTUVWXX"
	var letter rune
	if -80 <= latitude && latitude <= 84 {
		letter = rune(letters[int(latitude+80)>>3])
	}
	north := latitude >= 0
	if 56 <= latitude && latitude <= 64 && 3 <= longitude && longitude <= 12 {
		return Zone{Number: 32, Letter: letter, North: north}
	}
	if 72 <= latitude && latitude <= 84 && longitude >= 0 {
		if longitude <= 9 {
			return Zone{Number: 31, Letter: letter, North: north}
		} else if longitude <= 21 {
			return Zone{Number: 33, Letter: letter, North: north}
		} else if longitude <= 33 {
			return Zone{Number: 35, Letter: letter, North: north}
		} else if longitude <= 42 {
			return Zone{Number: 37, Letter: letter, North: north}
		}
	}
	return Zone{
		Number: int((longitude+180)/6) + 1,
		Letter: letter,
		North:  north,
	}
}

// ParseZone parses a zone number followed by a zone letter
func ParseZone(s string) (Zone, bool) {
	if len(s) < 2 {
		return Zone{}, false
	}
	last := len(s) - 1
	n, err := strconv.Atoi(s[:last])
	if err != nil || n < 1 || n > 60 {
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
	return Zone{Number: n, Letter: rune(s[last]), North: north}, true
}
