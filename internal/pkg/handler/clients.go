package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	p2m_api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
	"github.com/Haevnen/p2m_be/internal/pkg/registry"
	"github.com/Haevnen/p2m_be/internal/pkg/registry/interactorinterface"
)

type clientHandler struct {
	clientManagementInteractor interactorinterface.ClientManagementInterface
}

func newClientHandler(registry *registry.Registry) clientHandler {
	return clientHandler{
		clientManagementInteractor: registry.ClientManagementInteractor(),
	}
}

// Get all clients
// (GET /clients)
func (ci clientHandler) InternalGetAllClients(c *gin.Context, params p2m_api.InternalGetAllClientsParams) {
	clients, err := ci.clientManagementInteractor.GetAllClient(c, params.IncludingDeactivates)
	if err != nil {
		SendError(c, "get all clients error", err)
		return
	}

	c.JSON(http.StatusOK, clients)
}

// Register new client
// (POST /clients/register)
func (ci clientHandler) InternalRegisterClient(c *gin.Context, params p2m_api.InternalRegisterClientParams) {
	var client p2m_api.ClientBody
	if err := bindRequestBody(c, &client); err != nil {
		SendError(c, "bind error", err)
		return
	}

	res, err := ci.clientManagementInteractor.CreateClient(c, client)
	if err != nil {
		SendError(c, "create client error", err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete client by id
// (DELETE /clients/{id})
func (ci clientHandler) InternalRemoveClient(c *gin.Context, clientId string, params p2m_api.InternalRemoveClientParams) {
	err := ci.clientManagementInteractor.RemoveClient(c, clientId)
	if err != nil {
		SendError(c, "remove client error", err)
		return
	}

	c.JSON(http.StatusOK, "remove client successfully")
}
