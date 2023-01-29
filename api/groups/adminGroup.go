package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/api/errors"
	elrondApiShared "github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/dragos-rebegea/evaluare-tool/api/shared"
	"github.com/dragos-rebegea/evaluare-tool/authentication"
	"github.com/dragos-rebegea/evaluare-tool/core"
	"github.com/gin-gonic/gin"
)

type adminGroup struct {
	*baseGroup
	facade               shared.FacadeHandler
	mutFacade            sync.RWMutex
	database             *core.DatabaseHandler
	authenticationNeeded bool
}

// NewAdminGroup returns a new instance of adminGroup
func NewAdminGroup(facade shared.FacadeHandler, dbHandler *core.DatabaseHandler) (*adminGroup, error) {
	if check.IfNil(facade) {
		return nil, fmt.Errorf("%w for admin group", errors.ErrNilFacadeHandler)
	}
	if check.IfNil(dbHandler) {
		return nil, fmt.Errorf("%w for admin group", ErrNilDatabaseHandler)
	}

	ag := &adminGroup{
		facade:               facade,
		baseGroup:            &baseGroup{},
		database:             dbHandler,
		authenticationNeeded: true,
	}

	endpoints := []*elrondApiShared.EndpointHandlerData{
		{
			Path:    "/ping",
			Method:  http.MethodGet,
			Handler: ag.ping,
		},
		{
			Path:    "/createClass",
			Method:  http.MethodPost,
			Handler: ag.createClass,
		},
	}
	ag.endpoints = endpoints

	return ag, nil
}

func (ag *adminGroup) ping(c *gin.Context) {
	if !ag.checkIfAdmin(c) {
		return
	}

	c.JSON(
		http.StatusOK,
		elrondApiShared.GenericAPIResponse{
			Data:  "pong",
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)
}

func (ag *adminGroup) checkIfAdmin(c *gin.Context) bool {
	email := c.GetString(authentication.EmailKey)

	isAdmin, err := ag.database.IsAdmin(email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return false
	}

	if !isAdmin {
		c.JSON(
			http.StatusForbidden,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: "you are not an admin",
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return false
	}

	return true
}

func (ag *adminGroup) createClass(c *gin.Context) {
	if !ag.checkIfAdmin(c) {
		return
	}

	var class core.Class
	err := json.NewDecoder(c.Request.Body).Decode(&class)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return
	}

	err = ag.database.CreateClass(&class)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return
	}
}

// UpdateFacade will update the facade
func (ag *adminGroup) UpdateFacade(newFacade shared.FacadeHandler) error {
	if check.IfNil(newFacade) {
		return errors.ErrNilFacadeHandler
	}

	ag.mutFacade.Lock()
	ag.facade = newFacade
	ag.mutFacade.Unlock()

	return nil
}

// IsAuthenticationNeeded will return true if the group requires authentication
func (ag *adminGroup) IsAuthenticationNeeded() bool {
	return ag.authenticationNeeded
}

// IsInterfaceNil returns true if there is no value under the interface
func (ag *adminGroup) IsInterfaceNil() bool {
	return ag == nil
}
