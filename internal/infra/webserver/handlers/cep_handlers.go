package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/dto"
	clientDto "github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/dto/client"
	webserverDto "github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/dto/webserver"
	"github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/infra/client"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

type CepHandler struct {
	BrasilApiClient client.CepApiClient
	ViaCepApiClient client.CepApiClient
}

func NewCepHandler(brasilApiClient client.CepApiClient, viaCepApiClient client.CepApiClient) *CepHandler {
	return &CepHandler{
		BrasilApiClient: brasilApiClient,
		ViaCepApiClient: viaCepApiClient,
	}
}

// FindCep godoc
// @Summary Find a cep
// @Description Find a cep
// @Tags Cep
// @Accept  json
// @Produce  json
// @Param cep path int true "Cep"
// @Success 200 {object} webserverDto.FindCepResponseDto
// @Failure 500 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 408 {object} dto.Error
// @Router /busca-cep/{cep} [get]
func (c *CepHandler) FindCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := dto.Error{Message: "Cep é obrigatório"}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	if len(cep) != 8 {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := dto.Error{Message: "Cep deve conter 8 dígitos"}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}
	_, err := strconv.Atoi(cep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := dto.Error{Message: "Cep inválido"}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	chanelBrasilAPI := make(chan *clientDto.CepClientResponseDto)
	chanelViaCep := make(chan *clientDto.CepClientResponseDto)

	go func() {
		result, err := c.BrasilApiClient.FindCep(cep)
		if err != nil {
			fmt.Println(fmt.Errorf("Erro ao buscar CEP BrasilApi: %v\n", err))
			return
		}
		chanelBrasilAPI <- result
	}()

	go func() {
		result, err := c.ViaCepApiClient.FindCep(cep)
		if err != nil {
			fmt.Println(fmt.Errorf("Erro ao buscar CEP ViaCEP: %v\n", err))
			return
		}
		chanelViaCep <- result
	}()

	var result webserverDto.FindCepResponseDto

	for {
		select {
		case cepClientResponse := <-chanelBrasilAPI:
			result.APIOrigem = "BrasilAPI"
			result.DadosEndereco = cepClientResponse

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorMessage := dto.Error{Message: err.Error()}
				err = json.NewEncoder(w).Encode(errorMessage)
			} else {
				fmt.Printf("API: %+v\n", result.APIOrigem)
				fmt.Printf("Dados Endereço: %+v\n", result.DadosEndereco)
			}
			return
		case cepClientResponse := <-chanelViaCep:
			result.APIOrigem = "ViaCepAPI"
			result.DadosEndereco = cepClientResponse

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorMessage := dto.Error{Message: err.Error()}
				err = json.NewEncoder(w).Encode(errorMessage)
			} else {
				fmt.Printf("API: %+v\n", result.APIOrigem)
				fmt.Printf("Dados Endereço: %+v\n", result.DadosEndereco)
			}
			return
		case <-time.After(1 * time.Second):
			//w.WriteHeader(http.StatusRequestTimeout) não retorna no swagger, fica em loop, mas a aplicação retorna corretamente
			w.WriteHeader(http.StatusRequestTimeout)
			errorMessage := dto.Error{Message: "Tempo de requisição excedido"}
			err = json.NewEncoder(w).Encode(errorMessage)
			return
		}
	}

}
