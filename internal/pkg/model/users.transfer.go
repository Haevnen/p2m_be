package model

import (
	apiModel "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

func (u *User) ToDB(user apiModel.User) *User {

	u.NickName = user.NickName
	u.Email = user.Email

	u.ContractType = string(user.ContractType)

	return u
}

func (u *User) ToAPI() apiModel.User {
	return apiModel.User{
		NickName: u.NickName,
		Email:    u.Email,
	}
}
