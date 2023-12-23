package webserver

import "github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/dto/client"

type FindCepResponseDto struct {
	APIOrigem     string                       `json:"api-origem"`
	DadosEndereco *client.CepClientResponseDto `json:"dados-endereco"`
}
