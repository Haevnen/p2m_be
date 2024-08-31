package model

import (
	"github.com/go-openapi/swag"

	p2mapi "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

func (c *Client) ToClient(client p2mapi.ClientBody) error {
	c.ClientID = client.ClientId
	c.EditingStyle = swag.StringValue(client.EditingStyle)
	c.Others = swag.StringValue(client.Others)
	c.Requirements = swag.StringValue(client.Requirements)
	return nil
}

func (c *Client) FromClient() *p2mapi.ClientResponse {

	id := c.ID
	return &p2mapi.ClientResponse{
		ClientId:     c.ClientID,
		EditingStyle: &c.EditingStyle,
		Id:           id,
		IsActive:     c.IsActive,
		Others:       &c.Others,
		Requirements: &c.Requirements,
	}
}
