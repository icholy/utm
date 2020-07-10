# WGS48 UTM Conversion package

> Package for converting to and from the Universal Transverse Mercator coordinate system

## Examples:

### Lookup a zone by lat/lon or srid code
``` go
zone := utm.LatLonZone(50.77535, 6.008)
zone, _ := utm.LookupSRID(32601)
```

### convert from lat/lon to UTM
``` go
easting, northing, zone := utm.ToUTM(50.77535, 6.008)
```

### convert from UTM to lat/lon
``` go
latitude, longitude := utm.ToLatLon(294408.917, 5628897.997, zone)
```

## Credit:

This was mostly copied from: https://github.com/Turbo87/utm
