package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AddressViaCEP struct {
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

type AddressBrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type Result struct {
	AddressViaCEP    AddressViaCEP
	AddressBrasilAPI AddressBrasilAPI
	Source           string
}

// fetchFromBrasilAPI é uma função que busca o endereço de um CEP na API do BrasilAPI
func fetchFromBrasilAPI(cep string, ch chan<- Result) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	// Criando um cliente HTTP com timeout de 1 segundo
	client := http.Client{Timeout: 1 * time.Second}

	// Fazendo uma requisição GET para a URL
	resp, err := client.Get(url)
	if err != nil {
		ch <- Result{}
		return
	}
	defer resp.Body.Close()

	// Decodificando o corpo da resposta em um objeto Address
	var address AddressBrasilAPI
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		// Se houver erro, envia um resultado vazio para o canal
		ch <- Result{}
		return
	}

	// Enviando o resultado para o canal
	ch <- Result{AddressBrasilAPI: address, Source: "BrasilAPI"}
}

// fetchFromViaCEP é uma função que busca o endereço de um CEP na API do ViaCEP
func fetchFromViaCEP(cep string, ch chan<- Result) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	// Criando um cliente HTTP com timeout de 1 segundo
	client := http.Client{Timeout: 1 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		// Se houver erro, envia um resultado vazio para o canal
		ch <- Result{}
		return
	}
	defer resp.Body.Close()

	// Decodificando o corpo da resposta em um objeto Address
	var address AddressViaCEP
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		// Se houver erro, envia um resultado vazio para o canal
		ch <- Result{}
		return
	}

	// Enviando o resultado para o canal
	ch <- Result{AddressViaCEP: address, Source: "ViaCEP"}
}

func main() {
	cep := "24930024"

	// Criando um canal para receber o resultado
	ch := make(chan Result)

	// Iniciando as goroutines para buscar o endereço
	go fetchFromBrasilAPI(cep, ch)
	go fetchFromViaCEP(cep, ch)

	select {
	// Se o canal for preenchido, exibe o resultado
	case result := <-ch:
		if (result != Result{}) {
			fmt.Printf("Resposta mais rápida da API: %s\n", result.Source)
			if result.Source == "BrasilAPI" {
				fmt.Printf("Endereço: %+v\n", result.AddressBrasilAPI)
			} else {
				fmt.Printf("Endereço: %+v\n", result.AddressViaCEP)
			}
		} else {
			fmt.Println("Erro: Falha ao buscar o endereço.")
		}
	// Se o tempo limite for atingido, exibe uma mensagem de erro
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: Timeout ao buscar o endereço.")
	}
}
