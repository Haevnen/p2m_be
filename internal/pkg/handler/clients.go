package handler

import (
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
}

// Register new client
// (POST /clients/register)
func (ci clientHandler) InternalRegisterClient(c *gin.Context, params p2m_api.InternalRegisterClientParams) {
}

// Delete client by id
// (DELETE /clients/{id})
func (ci clientHandler) InternalRemoveClient(c *gin.Context, id int, params p2m_api.InternalRemoveClientParams) {
}
