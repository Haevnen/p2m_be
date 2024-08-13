package interactor

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/Haevnen/p2m_be/pkg/constants"
	"github.com/Haevnen/p2m_be/pkg/util"
)

type UserManagement struct {
	tokenMaker interactorinterface.Maker
}

func NewUserManagement(tokenMaker interactorinterface.Maker) *UserManagement {
	return &UserManagement{
		tokenMaker: tokenMaker,
	}
}

func (ci *UserManagement) GetAllUser(ctx context.Context, includeDeActive *bool) ([]*p2mapi.User, error) {
	u := dal.Q.User
	var users []*model.User
	var err error

	if includeDeActive == nil || !(*includeDeActive) {
		users, err = u.WithContext(ctx).Where(u.IsActive).Find()
	} else {
		users, err = u.WithContext(ctx).Find()
	}

	if err != nil {
		return nil, err
	}

	res := make([]*p2mapi.User, 0, len(users))
	for _, user := range users {
		res = append(res, user.FromUser())
	}

	return res, nil
}

func (ci *UserManagement) LoginUser(ctx context.Context, body p2mapi.UserLoginBody) (p2mapi.UserLoginResponse, error) {
	u := dal.Q.User

	user, err := u.WithContext(ctx).Where(u.NickName.Eq(body.NickName)).Where(u.IsActive.Is(true)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p2mapi.UserLoginResponse{}, apperror.ErrRecordNotFound
		}
		return p2mapi.UserLoginResponse{}, err
	}

	err = util.CheckPassword(body.Password, user.PasswordHashed)
	if err != nil {
		return p2mapi.UserLoginResponse{}, apperror.ErrInvalidPassword
	}

	accessToken, _, err := ci.tokenMaker.CreateToken(user.UserID, user.IsAdmin, constants.AccessTokenDuration)
	if err != nil {
		return p2mapi.UserLoginResponse{}, err
	}

	refreshToken, refreshPayload, err := ci.tokenMaker.CreateToken(user.UserID, user.IsAdmin, constants.RefreshTokenDuration)
	if err != nil {
		return p2mapi.UserLoginResponse{}, err
	}

	err = dal.Q.Session.WithContext(ctx).Create(&model.Session{
		SessionID:    refreshPayload.ID.String(),
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return p2mapi.UserLoginResponse{}, err
	}

	userWoPass := p2mapi.UserWithoutPass{
		ContractType: p2mapi.ContractType(user.ContractType),
		Email:        user.Email,
		IsActive:     user.IsActive,
		IsAdmin:      user.IsAdmin,
		NickName:     user.NickName,
	}

	return p2mapi.UserLoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshPayload.ExpiredAt.Format(constants.DateTimeFormat),
		SessionId:             refreshPayload.ID.String(),
		User:                  userWoPass,
	}, nil
}

func (ci *UserManagement) LogoutUser(ctx context.Context, body p2mapi.RefreshTokenBody) error {
	refreshPayload, err := ci.tokenMaker.VerifyToken(body.RefreshToken)
	if err != nil {
		return apperror.ErrInvalidRefreshToken
	}

	u := dal.Q.Session

	session, err := u.WithContext(ctx).Where(u.SessionID.Eq(refreshPayload.ID.String())).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}
		return err
	}

	if session.RefreshToken != body.RefreshToken {
		return apperror.ErrInvalidRefreshToken
	}

	if session.UserID != refreshPayload.UserID {
		return apperror.ErrInvalidRefreshToken
	}

	// remove all existing session relate to UserID
	_, err = u.WithContext(ctx).Where(u.UserID.Eq(refreshPayload.UserID)).Delete(&model.Session{})
	return err
}

func (ci *UserManagement) RefreshToken(ctx context.Context, body p2mapi.RefreshTokenBody) (string, error) {
	refreshPayload, err := ci.tokenMaker.VerifyToken(body.RefreshToken)
	if err != nil {
		return "", apperror.ErrInvalidRefreshToken
	}

	u := dal.Q.Session

	session, err := u.WithContext(ctx).Where(u.SessionID.Eq(refreshPayload.ID.String())).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.ErrRecordNotFound
		}
		return "", err
	}

	if session.RefreshToken != body.RefreshToken {
		return "", apperror.ErrInvalidRefreshToken
	}

	if session.UserID != refreshPayload.UserID {
		return "", apperror.ErrInvalidRefreshToken
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

func (ci *UserManagement) CreateUser(ctx context.Context, user p2mapi.User) (*p2mapi.User, error) {
	u := dal.Q.User

	var userDb model.User
	err := userDb.ToUser(user)
	if err != nil {
		return nil, err
	}

	// validate user
	// if user has nick_name and is_active is true, return error
	userCheck, err := u.WithContext(ctx).Where(u.NickName.Eq(userDb.NickName), u.IsActive).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if userCheck != nil {
		return nil, apperror.ErrUserHasNicknameExists
	}

	// if user has email and is_active is true, return error
	userCheck, err = u.WithContext(ctx).Where(u.Email.Eq(userDb.Email), u.IsActive).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if userCheck != nil {
		return nil, apperror.ErrUserHasEmailExists
	}

	err = u.WithContext(ctx).Save(&userDb)
	if err != nil {
		return nil, err
	}

	return userDb.FromUser(), nil
}

func (ci *UserManagement) RemoveUser(ctx context.Context, nickName string) error {
	u := dal.Q.User

	user, err := u.WithContext(ctx).Where(u.NickName.Eq(nickName), u.IsActive).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrRecordNotFound
		}

		return err
	}

	info, err := u.WithContext(ctx).Where(u.ID.Eq(user.ID)).Update(u.IsActive, false)
	if err != nil {
		return err
	}

	if info.RowsAffected != 1 {
		return apperror.ErrInternalServer
	}

	// Remove all existing session relate to deleted user
	s := dal.Q.Session
	_, err = s.WithContext(ctx).Where(s.UserID.Eq(user.UserID)).Delete(&model.Session{})
	return err
}
