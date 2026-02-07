package ip

import (
	"context"
	"errors"
	"log/slog"

	"github.com/artarts36/ip-scanner/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ipscannerapi "github.com/artarts36/ip-scanner/pkg/ip-scanner-grpc-api"
)

func (s *Service) ScanIP(ctx context.Context, req *ipscannerapi.ScanIPRequest) (*ipscannerapi.ScanIPResponse, error) {
	ip, err := s.ipRepository.Find(ctx, req.Address)
	if err != nil {
		var errIPInvalid *domain.InvalidIPError
		if errors.As(err, &errIPInvalid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "ip not found")
		}

		slog.ErrorContext(ctx, "unable to scan ip address", slog.String("ip", req.Address))

		return nil, status.Error(codes.Internal, "unable to scan ip address")
	}

	return &ipscannerapi.ScanIPResponse{
		Ip: &ipscannerapi.IP{
			Address: ip.Address,
			Country: &ipscannerapi.Country{
				IsoCode: ip.Country.ISOCode,
				Names:   ip.Country.Names,
			},
			City: &ipscannerapi.City{
				Name: ip.City.Name,
			},
			Location: &ipscannerapi.Location{
				Longitude: ip.Location.Longitude,
				Latitude:  ip.Location.Latitude,
			},
		},
	}, nil
}
