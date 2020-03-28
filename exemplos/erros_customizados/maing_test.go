package main

import "testing"

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
