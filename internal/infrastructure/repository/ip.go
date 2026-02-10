package repository

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/artarts36/ip-scanner/internal/infrastructure/storage"

	"github.com/artarts36/ip-scanner/internal/domain"
)

type IPRepository struct {
	db *storage.DB
}

func NewIPRepository(mmdb *storage.DB) *IPRepository {
	return &IPRepository{db: mmdb}
}

func (r *IPRepository) Find(_ context.Context, address string) (*domain.IP, error) {
	ipAddr, err := netip.ParseAddr(address)
	if err != nil {
		return nil, domain.NewErrIPInvalid(err)
	}

	// https://github.com/oschwald/geoip2-golang/blob/main/reader.go
	type record struct {
		Country struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`

		Location struct {
			Longitude float64 `maxminddb:"longitude"`
			Latitude  float64 `maxminddb:"latitude"`
		} `maxminddb:"location"`

		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
	}

	var rec record

	res := r.db.DB.Lookup(ipAddr)
	if err = res.Err(); err != nil {
		return nil, fmt.Errorf("failed to lookup address address: %w", err)
	}
	if !res.Found() {
		return nil, domain.ErrNotFound
	}

	err = res.Decode(&rec)
	if err != nil {
		return nil, fmt.Errorf("failed to decode address address: %w", err)
	}

	var city string
	err = res.DecodePath(&city, "location", "city", "names", "en")
	if err != nil {
		return nil, fmt.Errorf("failed to decode city: %w", err)
	}

	return &domain.IP{
		Address: address,

		Country: domain.Country{
			ISOCode: rec.Country.ISOCode,
			Names:   rec.Country.Names,
		},

		City: domain.City{
			Name: rec.City.Names["en"],
		},

		Location: domain.Location{
			Longitude: rec.Location.Longitude,
			Latitude:  rec.Location.Latitude,
		},
	}, nil
}
