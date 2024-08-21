package handler

import (
	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
	"github.com/gin-gonic/gin"
)

type linkHandler struct {
	linkManagementInteractor interactorinterface.LinkManagementInterface
}

func newLinkHandler(registry *registry.Registry) linkHandler {
	return linkHandler{
		linkManagementInteractor: registry.LinkManagementInterface(),
	}
}

// Update link by id
// (PUT /links/update/{link_id})
func (h linkHandler) InternalUpdateLink(c *gin.Context, linkId int, params p2m_api.InternalUpdateLinkParams) {
}

// Delete link by id
// (DELETE /links/{link_id})
func (h linkHandler) InternalRemoveLink(c *gin.Context, linkId int, params p2m_api.InternalRemoveLinkParams) {
}

// Get ticket links
// (GET /links/{ticket_id})
func (h linkHandler) InternalGetAllLinks(c *gin.Context, ticketId int) {
}

// Create new link
// (POST /links)
func (h linkHandler) InternalCreateLink(c *gin.Context, params p2m_api.InternalCreateLinkParams) {}
