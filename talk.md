# Tratamento de erros em go antes da versão 1.13

## Built-in

```golang
f, err := os.Open("filename.ext")
if err != nil {
    log.Fatal(err)
}
```
### O tipo `error`

```golang
type error interface {
    Error() string
}
```
O tipo `error`, assim como outros tipos interno de go, é [pré-declarado](https://golang.org/ref/spec#Predeclared_identifiers) no [bloco universal](https://golang.org/ref/spec#Blocks), que engloba todo o código fonte de Go.

## Implemetação do tipo `error`

A implementação mais trival de um tipo `error` é a implementação não exportada do pacote [`errors`](https://golang.org/src/errors/errors.go)

```golang
// errorString is a trivial implementation of error.
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

Ela pode ser construída a partir da função `errors.New`:

```golang
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}
```

Isso acaba se tornando uma dor de cabeça para quem está começando em go pois ao tentar realizar a comparação abaixo o resultado não é bem o esperado.

```golang

func main() {
erro1 := errors.New("Isto é um erro")
erro2 := errors.New("Isto é um erro")
fmt.Println(erro1 == erro2) //false
}
```
* https://play.golang.org/p/fqPUo8Rgglv3O

### Sentinelas

```golang
package main

import (
	"fmt"
	"errors"
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
		fmt.Printf("Ocorreu um erro ao dizer Olá: %s", err)
	}else{
		fmt.Println(str)
	}
}
````
* https://play.golang.org/p/Rt16cZiSJjD

### Tipos de erros Customizados

```golang
type StringVaziaError string

func (s StringVaziaError) Error() string {
	return "A string está vazia."
}

type StringNumericaError string

func (s StringNumericaError) Error() string {
	return "A string está " + string(s) + "contém apenas números."
}

type StringComCaracteresEspeciaisError string

func (s StringComCaracteresEspeciaisError) Error() string {
	return "A string está " + string(s) + "contém apenas números."
}
```

Além de mensagens customizadas é possível tratar os erros a partir do seu tipo:

```golang
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
```

### fontes
* [Error handling and Go - The go blog](https://blog.golang.org/error-handling-and-go)