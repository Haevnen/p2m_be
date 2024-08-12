package model

import (
	"github.com/go-openapi/swag"
	"github.com/google/uuid"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/pkg/logger"
	"github.com/Haevnen/p2m_be/pkg/password"
)

func (u *User) ToUser(user p2mapi.User) error {
	// has password
	passwordHash, err := password.HashPassword(swag.StringValue(user.Password))
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return err
	}

	u.NickName = user.NickName
	u.Email = user.Email
	u.PasswordHashed = passwordHash
	u.ContractType = string(user.ContractType)
	u.IsActive = user.IsActive
	u.IsAdmin = user.IsAdmin

	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	u.UserID = userID.String()

	return nil
}

func (u *User) FromUser() *p2mapi.User {

	return &p2mapi.User{
		NickName:     u.NickName,
		Email:        u.Email,
		ContractType: p2mapi.ContractType(u.ContractType),
		IsActive:     u.IsActive,
		IsAdmin:      u.IsAdmin,
	}
}
