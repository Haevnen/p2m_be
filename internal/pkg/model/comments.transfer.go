package model

import (
	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

func (c *Comment) ToComment(comment p2mapi.CreateCommentBody) {
	c.Comment = comment.Comment
	c.TicketID = comment.TicketId
}

func (c *Comment) FromComment() *p2mapi.CommentResponse {
	return &p2mapi.CommentResponse{
		Comment:  c.Comment,
		TicketId: c.TicketID,
		Id:       c.ID,
	}
}
