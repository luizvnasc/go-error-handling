package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type StringVaziaError string

func (s StringVaziaError) Error() string {
	return "A string está vazia."
}

type StringNumericaError string

func (s StringNumericaError) Error() string {
	return "A string " + string(s) + " contém apenas números."
}

type StringComCaracteresEspeciaisError string

func (s StringComCaracteresEspeciaisError) Error() string {
	return "A string " + string(s) + " contém apenas caracteres especiais."
}

// BemVindo constrói uma menssagem de boas vindas desejada para um nome passado por parâmetro.
func BemVindoCustom(nome string) (string, error) {
	// verifica se a string é vazia
	if s := strings.Trim(nome, " "); len(s) == 0 {
		return "", StringVaziaError(nome)
	}
	// verifica se a string possui apenas números
	if _, err := strconv.ParseFloat(nome, 64); err == nil {
		return "", StringNumericaError(nome)
	}
	// verifica se a string possui caracteres especiais
	if strings.ContainsAny(nome, `,.|!@#$%&*+_-=[]{};:/?\\'"()`) {
		return "", StringComCaracteresEspeciaisError(nome)
	}

	return "Bem Vindo ao meetup da comunidade Golang CWB, " + nome + ".", nil
}

// DigaBemVindoCustom imprime uma mensagem de bem vindo para um participante do meetup.
func DigaBemVindoCustom(w io.Writer, nome string) {
	msgBoasVindas, err := BemVindoCustom(nome)
	if err != nil {
		switch err.(type) {
		case StringVaziaError:
			fmt.Fprintln(w, "Não aceitamos pessoas anônimas!")
		case StringNumericaError:
			fmt.Fprintln(w, "Te entendo, somos todos apenas números.")
		case StringComCaracteresEspeciaisError:
			fmt.Fprintln(w, "Você ainda usa hotmail?")
		}
	}
	fmt.Fprintln(w, msgBoasVindas)
}

// BemVindo constrói uma menssagem de boas vindas desejada para um nome passado por parâmetro.
func BemVindo(nome string) (string, error) {
	// verifica se a string é vazia
	if s := strings.Trim(nome, " "); len(s) == 0 {
		return "", errors.New("O nome informado está vazio")
	}
	// verifica se a string possui apenas números
	if _, err := strconv.ParseFloat(nome, 64); err == nil {
		return "", errors.New("O nome informado é um número")
	}
	// verifica se a string possui caracteres especiais
	if strings.ContainsAny(nome, `,.|!@#$%&*+_-=[]{};:/?\\'"()`) {
		return "", errors.New("O nome informado contém caracteres especiais")
	}

	return "Bem Vindo ao meetup da comunidade Golang CWB, " + nome + ".", nil
}

// DigaBemVindo imprime uma mensagem de bem vindo para um participante do meetup.
func DigaBemVindo(w io.Writer, nome string) {
	msgBoasVindas, err := BemVindo(nome)
	if err != nil {
		switch err.Error() {
		case "O nome informado está vazio":
			fmt.Fprintln(w, "Não aceitamos pessoas anônimas!")
		case "O nome informado é um número":
			fmt.Fprintln(w, "Te entendo, somos todos apenas números.")
		case "O nome informado contém caracteres especiais":
			fmt.Fprintln(w, "Você ainda usa hotmail?")
		}
	}
	fmt.Fprintln(w, msgBoasVindas)
}

func main() {
	nome := flag.String("nome", "folks", "Nome do participante do meetup")
	flag.Parse()
	DigaBemVindo(os.Stdout, *nome)
}
