package client

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBrasilApiCepClient_FindCep(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://brasilapi.com.br/api/cep/v1/12345678",
		httpmock.NewStringResponder(200, `{
			"cep": "01001-000",
			"state": "SP",
			"city": "São Paulo",
			"neighborhood": "Sé",
			"street": "Praça da Sé",
			"service": "correios"
		}`))

	brasilApiClient := NewBrasilApiClient("https://brasilapi.com.br/api")
	result, err := brasilApiClient.FindCep(12345678)
	assert.NotNil(t, result)
	assert.NoError(t, err)
	assert.Equal(t, "01001-000", result.Cep)
	assert.Equal(t, "São Paulo", result.Cidade)
	assert.Equal(t, "Sé", result.Bairro)
	assert.Equal(t, "Praça da Sé", result.Logradouro)
	assert.Equal(t, "SP", result.Uf)
}

func TestBrasilApiCepClient_FindCep_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://brasilapi.com.br/api/cep/v1/12345678",
		httpmock.NewStringResponder(500, `{
			"code": 500,
			"message": "Internal Server Error"
		}`))

	brasilApiClient := NewBrasilApiClient("https://brasilapi.com.br/api")
	result, err := brasilApiClient.FindCep(12345678)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "Erro 500 ao fazer requisição\n", err.Error())
}

func TestBrasilApiCepClient_FindCep_Timeout(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://brasilapi.com.br/api/cep/v1/12345678",
		httpmock.NewStringResponder(200, `{
			"cep": "01001-000",
			"state": "SP",
			"city": "São Paulo",
			"neighborhood": "Sé",
			"street": "Praça da Sé",
			"service": "correios"
		}`).Delay(1000*time.Millisecond))

	brasilApiClient := NewBrasilApiClient("https://brasilapi.com.br/api")
	result, err := brasilApiClient.FindCep(12345678)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "Erro ao fazer requisição: Get \"https://brasilapi.com.br/api/cep/v1/12345678\": context deadline exceeded\n", err.Error())
}
