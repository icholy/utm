package utm

import (
	"math"
)

const (
	_K0  = 0.9996
	_E   = 0.00669438
	_E2  = _E * _E
	_E3  = _E2 * _E
	_EP2 = _E / (1.0 - _E)

	_SqrtE = 0.9966471893303066 // math.Sqrt(1 - _E)

	_xE  = (1 - _SqrtE) / (1 + _SqrtE)
	_xE2 = _xE * _xE
	_xE3 = _xE2 * _xE
	_xE4 = _xE3 * _xE
	_xE5 = _xE4 * _xE

	_M1 = 1 - _E/4 - 3*_E2/64 - 5*_E3/256
	_M2 = 3*_E/8 + 3*_E2/32 + 45*_E3/1024
	_M3 = 15*_E2/256 + 45*_E3/1024
	_M4 = 35 * _E3 / 3072

	_P2 = 3./2*_xE - 27./32*_xE3 + 269./512*_xE5
	_P3 = 21./16*_xE2 - 55./32*_xE4
	_P4 = 151./96*_xE3 - 417./128*_xE5
	_P5 = 1097. / 512 * _xE4

	_R = 6378137
)

// ToLatLon converts UTM coordinates to EPSG:4326 latitude/longitude
// Note: the zone's North field must be correctly set, the letter is ignored
func (z Zone) ToLatLon(easting, northing float64) (latitude, longitude float64) {

	x := easting - 500000
	y := northing

	if !z.North {
		y -= 10000000
	}

	m := y / _K0
	mu := m / (_R * _M1)

	pRad := mu +
		_P2*math.Sin(2*mu) +
		_P3*math.Sin(4*mu) +
		_P4*math.Sin(6*mu) +
		_P5*math.Sin(8*mu)

	pSin := math.Sin(pRad)
	pSin2 := pSin * pSin

	pCos := math.Cos(pRad)

	pTan := pSin / pCos
	pTan2 := pTan * pTan
	pTan4 := pTan2 * pTan2

	epSin := 1 - _E*pSin2
	epSinSqrt := math.Sqrt(1 - _E*pSin2)

	n := _R / epSinSqrt
	rad := (1 - _E) / epSin

	c := _EP2 * pCos * pCos
	c2 := c * c

	d := x / (n * _K0)
	d2 := d * d
	d3 := d2 * d
	d4 := d3 * d
	d5 := d4 * d
	d6 := d5 * d

	latitude = pRad - (pTan/rad)*
		(d2/2-
			d4/24*(5+3*pTan2+10*c-4*c2-9*_EP2)) +
		d6/720*(61+90*pTan2+298*c+45*pTan4-252*_EP2-3*c2)

	longitude = (d -
		d3/6*(1+2*pTan2+c) +
		d5/120*(5-2*c+28*pTan2-3*c2+8*_EP2+24*pTan4)) / pCos

	longitude = modAngle(longitude + toRad(z.CentralMeridian()))

	latitude = toDeg(latitude)
	longitude = toDeg(longitude)

	return latitude, longitude
}

// ToUTM convert a EPSG:4326 latitude/longitude to UTM.
func ToUTM(latitude, longitude float64) (easting, northing float64, zone Zone) {
	zone = LatLonZone(latitude, longitude)
	easting, northing = zone.ToUTM(latitude, longitude)
	return easting, northing, zone
}

// ToUTM convert a EPSG:4326 latitude/longitude to UTM.
func (z Zone) ToUTM(latitude, longitude float64) (easting, northing float64) {

	latRad := toRad(latitude)
	latSin := math.Sin(latRad)
	latCos := math.Cos(latRad)

	latTan := latSin / latCos
	latTan2 := latTan * latTan
	latTan4 := latTan2 * latTan2

	lonRad := toRad(longitude)
	centralLonRad := toRad(z.CentralMeridian())

	n := _R / math.Sqrt(1-_E*math.Pow(latSin, 2))
	c := _EP2 * latCos * latCos

	a := latCos * modAngle(lonRad-centralLonRad)
	a2 := a * a
	a3 := a2 * a
	a4 := a3 * a
	a5 := a4 * a
	a6 := a5 * a

	m := _R * (_M1*latRad -
		_M2*math.Sin(2*latRad) +
		_M3*math.Sin(4*latRad) -
		_M4*math.Sin(6*latRad))

	easting = _K0*n*(a+
		a3/6*(1-latTan2+c)+
		a5/120*(5-18*latTan2+latTan4+72*c-58*_EP2)) + 500000

	northing = _K0 * (m + n*latTan*(a2/2+
		a4/24*(5-latTan2+9*c+4*math.Pow(c, 2))+
		a6/720*(61-58*latTan2+latTan4+600*c-330*_EP2)))

	if latitude < 0 {
		northing += 10000000
	}

	return easting, northing
}

func toDeg(r float64) float64 { return r / (math.Pi / 180) }
func toRad(d float64) float64 { return d * (math.Pi / 180) }

// modAngle returns angle in radians to be between -pi and pi
func modAngle(value float64) float64 {
	return mod(value+math.Pi, 2.0*math.Pi) - math.Pi
}

// mod acts like python's % operator
func mod(d, m float64) float64 {
	res := math.Mod(d, m)
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}
