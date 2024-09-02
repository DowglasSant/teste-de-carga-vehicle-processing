package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Vehicle struct {
	Type             int    `json:"type"`
	LicensePlate     string `json:"license_plate"`
	Year             int    `json:"year"`
	OwnerCPF         string `json:"owner_cpf"`
	RegistrationCity string `json:"registration_city"`
	Color            string `json:"color"`
	Brand            string `json:"brand"`
	Model            string `json:"model"`
}

const (
	url            = "http://localhost:8000/test"
	numRequisicoes = 1000
	numGoroutines  = 50
)

func gerarCPF(id int) string {
	cpf := fmt.Sprintf("12345678%02d", id%100)
	return cpf
}

func fazerRequisicao(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	vehicle := Vehicle{
		Type:             1,
		LicensePlate:     fmt.Sprintf("ABC-%04d", id),
		Year:             2020,
		OwnerCPF:         gerarCPF(id),
		RegistrationCity: "São Paulo",
		Color:            "Preto",
		Brand:            "Toyota",
		Model:            "Corolla",
	}

	vehicleJSON, err := json.Marshal(vehicle)
	if err != nil {
		fmt.Printf("Erro ao converter o veículo %d para JSON: %v\n", id, err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(vehicleJSON))
	if err != nil {
		fmt.Printf("Erro na requisição %d: %v\n", id, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Requisição %d: %s\n", id, resp.Status)
}

func main() {
	var wg sync.WaitGroup
	inicio := time.Now()

	sem := make(chan struct{}, numGoroutines)

	for i := 0; i < numRequisicoes; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer func() { <-sem }()
			fazerRequisicao(id, &wg)
		}(i)
	}

	wg.Wait()
	close(sem)

	fim := time.Now()
	fmt.Printf("Teste concluído em %.2f segundos\n", fim.Sub(inicio).Seconds())
}
