# Tratamento de erros em go antes da versão 1.13

Tratamento de erros em Go é um tópico que gera muita discussão, geralmente reclamações sobre a quantidade o bloco

```golang
if err != nil {
	return err
}
```

aparece no código. Isso geralmente ocorre porque um programador novo em golang acredita que existe um padrão "bala de prata" para tratar erros ou apenas pensa que substituindo todo bloco `try-catch` pelo bloco acima resolverá seus problemas. Entretanto, como Rob Pike diz [neste artigo](https://blog.golang.org/errors-are-values), uma coisa fundamental que estes programadores esquecem é que:

> Erros são valores!
>
> Valores podem ser programados, e como erros são valores, erros podem ser programados.

Apesar de a comparação de um erro com `nil` ser o tratamento mais óbvio possível, existem diversas outras formas de ser tratar um erro, e isso é algo que eu quero mostrar neste artigo, mas antes vamos nos aprofundar um pouco mais para entender como funciona o tipo `error` em go.

### O tipo `error`

O tipo `error` nada mais é do que a interface abaixo

```golang
type error interface {
    Error() string
}
```

e, assim como outros tipos interno de go, é [pré-declarado](https://golang.org/ref/spec#Predeclared_identifiers) no [bloco universal](https://golang.org/ref/spec#Blocks), que engloba todo o código fonte de Go.

### Implemetação do tipo `error`

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

A construção pela função acima, e também pela função `fmt.Errorf()`, acabam se tornando uma dor de cabeça para quem está começando em go pois ao tentar realizar a comparação abaixo o resultado não é bem o esperado.

```golang
package main

import (
	"errors"
	"fmt"
)

func div(dividendo, divisor int) (int, error) {
	if divisor == 0 {
		return 0, errors.New("Erro: divisão por 0")
	}
	return dividendo / divisor, nil
}

func main() {
	_, err1 := div(10, 0)
	_, err2 := div(10, 0)
	fmt.Println(err1 == err2) //false
}
```

Isso ocorre porque, como pode ser visto na implementação da função `errors.New`, essas funções retornam um ponteiro da interface e error, e quando a comparação é feita, ela é feita em cima do endereço dos ponteiros e não em seus valores. Uma forma de realizar essa comparação seria da seguinte forma:

```golang
func main() {
	_, err1 := div(10, 0)
	_, err2 := div(10, 0)
	fmt.Println(err1.Error() == err2.Error()) //true
}
```

No exemplo acima a comparação está sendo feita entre as mensagens dos erros e não em cima dos seus endereços.

Dito isso, como poderiamos tratar melhor os nossos erros em go?

<!--

NA: isso é válido?

Primeiramente entendendo que, isso não apenas em go mas em qualquer linguagem, erros não são efeitos colaterais ou algo inesperado que ocorreu no seu programa, erros são parte da sua aplicação e devem ser tratados como tais. Generalizar um erro comparando ele com `nil` apenas pode ser uma grande dor de cabeça com o crescimento da sua aplicação.

Lógico que quanto mais perto da fronteira do `core` da sua aplicação, mais difícil é tratar um erro pois não sabemos quais erros são gerados em um package de terceiro. Mesmo assim, todo erro que atravessar esta fronteira adentrando no `core` da sua aplicação devem ser tratados como parte da regra de negócio da mesma. -->

### Sentinelas

Para evitar o problema de comparaç~]ao de erros dito anteriormente, uma solução é a utilização de sentinelas que são apenas a declaração dos erros da aplicação em variáveis ou constantes. Tendo feito isso, a comparação não precisa mais ser entre o texto do `error`, através a função `Error()`, e pode ser feita através do endereço dele, pois agora o erro irá se encontrar no mesmo bloco de memória sempre.

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
		//trata o erro aqui
	}else{
		fmt.Println(str)
	}
}
```

Como pode ser visto no exemplo acima, a comparação do erro retornado não é mais com o texto do erro ou com `nil`, mas sim com a variável `ErroNomeVazio`.

### Tipos de erros Customizados

Caso seja necessário tratar os erros de forma específica, uma alteranitiva é a criação de erros personalidados implementando a interface `error` apresentada anteriormente. As vantagens desta abordagem é a criação de mensagens personalidadas e o tratamento  do erro através da comparação do seu tipo.

```golang
// O exemplo completo você pode ver em https://github.com/luizvnasc/go-error-handling/tree/master/exemplos/3_erros_customizados
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
```

No exemplo acima não precisamos da utilização de sentinelas para verificar o erro que foi retornado para realizarmos seu devido tratamento. Em vez disso comparamos o seu tipo com os tipos de erro da aplicação e, dependendo do erro retornado, é feita a impressão da mensagem de erro.

### Stack trace e 3th parties

Até a versão o go 1.13, o pacote padrão `errors` nãop possuía nenhuma implementação de _stacktraces_. Por este motivo foram criados algums pacotes de terceiros para solucionar este problema como [palantir/stacktrace](https://github.com/palantir/stacktrace), [go-erros/errors](https://github.com/go-errors/errors) e [pkg/errors](https://github.com/pkg/errors), este último o mais popular entre eles com aproximadamente 5700 stars e 419 forks até a data de publicação deste artigo.

Modificando um pouco nosso exemplo anterior para utilizar a biblioteca `pkg/errors` temos o seguinte código:

```golang
// O exemplo completo você pode visualizar em https://github.com/luizvnasc/go-error-handling/tree/master/exemplos/4_stacktrace
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

// Os erros customizados a função que gera a mensagem de bem vindo não foram alterados.

// DigaBemVindo imprime uma mensagem de bem vindo para um participante do meetup.
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
```

Como pode ser visto no exemplo, utilizamos duas funções da bibliteca `pkg/errors`. A função `Wrap(err error, message string)` embrulha um erro em um novo erro para que seja criada a _stacktrace_, já a função `Cause(err error) error` percorre a pilha de forma recursiva até chegar a causa do problema, ou seja, aquele erro que embrulha nenhum outro erro.

### A vida após o go 1.13

### Referências bibliográficas

-   [Error handling and Go - The go blog](https://blog.golang.org/error-handling-and-go)
-   [Nerdgirlz #30 - Go Go Go!](https://www.youtube.com/watch?v=ZAmESdN5alo)
-   [Errors are values](https://blog.golang.org/errors-are-values)
-   [Error Handling in Go that Every Beginner should Know](https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76)
