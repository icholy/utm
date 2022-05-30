package utm

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

var approx = cmpopts.EquateApprox(0, 0.001)

func TestFromLatLon(t *testing.T) {
	tests := []struct {
		name                string
		latitude, longitude float64
		easting, northing   float64
		zone                Zone
	}{
		{
			name:      "Aachen, Germany",
			latitude:  50.77535,
			longitude: 6.08389,
			easting:   294408.917,
			northing:  5628897.997,
			zone:      Zone{Number: 32, Letter: 'U', North: true},
		},
		{
			name:      "New York, USA",
			latitude:  40.71435,
			longitude: -74.00597,
			easting:   583959.959,
			northing:  4507523.087,
			zone:      Zone{Number: 18, Letter: 'T', North: true},
		},
		{
			name:      "Wellington, New Zealand",
			latitude:  -41.28646,
			longitude: 174.77624,
			easting:   313784.305,
			northing:  5427057.321,
			zone:      Zone{Number: 60, Letter: 'G', North: false},
		},
		{
			name:      "Capetown, South Africa",
			latitude:  -33.92487,
			longitude: 18.42406,
			easting:   261877.816,
			northing:  6243185.589,
			zone:      Zone{Number: 34, Letter: 'H', North: false},
		},
		{
			name:      "Mendoza, Argentina",
			latitude:  -32.89018,
			longitude: -68.84405,
			easting:   514586.227,
			northing:  6360876.824,
			zone:      Zone{Number: 19, Letter: 'H', North: false},
		},
		{
			name:      "Fairbanks, Alaska, USA",
			latitude:  64.83778,
			longitude: -147.71639,
			easting:   466013.272,
			northing:  7190568,
			zone:      Zone{Number: 6, Letter: 'W', North: true},
		},
		{
			name:      "Ben Nevis, Scotland, UK",
			latitude:  56.79680,
			longitude: -5.00601,
			easting:   377485.765,
			northing:  6296561.854,
			zone:      Zone{Number: 30, Letter: 'V', North: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("LatLonZone", func(t *testing.T) {
				zone := LatLonZone(tt.latitude, tt.longitude)
				assert.DeepEqual(t, zone, tt.zone)
			})
			t.Run("ToUTM", func(t *testing.T) {
				easting, northing, zone := ToUTM(tt.latitude, tt.longitude)
				assert.DeepEqual(t, zone, tt.zone)
				assert.DeepEqual(t, easting, tt.easting, approx)
				assert.DeepEqual(t, northing, tt.northing, approx)
			})
			t.Run("ToLatLon", func(t *testing.T) {
				latitude, longitude := tt.zone.ToLatLon(tt.easting, tt.northing)
				assert.DeepEqual(t, latitude, tt.latitude, approx)
				assert.DeepEqual(t, longitude, tt.longitude, approx)
			})
		})
	}
}

func TestSRID(t *testing.T) {
	tests := []struct {
		srid int
		zone Zone
	}{
		{srid: 32610, zone: Zone{Number: 10, North: true}},
		{srid: 32659, zone: Zone{Number: 59, North: true}},
		{srid: 32734, zone: Zone{Number: 34, North: false}},
		{srid: 32701, zone: Zone{Number: 1, North: false}},
	}
	for _, tt := range tests {
		t.Run(tt.zone.String(), func(t *testing.T) {
			zone, ok := LookupSRID(tt.srid)
			assert.Assert(t, ok)
			assert.DeepEqual(t, zone, tt.zone)
			assert.Equal(t, zone.SRID(), tt.srid)
		})
	}
}

func TestParseZone(t *testing.T) {
	tests := []struct {
		input string
		zone  Zone
		valid bool
	}{
		{input: "45N", zone: Zone{Number: 45, Letter: 'N', North: true}, valid: true},
		{input: "12J", zone: Zone{Number: 12, Letter: 'J', North: false}, valid: true},
		{input: "3G", zone: Zone{Number: 3, Letter: 'G', North: false}, valid: true},
		{input: "5R", zone: Zone{Number: 5, Letter: 'R', North: true}, valid: true},
		{input: "5", zone: Zone{}, valid: false},
		{input: "RR", zone: Zone{}, valid: false},
		{input: "555R", zone: Zone{}, valid: false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			zone, ok := ParseZone(tt.input)
			assert.Equal(t, ok, tt.valid)
			if tt.valid {
				assert.DeepEqual(t, zone, tt.zone)
			}
		})
	}
}

func TestZoneValid(t *testing.T) {
	tests := []struct {
		zone  Zone
		valid bool
	}{
		{zone: Zone{Number: 1, Letter: 'S', North: true}, valid: true},
		{zone: Zone{Number: 8, Letter: 'S', North: false}, valid: false},
		{zone: Zone{Number: 70, Letter: 'S', North: true}, valid: false},
		{zone: Zone{Number: 34, Letter: 'O', North: true}, valid: false},
	}
	for _, tt := range tests {
		t.Run(tt.zone.String(), func(t *testing.T) {
			assert.Equal(t, tt.zone.Valid(), tt.valid)
		})
	}
}

func TestForcingAntiMeridian(t *testing.T) {
	// Force point just west of anti-meridian to east zone 1
	zone, _ := ParseZone("1N")
	easting, northing := zone.ToUTM(0, 179.9)
	_, lon := zone.ToLatLon(easting, northing)
	assert.Assert(t, math.Abs(179.9-lon) < 0.00001)

	// Force point just east of anti-meridian to west zone 60
	zone, _ = ParseZone("60N")
	easting, northing = zone.ToUTM(0, -179.9)
	_, lon = zone.ToLatLon(easting, northing)
	assert.Assert(t, math.Abs(-179.9-lon) < 0.00001)
}
