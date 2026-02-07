package app

import (
	"github.com/artarts36/ip-scanner/internal/port/grpc/handlers/ip"
	ipscannerapi "github.com/artarts36/ip-scanner/pkg/ip-scanner-grpc-api"
)

func (app *App) registerServices() {
	ipscannerapi.RegisterIPServiceServer(app.gRPCServer, ip.NewService(app.infrastructure.repositories.ipRepository))
}
