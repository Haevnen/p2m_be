package interactor

import (
	"context"
	"errors"

	"gorm.io/gorm"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

type UserManagement struct{}

func NewUserManagement() *UserManagement {
	return &UserManagement{}
}

func (ci *UserManagement) GetAllUser(ctx context.Context, includeDeActive *bool) ([]*p2mapi.User, error) {
	u := dal.Q.User
	var users []*model.User
	var err error

	if includeDeActive == nil || *includeDeActive == false {
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

	if user == nil {
		return apperror.ErrRecordNotFound
	}

	info, err := u.WithContext(ctx).Where(u.ID.Eq(user.ID)).Update(u.IsActive, false)
	if err != nil {
		return err
	}

	if info.RowsAffected == 0 {
		return apperror.ErrInternalServer
	}

	return nil
}
