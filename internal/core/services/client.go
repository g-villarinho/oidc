package services

import "github.com/g-villarinho/oidc-server/internal/core/ports"

type ClientService struct {
	clientRepository ports.ClientRepository
}

func NewClientService(clientRepository ports.ClientRepository) *ClientService {
	return &ClientService{
		clientRepository: clientRepository,
	}
}
