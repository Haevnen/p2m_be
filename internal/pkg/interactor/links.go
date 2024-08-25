package interactor

import (
	"context"
	"errors"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
	"github.com/Haevnen/p2m_be/internal/pkg/dal"
	"github.com/Haevnen/p2m_be/internal/pkg/model"
	"gorm.io/gorm"
)

type LinkManagement struct {
}

func NewLinkManagement() *LinkManagement {
	return &LinkManagement{}
}

func (li *LinkManagement) CreateLink(ctx context.Context, link p2mapi.LinkBody) (*p2mapi.LinkResponse, error) {
	var linkModel model.Link
	linkModel.ToLink(link)

	l := dal.Q.Link
	// gorm will automatically fetch the the value of PK and populate to linkModel.ID
	err := l.WithContext(ctx).Create(&linkModel)
	if err != nil {
		return nil, err
	}

	return linkModel.FromLink(), nil
}
func (li *LinkManagement) GetAllLink(ctx context.Context, ticketID int64) ([]*p2mapi.LinkResponse, error) {
	l := dal.Q.Link
	links, err := l.WithContext(ctx).Where(l.TicketID.Eq(ticketID)).Find()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	var linkRes []*p2mapi.LinkResponse
	for _, v := range links {
		linkRes = append(linkRes, v.FromLink())
	}
	return linkRes, nil

}
func (li *LinkManagement) RemoveLink(ctx context.Context, linkID int64) error {
	l := dal.Q.Link
	_, err := l.WithContext(ctx).Where(l.ID.Eq(linkID)).Delete(&model.Link{})
	return err
}

func (li *LinkManagement) UpdateLink(ctx context.Context, linkID int64, body p2mapi.LinkBody) error {
	l := dal.Q.Link

	_, err := l.WithContext(ctx).Where(l.ID.Eq(int64(linkID))).UpdateColumn(l.Link, body.Url)
	return err
}
