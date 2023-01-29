package groups

import (
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

const tokenPath = "/token"

type evaluationGroup struct {
	*baseGroup
	facade               shared.FacadeHandler
	mutFacade            sync.RWMutex
	database             *core.DatabaseHandler
	authenticationNeeded bool
}

// NewEvaluationGroup returns a new instance of evaluationGroup
func NewEvaluationGroup(facade shared.FacadeHandler, dbHandler *core.DatabaseHandler) (*evaluationGroup, error) {
	if check.IfNil(facade) {
		return nil, fmt.Errorf("%w for evaluation group", errors.ErrNilFacadeHandler)
	}
	if check.IfNil(dbHandler) {
		return nil, fmt.Errorf("%w for evaluation group", ErrNilDatabaseHandler)
	}
	eg := &evaluationGroup{
		facade:               facade,
		baseGroup:            &baseGroup{},
		database:             dbHandler,
		authenticationNeeded: true,
	}

	endpoints := []*elrondApiShared.EndpointHandlerData{
		{
			Path:    "/getStudentsByClass/:class",
			Method:  http.MethodGet,
			Handler: eg.getStudentsByClass,
		},
	}
	eg.endpoints = endpoints

	return eg, nil
}

// sendTransaction returns will send the transaction signed by the guardian if the verification passed
func (eg *evaluationGroup) getStudentsByClass(c *gin.Context) {
	if !eg.checkIfProfesor(c) {
		return
	}

	class := c.Param("class")
	elevi, err := eg.database.GetStudentsByClass(class)
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

	c.JSON(
		http.StatusOK,
		elrondApiShared.GenericAPIResponse{
			Data:  elevi,
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)
}

// UpdateFacade will update the facade
func (eg *evaluationGroup) UpdateFacade(newFacade shared.FacadeHandler) error {
	if check.IfNil(newFacade) {
		return errors.ErrNilFacadeHandler
	}

	eg.mutFacade.Lock()
	eg.facade = newFacade
	eg.mutFacade.Unlock()

	return nil
}

func (eg *evaluationGroup) checkIfProfesor(c *gin.Context) bool {
	email := c.GetString(authentication.EmailKey)

	isProfesor, err := eg.database.IsProfesor(email)
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

	if !isProfesor {
		c.JSON(
			http.StatusForbidden,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: "Nu esti un profesor",
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return false
	}

	return true
}

// IsAuthenticationNeeded will return true if the group requires authentication
func (eg *evaluationGroup) IsAuthenticationNeeded() bool {
	return eg.authenticationNeeded
}

// IsInterfaceNil returns true if there is no value under the interface
func (eg *evaluationGroup) IsInterfaceNil() bool {
	return eg == nil
}
