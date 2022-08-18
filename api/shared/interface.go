package shared

import (
	"github.com/ElrondNetwork/multi-factor-auth-go-service/api/groups"
	"github.com/ElrondNetwork/multi-factor-auth-go-service/config"
	"github.com/gin-gonic/gin"
)

// GroupHandler defines the actions needed to be performed by an gin API group
type GroupHandler interface {
	UpdateFacade(newFacade FacadeHandler) error
	RegisterRoutes(
		ws *gin.RouterGroup,
		apiConfig config.ApiRoutesConfig,
	)
	IsInterfaceNil() bool
}

// FacadeHandler defines all the methods that a facade should implement
type FacadeHandler interface {
	RestApiInterface() string
	PprofEnabled() bool
	Validate(guardianValidateRequest groups.GuardianValidateRequest) (bool, error)
	RegisterUser(guardianRegisterRequest groups.GuardianRegisterRequest) ([]byte, error)
	IsInterfaceNil() bool
}

// UpgradeableHttpServerHandler defines the actions that an upgradeable http server need to do
type UpgradeableHttpServerHandler interface {
	StartHttpServer() error
	UpdateFacade(facade FacadeHandler) error
	Close() error
	IsInterfaceNil() bool
}