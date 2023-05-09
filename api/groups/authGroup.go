package groups

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/dragos-rebegea/evaluare-tool/api/shared"
	"github.com/dragos-rebegea/evaluare-tool/authentication"
	"github.com/dragos-rebegea/evaluare-tool/core"
	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/api/errors"
	elrondApiShared "github.com/multiversx/mx-chain-go/api/shared"
)

const (
	registerPath = "/register"
)

type authGroup struct {
	*baseGroup
	facade               shared.FacadeHandler
	mutFacade            sync.RWMutex
	database             *core.DatabaseHandler
	authenticationNeeded bool
}

// NewAuthGroup returns a new instance of evaluationGroup
func NewAuthGroup(facade shared.FacadeHandler, dbHandler *core.DatabaseHandler) (*authGroup, error) {
	if check.IfNil(facade) {
		return nil, fmt.Errorf("%w for auth group", errors.ErrNilFacadeHandler)
	}
	if check.IfNil(dbHandler) {
		return nil, fmt.Errorf("%w for auth group", ErrNilDatabaseHandler)
	}
	ag := &authGroup{
		facade:               facade,
		baseGroup:            &baseGroup{},
		database:             dbHandler,
		authenticationNeeded: false,
	}

	endpoints := []*elrondApiShared.EndpointHandlerData{
		{
			Path:    tokenPath,
			Method:  http.MethodPost,
			Handler: ag.generateToken,
		},
		{
			Path:    registerPath,
			Method:  http.MethodPost,
			Handler: ag.registerAdmin,
		},
	}
	ag.endpoints = endpoints

	return ag, nil
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ag *authGroup) registerAdmin(context *gin.Context) {
	var admin authentication.Profesor
	if err := context.ShouldBindJSON(&admin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	admin.IsAdmin = true
	err := ag.database.CreateProfesor(&admin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"userId":   admin.ID,
		"email":    admin.Email,
		"username": admin.Username,
		"password": admin.Password,
	})

}

func (ag *authGroup) generateToken(context *gin.Context) {
	var request TokenRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user, err := ag.database.GetProfesorByEmail(request.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := authentication.GenerateJWT(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// UpdateFacade will update the facade
func (ag *authGroup) UpdateFacade(newFacade shared.FacadeHandler) error {
	if check.IfNil(newFacade) {
		return errors.ErrNilFacadeHandler
	}

	ag.mutFacade.Lock()
	ag.facade = newFacade
	ag.mutFacade.Unlock()

	return nil
}

// IsAuthenticationNeeded will return true if the group requires authentication
func (ag *authGroup) IsAuthenticationNeeded() bool {
	return ag.authenticationNeeded
}

// IsInterfaceNil returns true if there is no value under the interface
func (ag *authGroup) IsInterfaceNil() bool {
	return ag == nil
}
