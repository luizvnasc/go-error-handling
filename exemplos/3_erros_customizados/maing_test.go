package main

import (
	"testing"
)

func TestOla(t *testing.T) {
	t.Run("Passando uma string numérica como parâmetro", func(t *testing.T) {
		msgBoasVindas, err := BemVindo("01")
		if err == nil {
			t.Error("Erro ao criar mensagem de boas vindas: Esperado um erro, obitido nil")
		}
		if msgBoasVindas != "" {
			t.Errorf("Erro ao criar mensagem de boas vindas: Esperado uma string vazia, obitido %q", msgBoasVindas)
		}
	})
	t.Run("Passando uma string float como parâmetro", func(t *testing.T) {
		msgBoasVindas, err := BemVindo("01.10")
		if err == nil {
			t.Error("Erro ao criar mensagem de boas vindas: Esperado um erro, obitido nil")
		}
		if msgBoasVindas != "" {
			t.Errorf("Erro ao criar mensagem de boas vindas: Esperado uma string vazia, obitido %q", msgBoasVindas)
		}
	})
	t.Run("Passando uma string com caracteres especiais como parâmetro", func(t *testing.T) {
		msgBoasVindas, err := BemVindo("#Golang@CWB")
		if err == nil {
			t.Error("Erro ao criar mensagem de boas vindas: Esperado um erro, obitido nil")
		}
		if msgBoasVindas != "" {
			t.Errorf("Erro ao criar mensagem de boas vindas: Esperado uma string vazia, obitido %q", msgBoasVindas)
		}
	})
	t.Run("Passando uma string vazia como parâmetro", func(t *testing.T) {
		msgBoasVindas, err := BemVindo("")
		if err == nil {
			t.Error("Erro ao criar mensagem de boas vindas: Esperado um erro, obitido nil")
		}
		if msgBoasVindas != "" {
			t.Errorf("Erro ao criar mensagem de boas vindas: Esperado uma string vazia, obitido %q", msgBoasVindas)
		}
	})
	t.Run("Passando uma string válida", func(t *testing.T) {
		msgBoasVindas, err := BemVindo("folks")
		if err != nil {
			t.Errorf("Erro ao criar mensagem de boas vindas: Esperado nil, obitido %v", err)
		}
		if msgBoasVindas == "" {
			t.Error("Erro ao criar mensagem de boas vindas: Esperado uma mensagem, obitido nil")
		}
	})
}

func BenchmarkDigaBemVindoCustom(b *testing.B) {
	nomes := []string{"1", "John Doe", "Golang@CWB", ""}
	for i := 0; i < b.N; i++ {
		nome := nomes[i%len(nomes)]
		DigaBemVindoCustom(&spyWriter{}, nome)
	}
}

func BenchmarkDigaBemVindo(b *testing.B) {
	nomes := []string{"1", "John Doe", "Golang@CWB", ""}
	for i := 0; i < b.N; i++ {
		nome := nomes[i%len(nomes)]
		DigaBemVindo(&spyWriter{}, nome)
	}
}

type spyWriter struct {
	counter int
}

func (spy *spyWriter) Write(p []byte) (n int, err error) {
	spy.counter++
	return spy.counter, nil
}
