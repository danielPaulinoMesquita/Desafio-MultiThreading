package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApiCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	requestApiExternal()
}

func requestApiExternal() {
	cep := "01153000"
	chanelBrasilApi := make(chan BrasilApiCep)
	chanelViaCepApi := make(chan ViaCep)

	brasilApi := "https://brasilapi.com.br/api/cep/v1/" + cep
	viaApi := "https://viacep.com.br/ws/" + cep + "/json/"

	go func() {
		req, err := http.Get(viaApi)
		if err != nil {
			panic(err)
		}

		defer req.Body.Close()

		res, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		var resData ViaCep
		err = json.Unmarshal(res, &resData)
		if err != nil {
			panic(err)
		}

		chanelViaCepApi <- resData
	}()

	go func() {
		req, err := http.Get(brasilApi)
		if err != nil {
			panic(err)
		}

		defer req.Body.Close()

		res, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		var resData BrasilApiCep
		err = json.Unmarshal(res, &resData)
		if err != nil {
			panic(err)
		}

		chanelBrasilApi <- resData
	}()

	select {
	case chanelBrasilApiMsg := <-chanelBrasilApi:
		fmt.Printf("Resultado: %+v API: BrasilApi\n", chanelBrasilApiMsg)

	case chanelViaCepApiMsg := <-chanelViaCepApi:
		fmt.Printf("Resultado: %+v API: ViaCep\n", chanelViaCepApiMsg)

	case <-time.After(time.Second * 1):
		println("Time Out!")
	}

}
