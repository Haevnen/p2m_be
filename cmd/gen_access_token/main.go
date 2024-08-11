package main

import (
	"fmt"
	"time"

	"github.com/Haevnen/p2m_be/internal/pkg/interactor"
)

const (
	secretKey = "12345678901234567890123456789012"
)

func main() {
	pasetoMaker, err := interactor.NewPasetoMaker(secretKey)
	if err != nil {
		panic(err)
	}

	token, _, err := pasetoMaker.CreateToken("admin", true, 24*time.Hour)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
