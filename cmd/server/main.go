package main

import (
	"fmt"

	"github.com/waanvieira/api-users/internal/entity"
)

func main() {
	user, err := entity.NewUser("user test", "email.br", "123")
	if err != nil {
		panic("erro ao criar user")
	}
	fmt.Println(user)
}
