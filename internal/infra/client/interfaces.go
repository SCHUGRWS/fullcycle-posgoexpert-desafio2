package client

import (
	clientDto "github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/dto/client"
)

type CepApiClient interface {
	FindCep(cep string) (*clientDto.CepClientResponseDto, error)
}
