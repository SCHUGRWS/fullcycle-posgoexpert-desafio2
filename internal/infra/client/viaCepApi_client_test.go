package client

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestViaCepApiClient_FindCep(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/12345678/json/",
		httpmock.NewStringResponder(200, `{
			"cep": "01001-000",
			"logradouro": "Praça da Sé",
			"complemento": "lado ímpar",
			"bairro": "Sé",
			"localidade": "São Paulo",
			"uf": "SP",
			"ibge": "3550308",
			"gia": "1004",
			"ddd": "11",
			"siafi": "7107"
		}`))

	viaCepApiClient := NewViaCepClient("https://viacep.com.br/ws")
	result, err := viaCepApiClient.FindCep(12345678)
	assert.NotNil(t, result)
	assert.NoError(t, err)
	assert.Equal(t, "01001-000", result.Cep)
	assert.Equal(t, "São Paulo", result.Cidade)
	assert.Equal(t, "Sé", result.Bairro)
	assert.Equal(t, "Praça da Sé", result.Logradouro)
	assert.Equal(t, "SP", result.Uf)
}

func TestViaCepApiClient_FindCep_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/12345678/json/",
		httpmock.NewStringResponder(500, `{
			"erro": true
		}`))

	viaCepApiClient := NewViaCepClient("https://viacep.com.br/ws")
	result, err := viaCepApiClient.FindCep(12345678)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "Erro 500 ao fazer requisição\n", err.Error())
}

func TestViaCepApiClient_FindCep_Timeout(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/12345678/json/",
		httpmock.NewStringResponder(200, `{
			"cep": "01001-000",
			"logradouro": "Praça da Sé",
			"complemento": "lado ímpar",
			"bairro": "Sé",
			"localidade": "São Paulo",
			"uf": "SP",
			"ibge": "3550308",
			"gia": "1004",
			"ddd": "11",
			"siafi": "7107"
		}`).Delay(1000*time.Millisecond))

	viaCepApiClient := NewViaCepClient("https://viacep.com.br/ws")
	result, err := viaCepApiClient.FindCep(12345678)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "Erro ao fazer requisição: Get \"https://viacep.com.br/ws/12345678/json/\": context deadline exceeded\n", err.Error())
}
