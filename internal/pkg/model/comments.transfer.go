package model

import (
	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

type CommentWithName struct {
	Comment
	NickName string `gorm:"column:nick_name"`
}

func (c *Comment) ToComment(comment p2mapi.CreateCommentBody) {
	c.Comment = comment.Comment
	c.TicketID = comment.TicketId
}

func (c *Comment) FromComment(nickName string) *CommentWithName {
	return &CommentWithName{
		Comment:  *c,
		NickName: nickName,
	}
}

func (cn *CommentWithName) FromCommentWithName() *p2mapi.CommentResponse {
	return &p2mapi.CommentResponse{
		Comment:  cn.Comment.Comment,
		TicketId: cn.TicketID,
		Id:       cn.ID,
		NickName: &cn.NickName,
	}
}
