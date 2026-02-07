package domain

import "context"

type IPRepository interface {
	Find(ctx context.Context, ip string) (*IP, error)
}

type IP struct {
	Address  string
	Country  Country
	City     City
	Location Location
}

type Country struct {
	ISOCode string

	Names map[string]string
}

type City struct {
	Name string
}

type Location struct {
	Longitude float64
	Latitude  float64
}
