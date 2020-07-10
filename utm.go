package utm

import (
	"math"
)

const (
	_K0   = 0.9996
	_E    = 0.00669438
	_E2   = _E * _E
	_E3   = _E2 * _E
	_E_P2 = _E / (1.0 - _E)

	_SQRT_E = 0.9966471893303066 // math.Sqrt(1 - _E)

	__E  = (1 - _SQRT_E) / (1 + _SQRT_E)
	__E2 = __E * __E
	__E3 = __E2 * __E
	__E4 = __E3 * __E
	__E5 = __E4 * __E

	_M1 = 1 - _E/4 - 3*_E2/64 - 5*_E3/256
	_M2 = 3*_E/8 + 3*_E2/32 + 45*_E3/1024
	_M3 = 15*_E2/256 + 45*_E3/1024
	_M4 = 35 * _E3 / 3072

	_P2 = 3./2*__E - 27./32*__E3 + 269./512*__E5
	_P3 = 21./16*__E2 - 55./32*__E4
	_P4 = 151./96*__E3 - 417./128*__E5
	_P5 = 1097. / 512 * __E4

	_R = 6378137
)

// ToLatLon converts UTM coordinates to EPSG:4326 latitude/longitude
// Note: the zone's North field must be correctly set, the letter is ignored
func (zone Zone) ToLatLon(easting, northing float64) (latitude, longitude float64) {

	x := easting - 500000
	y := northing

	if !zone.North {
		y -= 10000000
	}

	m := y / _K0
	mu := m / (_R * _M1)

	p_rad := mu +
		_P2*math.Sin(2*mu) +
		_P3*math.Sin(4*mu) +
		_P4*math.Sin(6*mu) +
		_P5*math.Sin(8*mu)

	p_sin := math.Sin(p_rad)
	p_sin2 := p_sin * p_sin

	p_cos := math.Cos(p_rad)

	p_tan := p_sin / p_cos
	p_tan2 := p_tan * p_tan
	p_tan4 := p_tan2 * p_tan2

	ep_sin := 1 - _E*p_sin2
	ep_sin_sqrt := math.Sqrt(1 - _E*p_sin2)

	n := _R / ep_sin_sqrt
	rad := (1 - _E) / ep_sin

	c := __E * p_cos * p_cos
	c2 := c * c

	d := x / (n * _K0)
	d2 := d * d
	d3 := d2 * d
	d4 := d3 * d
	d5 := d4 * d
	d6 := d5 * d

	latitude = p_rad - (p_tan/rad)*
		(d2/2-
			d4/24*(5+3*p_tan2+10*c-4*c2-9*_E_P2)) +
		d6/720*(61+90*p_tan2+298*c+45*p_tan4-252*_E_P2-3*c2)

	longitude = (d -
		d3/6*(1+2*p_tan2+c) +
		d5/120*(5-2*c+28*p_tan2-3*c2+8*_E_P2+24*p_tan4)) / p_cos

	latitude = toDeg(latitude)
	longitude = toDeg(longitude) + zone.CentralMeridian()

	return latitude, longitude
}

// ToUTM convert a EPSG:4326 latitude/longitude to UTM.
func ToUTM(latitude, longitude float64) (easting, northing float64, zone Zone) {

	lat_rad := toRad(latitude)
	lat_sin := math.Sin(lat_rad)
	lat_cos := math.Cos(lat_rad)

	lat_tan := lat_sin / lat_cos
	lat_tan2 := lat_tan * lat_tan
	lat_tan4 := lat_tan2 * lat_tan2

	lon_rad := toRad(longitude)
	zone = LatLonZone(latitude, longitude)
	central_lon_rad := toRad(zone.CentralMeridian())

	n := _R / math.Sqrt(1-_E*math.Pow(lat_sin, 2))
	c := _E_P2 * lat_cos * lat_cos

	a := lat_cos * (lon_rad - central_lon_rad)
	a2 := a * a
	a3 := a2 * a
	a4 := a3 * a
	a5 := a4 * a
	a6 := a5 * a

	m := _R * (_M1*lat_rad -
		_M2*math.Sin(2*lat_rad) +
		_M3*math.Sin(4*lat_rad) -
		_M4*math.Sin(6*lat_rad))

	easting = _K0*n*(a+
		a3/6*(1-lat_tan2+c)+
		a5/120*(5-18*lat_tan2+lat_tan4+72*c-58*_E_P2)) + 500000

	northing = _K0 * (m + n*lat_tan*(a2/2+
		a4/24*(5-lat_tan2+9*c+4*math.Pow(c, 2))+
		a6/720*(61-58*lat_tan2+lat_tan4+600*c-330*_E_P2)))

	if latitude < 0 {
		northing += 10000000
	}

	return easting, northing, zone
}

func toDeg(r float64) float64 { return r / (math.Pi / 180) }
func toRad(d float64) float64 { return d * (math.Pi / 180) }
