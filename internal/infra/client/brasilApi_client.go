package client

import (
	"context"
	"encoding/json"
	"fmt"
	clientDto "github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/dto/client"
	"io"
	"net/http"
	"time"
)

type BrasilApiClient struct {
	URI string
}

func NewBrasilApiClient(uri string) *BrasilApiClient {
	return &BrasilApiClient{URI: uri}
}

func (client *BrasilApiClient) FindCep(cep int) (*clientDto.CepClientResponseDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*999)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/cep/v1/%v", client.URI, cep), nil)
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar requisição: %v\n", err)
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("Erro ao fazer requisição: %v\n", err)
	}
	defer res.Body.Close()

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("Erro %v ao fazer requisição\n", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler resposta: %v\n", err)
	}

	var data clientDto.BrasilApiCepResponseDto
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("Erro ao converter resposta: %v\n", err)
	}

	var result clientDto.CepClientResponseDto
	result.Cep = data.Cep
	result.Logradouro = data.Street
	result.Bairro = data.Neighborhood
	result.Cidade = data.City
	result.Uf = data.State

	return &result, nil
}
