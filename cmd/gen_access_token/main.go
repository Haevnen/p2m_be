package main

import (
	"fmt"
	"time"

	"github.com/Haevnen/p2m_be/internal/pkg/interactor"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/util"
	"github.com/google/uuid"
)

const (
	secretKey = "ffad96f69973f0cecafe52ab794d5556"
)

func decodeToken(maker interactorinterface.Maker, token string) {
	payload, _ := maker.VerifyToken(token)
	fmt.Println("Expired at: ", payload.ExpiredAt)
	fmt.Println("Issued at: ", payload.IssuedAt)
	fmt.Println("ID: ", payload.ID)
	fmt.Println("User ID: ", payload.UserID)
}

func main() {
	uuid, _ := uuid.NewRandom()
	fmt.Println("UUID: ", uuid.String())

	pasetoMaker, err := interactor.NewPasetoMaker(secretKey)
	if err != nil {
		panic(err)
	}

	token, _, err := pasetoMaker.CreateToken("admin", true, 24*time.Hour)
	if err != nil {
		panic(err)
	}

	fmt.Println("Token: ", token)

	decodeToken(pasetoMaker, token)

	password := "Ck27082021$"

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		panic(err)
	}
	fmt.Println(hashedPassword)
}
