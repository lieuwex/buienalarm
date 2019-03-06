package location

import (
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/locationiq"
)

const locationiqKey = ""

func init() {
	if locationiqKey == "" {
		panic("locationiqKey is empty")
	}
}

func Geocode(addr string) (*geo.Location, error) {
	coder := locationiq.Geocoder(locationiqKey, 18)
	return coder.Geocode(addr)
}
