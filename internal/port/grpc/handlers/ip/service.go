package ip

import (
	"github.com/artarts36/ip-scanner/internal/domain"
	ipscannerapi "github.com/artarts36/ip-scanner/pkg/ip-scanner-grpc-api"
)

type Service struct {
	ipscannerapi.UnimplementedIPServiceServer

	ipRepository domain.IPRepository
}

func NewService(ipRepository domain.IPRepository) *Service {
	return &Service{ipRepository: ipRepository}
}
