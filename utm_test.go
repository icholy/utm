package utm

import (
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
			zone:      Zone{N: 32, North: true},
		},
		{
			name:      "New York, USA",
			latitude:  40.71435,
			longitude: -74.00597,
			easting:   583959.959,
			northing:  4507523.087,
			zone:      Zone{N: 18, North: true},
		},
		{
			name:      "Wellington, New Zealand",
			latitude:  -41.28646,
			longitude: 174.77624,
			easting:   313784.305,
			northing:  5427057.321,
			zone:      Zone{N: 60, North: false},
		},
		{
			name:      "Capetown, South Africa",
			latitude:  -33.92487,
			longitude: 18.42406,
			easting:   261877.816,
			northing:  6243185.589,
			zone:      Zone{N: 34, North: false},
		},
		{
			name:      "Mendoza, Argentina",
			latitude:  -32.89018,
			longitude: -68.84405,
			easting:   514586.227,
			northing:  6360876.824,
			zone:      Zone{N: 19, North: false},
		},
		{
			name:      "Fairbanks, Alaska, USA",
			latitude:  64.83778,
			longitude: -147.71639,
			easting:   466013.272,
			northing:  7190568,
			zone:      Zone{N: 6, North: true},
		},
		{
			name:      "Ben Nevis, Scotland, UK",
			latitude:  56.79680,
			longitude: -5.00601,
			easting:   377485.765,
			northing:  6296561.854,
			zone:      Zone{N: 30, North: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("ToUTM", func(t *testing.T) {
				easting, northing, zone := ToUTM(tt.latitude, tt.longitude)
				assert.DeepEqual(t, zone, tt.zone)
				assert.DeepEqual(t, easting, tt.easting, approx)
				assert.DeepEqual(t, northing, tt.northing, approx)
			})
			t.Run("ToLatLon", func(t *testing.T) {
				latitude, longitude := ToLatLon(tt.easting, tt.northing, tt.zone)
				assert.DeepEqual(t, latitude, tt.latitude, approx)
				assert.DeepEqual(t, longitude, tt.longitude, approx)
			})
		})
	}
}
