package main

import (
	"errors"
	"fmt"
)

var ErroNomeVazio = errors.New("O nome informado está vazio.")

func DigaOla(nome string) (string, error) {
	if len(nome) == 0 {
		return "", ErroNomeVazio
	}
	return "Olá " + nome, nil
}

func main() {
	if str, err := DigaOla(""); err == ErroNomeVazio {
		fmt.Printf("Ocorreu um erro ao dizer Olá: %s\n", err)
	} else {
		fmt.Println(str)
	}
}
