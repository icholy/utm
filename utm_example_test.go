package utm_test

import (
	"fmt"

	"github.com/icholy/utm"
)

func ExampleLookupSRID() {
	if zone, ok := utm.LookupSRID(32617); ok {
		fmt.Println(zone)
	}
	// Output: 17 (north)
}

func ExampleLatLonZone() {
	fmt.Println(utm.LatLonZone(50.77535, 6.008))
	// Output: 32 (north)
}

func ExampleToUTM() {
	easting, northing, zone := utm.ToUTM(50.77535, 6.008)
	fmt.Println("Zone:", zone)
	fmt.Printf("Easting: %f\n", easting)
	fmt.Printf("Northing: %f\n", northing)
	// Output: Zone: 32 (north)
	// Easting: 289059.493943
	// Northing: 5629111.846925
}

func ExampleToLatLon() {
	zone, _ := utm.LookupSRID(32632)
	latitude, longitude := utm.ToLatLon(289059.493943, 5629111.846925, zone)
	fmt.Printf("Latitude: %f\n", latitude)
	fmt.Printf("Longitude: %f\n", longitude)
	// Output: Latitude: 50.775350
	// Longitude: 6.007999
}
