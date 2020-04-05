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
	fmt.Println(err1 == err2)
}
