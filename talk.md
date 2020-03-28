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
func comparaErros(erro1,erro2 error) bool {
    return erro1 == erro2
}

func main() {
erro1 := errors.New("Isto é um erro")
erro2 := errors.New("Isto é um erro")
fmt.Println(comparaErros(erro1,erro2)) //false
}
```
* https://play.golang.org/p/rFt8xI9Dd3O

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

### fontes
* [Error handling and Go - The go blog](https://blog.golang.org/error-handling-and-go)