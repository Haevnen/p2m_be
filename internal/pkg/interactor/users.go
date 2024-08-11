package interactor

import (
	"context"
	"errors"
	"time"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/pkg/constants"
	"github.com/Haevnen/p2m_be/pkg/util"
	"gorm.io/gorm"

	"github.com/Haevnen/p2m_be/internal/pkg/interactor/interactorinterface"
)

type UserManagement struct {
	tokenMaker interactorinterface.Maker
}

func NewUserManagement(tokenMaker interactorinterface.Maker) *UserManagement {
	return &UserManagement{
		tokenMaker: tokenMaker,
	}
}

func (ci *UserManagement) GetAllUser(ctx context.Context) (*model.User, error) {
	user, err := dal.Q.User.WithContext(ctx).First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ci *UserManagement) LoginUser(ctx context.Context, body p2m_api.UserLoginBody) (p2m_api.UserLoginResponse, error) {
	u := dal.Q.User

	user, err := u.WithContext(ctx).Where(u.NickName.Eq(body.NickName)).Where(u.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p2m_api.UserLoginResponse{}, apperror.ErrRecordNotFound
		}
		return p2m_api.UserLoginResponse{}, err
	}

	err = util.CheckPassword(body.Password, user.PasswordHashed)
	if err != nil {
		return p2m_api.UserLoginResponse{}, apperror.ErrInvalidPassword
	}

	accessToken, _, err := ci.tokenMaker.CreateToken(user.NickName, user.IsAdmin, constants.AccessTokenDuration)
	if err != nil {
		return p2m_api.UserLoginResponse{}, err
	}

	refreshToken, refreshPayload, err := ci.tokenMaker.CreateToken(user.NickName, user.IsAdmin, constants.RefreshTokenDuration)
	if err != nil {
		return p2m_api.UserLoginResponse{}, err
	}

	err = dal.Q.Session.WithContext(ctx).Create(&model.Session{
		SessionID:    refreshPayload.ID,
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return p2m_api.UserLoginResponse{}, err
	}

	userWoPass := p2m_api.UserWithoutPass{
		ContractType: p2m_api.ContractType(user.ContractType),
		Email:        user.Email,
		IsActive:     user.IsActive,
		IsAdmin:      user.IsAdmin,
		NickName:     user.NickName,
	}

	return p2m_api.UserLoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshPayload.ExpiredAt,
		SessionId:             refreshPayload.ID,
		User:                  userWoPass,
	}, nil
}

func (ci *UserManagement) LogoutUser(ctx context.Context, body p2m_api.RefreshTokenBody) error {
}

func (ci *UserManagement) RefreshToken(ctx context.Context, body p2m_api.RefreshTokenBody) (string, error) {
	refreshPayload, err := ci.tokenMaker.VerifyToken(body.RefreshToken)
	if err != nil {
		return "", apperror.ErrInvalidToken
	}

	u := dal.Q.Session

	session, err := u.WithContext(ctx).Where(u.SessionID.Eq(refreshPayload.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.ErrRecordNotFound
		}
		return "", err
	}

	if session.RefreshToken != body.RefreshToken {
		return "", apperror.ErrInvalidToken
	}

	if session.UserID != refreshPayload.UserID {
		return "", apperror.ErrInvalidToken
	}

	if time.Now().After(session.ExpiredAt) {
		return "", apperror.ErrExpiredRefreshToken
	}

	newAccessToken, _, err := ci.tokenMaker.CreateToken(
		refreshPayload.UserID,
		refreshPayload.IsAdmin,
		constants.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
