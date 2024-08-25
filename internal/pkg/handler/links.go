package handler

import (
	"net/http"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/apperror"
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
func (h linkHandler) InternalUpdateLink(c *gin.Context, linkId int64, params p2m_api.InternalUpdateLinkParams) {
	var updateLinkBody p2m_api.LinkBody
	if err := bindRequestBody(c, &updateLinkBody); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	err := h.linkManagementInteractor.UpdateLink(c, linkId, updateLinkBody)
	if err != nil {
		SendError(c, "update link error", err)
		return
	}

	c.JSON(http.StatusOK, "update link successfully")
}

// Delete link by id
// (DELETE /links/{link_id})
func (h linkHandler) InternalRemoveLink(c *gin.Context, linkId int64, params p2m_api.InternalRemoveLinkParams) {
	err := h.linkManagementInteractor.RemoveLink(c, linkId)
	if err != nil {
		SendError(c, "remove link error", err)
		return
	}

	c.JSON(http.StatusOK, "remove link successfully")
}

// Get ticket links
// (GET /links/{ticket_id})
func (h linkHandler) InternalGetAllLinks(c *gin.Context, ticketId int64) {
	links, err := h.linkManagementInteractor.GetAllLink(c, ticketId)
	if err != nil {
		SendError(c, "get all link error", err)
		return
	}

	c.JSON(http.StatusOK, links)
}

// Create new link
// (POST /links)
func (h linkHandler) InternalCreateLink(c *gin.Context, params p2m_api.InternalCreateLinkParams) {
	var createLinkBody p2m_api.LinkBody
	if err := bindRequestBody(c, &createLinkBody); err != nil {
		SendError(c, err.Error(), apperror.ErrInvalidRequestInput)
		return
	}

	resp, err := h.linkManagementInteractor.CreateLink(c, createLinkBody)
	if err != nil {
		SendError(c, "create link error", err)
		return
	}

	c.JSON(http.StatusOK, resp)

}
