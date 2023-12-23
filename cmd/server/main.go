package main

import (
	"github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/configs"
	_ "github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/docs"
	"github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/infra/client"
	"github.com/SCHUGRWS/fullcycle-posgoexpert-desafio2/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title GO Expert Desafio 2
// @description This is a cep finder API
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @contact.name:Richard Schug
// @contact.email:richard.schug@outlook.com

// @host localhost:8000
// @BasePath /
func main() {
	config := configs.NewConfig()

	brasilApiClient := client.NewBrasilApiClient(config.BrasilAPIURL)
	viaCepApiClient := client.NewViaCepClient(config.ViaCEPURL)
	cepHandler := handlers.NewCepHandler(brasilApiClient, viaCepApiClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	//Product routes
	r.Route("/busca-cep", func(r chi.Router) {
		r.Get("/{cep}", cepHandler.FindCep)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}

}
