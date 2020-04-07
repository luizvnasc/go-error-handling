package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
func BemVindo(nome string) (string, error) {
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
func DigaBemVindo(w io.Writer, nome string) error {
	msgBoasVindas, err := BemVindo(nome)
	if err != nil {
		return errors.Wrap(err, "Erro ao criar mensagem de boas vindas")
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
		switch errors.Cause(err).(type) {
		case StringVaziaError:
			fmt.Fprintln(os.Stdout, "Não aceitamos pessoas anônimas!")
		case StringNumericaError:
			fmt.Fprintln(os.Stdout, "Te entendo, somos todos apenas números.")
		case StringComCaracteresEspeciaisError:
			fmt.Fprintln(os.Stdout, "Você ainda usa hotmail?")
		}
	}
}
