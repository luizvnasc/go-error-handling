package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"errors"
)

type StringVaziaError string

var(
	errNomeVazio = errors.New("O nome informado está vazio")
	errNomeNumerico = errors.New("O nome informado é um número")
	errNomeCaracterEspecial = errors.New("O nome informado contém caracteres especiais")
)

// BemVindo constrói uma menssagem de boas vindas desejada para um nome passado por parâmetro.
func BemVindo(nome string) (string, error) {
	// verifica se a string é vazia
	if s := strings.Trim(nome, " "); len(s) == 0 {
		return "", errNomeVazio
	}
	// verifica se a string possui apenas números
	if _, err := strconv.ParseFloat(nome, 64); err == nil {
		return "", errNomeNumerico
	}
	// verifica se a string possui caracteres especiais
	if strings.ContainsAny(nome, `,.|!@#$%&*+_-=[]{};:/?\\'"()`) {
		return "", errNomeCaracterEspecial
	}

	return "Bem Vindo ao meetup da comunidade Golang CWB, " + nome + ".", nil
}

// DigaBemVindo imprime uma mensagem de bem vindo para um participante do meetup.
func DigaBemVindo(w io.Writer, nome string) error {
	msgBoasVindas, err := BemVindo(nome)
	if err != nil {
		return fmt.Errorf("Erro ao criar mensagem de boas vindas: %w", err)
	}
	fmt.Fprintln(w, msgBoasVindas)
	return nil
}

func main() {
	nome := flag.String("nome", "folks", "Nome do participante do meetup")
	flag.Parse()
	err := DigaBemVindo(os.Stdout, *nome)

	if err != nil {
		log.Println(err)
		if errors.Is(err, errNomeVazio){
			fmt.Fprintln(os.Stdout, "Não aceitamos pessoas anônimas!")
		}
		if errors.Is(err, errNomeNumerico){
			fmt.Fprintln(os.Stdout, "Te entendo, somos todos apenas números.")
		}
		if errors.Is(err, errNomeCaracterEspecial){
			fmt.Fprintln(os.Stdout, "Você ainda usa hotmail?")
		}
	}
}
