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

## O tipo `error`

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

## Sentinelas

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

## Tipos de erros Customizados

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

## Stack trace e 3th parties

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

## A vida após o go 1.13

O lançamento do go 1.13 trouxe algumas funcionalidades para os pacotes padrões `errors` e `fmt` para tratar erros que embrulham outros erros. Dentre elas a convenção de que um erro que embrulha outro deve implementar a função Unwrap que retorna o erro embrulhado, conforme o exemplo abaixo:

```golang
type AppError struct {
    msg string
    Err error
}

func (a *AppError) Unwrap() error { return a.Err }

```

Como pode ser visto, diferentemente do pacote `pkg/errors`, nesta versão de go não foi implementada nenhuma função `Cause` que percorreria toda a _stacktrace_ até a raiz do erro.

### Embrulhando erros com `%w`

Como dito anteriormente, o pacote `fmt` ganhou uma nova funcionalidade para o embrulho de erros. Agora a função `fmt.Errorf` suporta a expressão `%w` que é responsável por embrulhar o erro informado dentro do erro criado. Modificando um pouco nosso exemplo de mensagem de boas vindas teremos:

```golang
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

```
Veja que agora a função `DigaBemVindo` embrulha o erro retornado pela função `BemVindo` em um novo erro utilizando a função `fmt.Errorf`. Também criamos sentinelas em vez de utilizar implementações customizadas. Isso foi apenas para podr mostrar outra funcionalidade implementada nesta versão do go que é a comparação de erros com a função `Is` da biblioteca padrão `errors`. Ela basicamente verifica se um erro, ou os erros contidos nele, são o mesmo que o sentinela.

Caso queira fazer uma comparação de tipo com dioferentes implementações de erro utilizamos a função `As`,também da biblioteca padrão `errors`. Vejamos como fica a implementação acima com utilizando a função `errors.As`.


```golang
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
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

var (
	errNomeVazio            StringVaziaError
	errNomeNumerico         StringNumericaError
	errNomeCaracterEspecial StringComCaracteresEspeciaisError
)

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
		if errors.As(err, &errNomeVazio) {
			fmt.Fprintln(os.Stdout, "Não aceitamos pessoas anônimas!")
		}
		if errors.As(err, &errNomeNumerico) {
			fmt.Fprintln(os.Stdout, "Te entendo, somos todos apenas números.")
		}
		if errors.As(err, &errNomeCaracterEspecial) {
			fmt.Fprintln(os.Stdout, "Você ainda usa hotmail?")
		}
	}
}
```

Veja que, mesmo utilizando tipos diferentes de erros, a função `As` precisa de sentinelas para que haja a comparação, entretanto desta vez os sentinelas são tipos.


## `defer`, `panic` e `recover`

Outra forma de tratar erros em go é a utilização da declaração `defer` e das funções panic e recover.

a declaração `defer` empilha uma função na pilha de execução para ser executada ao fim da função na qual ela foi chamada. Essa declaração é extremamente quando se precisa fazer alguma limpeza no final da execução de alguma função.

```golang
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}
```
A função acima faz a copia de um arquivo. Nela é possível ver que a declaração `defer` é utilizada para fechar os arquivos de origem e destino. No caso do arquivo de origem, a declaração `defer`garante que ele seja fechado mesmo que ocorra um erro durante a cópia. Outra vantagem da declaração `defer` é na organização do código, pedir para fechar um arquivo que foi recém aberto é mais legível que lembrar de fechar ele no fim da função.

A função `panic` para todo o fluxo de execução do go e entra em "pânico". Isso quer dizer que a se uma função **F** entra em pânico, ela irá para o restante da sua execução, a função declarada com `defer` será executada normalmente, e irá retornar para quem à chamou. Isso ocorrerá até que todas as funções chamadas sejam retornadas e o programa quebre. 

A função `recover` recupera o programa em pânico e retorna a execução normal a partir dela. Ela deve ser utilizada dentro de uma declaração `defer`, pois gantante que será executada em um programa em pânico.

```golang
package main

import "fmt"

func main() {
    f()
    fmt.Println("Returned normally from f.")
}

func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    fmt.Println("Calling g.")
    g(0)
    fmt.Println("Returned normally from g.")
}

func g(i int) {
    if i > 3 {
        fmt.Println("Panicking!")
        panic(fmt.Sprintf("%v", i))
    }
    defer fmt.Println("Defer in g", i)
    fmt.Println("Printing in g", i)
    g(i + 1)
}
```

No exemplo acima, a função `f()` chama a função `g(0)` que é executada recursivamente até o valor de i ser igual a 3. Após isso ela entra em pânico e retorna até ser recuperada na função `f()` que pega o valor recuperado e imprime no terminal.

## Conclusão

 Considerar que o tratamento de erro em go é simplório é uma interpretação erronea pois, pesar da linguagem go não possuir exceções como outras linguagem, ela possui diversas formas de tratamento de erros que da um leque de possibilidades para o desenvolvedor. Analisar como os erros da sua aplicação devem ser tratados é um desafio e tanto e eu espero que este artigo tenha ajudado quem busca a melhor maneira de ultrapassá-lo.

### Fontes

-   [Error handling and Go - The go blog](https://blog.golang.org/error-handling-and-go)
-   [Nerdgirlz #30 - Go Go Go!](https://www.youtube.com/watch?v=ZAmESdN5alo)
-   [Errors are values](https://blog.golang.org/errors-are-values)
-   [Error Handling in Go that Every Beginner should Know](https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76)
-   [Working with Errors in Go 1.13](https://blog.golang.org/go1.13-errors)
-   [defer, panic and recover](https://blog.golang.org/defer-panic-and-recover)
