## Projeto: API de Consulta de CEP em Golang

### Descrição
Este projeto é um desafio que utiliza multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas de consulta de CEP.
qs duas requisições são feitas simultaneamente para as seguintes APIs:

- BrasilAPI
- ViaCEP

O programa acata a resposta da API mais rápida e descarta a resposta mais lenta, exibindo o resultado no terminal com os dados do endereço e a API que enviou a resposta.
caso nenhuma das APIs responda dentro de 1 segundo, um erro de timeout é exibido.
