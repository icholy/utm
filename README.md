# WGS84 UTM Conversion

> Package for converting to and from the Universal Transverse Mercator coordinate system

## Examples:

**Lookup a zone by lat/lon, srid code, or text**
``` go
zone := utm.LatLonZone(50.77535, 6.008)
zone, _ := utm.LookupSRID(32601)
zone, _ := utm.ParseZone("32U")
```

**Convert from lat/lon to UTM**
``` go
easting, northing, zone := utm.ToUTM(50.77535, 6.008)
```

**Convert from UTM to lat/lon**
``` go
latitude, longitude := zone.ToLatLon(294408.917, 5628897.997)
```

## Credit:

This was mostly copied from: https://github.com/Turbo87/utm
