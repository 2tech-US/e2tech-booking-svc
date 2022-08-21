package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/kelvins/geocoder"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func LocationToAddress(location Location) (string, error) {
	geLocation := geocoder.Location{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}

	addresses, err := geocoder.GeocodingReverse(geLocation)
	if err != nil {
		return "", err
	}

	return addresses[0].FormattedAddress, nil
}

type Address struct {
	HouseNumber string `json:"house_number"`
	Street      string `json:"street"`
	Ward        string `json:"ward"`
	District    string `json:"district"`
	City        string `json:"city"`
}

func AddressToLocation(address Address) (Location, error) {
	houseNumber, err := strconv.Atoi(address.HouseNumber)
	if err != nil {
		return Location{}, fmt.Errorf("invalid house number: %w", err)
	}

	geAddress := geocoder.Address{
		Number:   houseNumber,
		Street:   address.Street,
		District: address.District,
		City:     address.City,
	}

	geLocation, err := geocoder.Geocoding(geAddress)
	if err != nil {
		return Location{}, err
	}

	return Location{
		Latitude:  geLocation.Latitude,
		Longitude: geLocation.Longitude,
	}, nil
}

const grabPrice = 0.68

func CalculatePrice(pickup_lat, pickup_lng, dropoff_lat, dropoff_lng float64) float64 {
	distance := math.Sqrt(
		math.Pow((pickup_lat-dropoff_lat)*69.1, 2) +
			math.Pow((pickup_lng-dropoff_lng*69.1*math.Cos(pickup_lat/57.3)), 2),
	)
	return distance * grabPrice
}
