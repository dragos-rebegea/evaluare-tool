package groups

import (
	"encoding/json"
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
		{
			Path:    "/getAllClasses",
			Method:  http.MethodGet,
			Handler: eg.getAllClasses,
		},
		{
			Path:    "/addCalificativ",
			Method:  http.MethodPost,
			Handler: eg.addCalificativ,
		},
		{
			Path:    "/updateCalificativ",
			Method:  http.MethodPost,
			Handler: eg.updateCalificativ,
		},
		{
			Path:    "/getCalificative/:student",
			Method:  http.MethodGet,
			Handler: eg.getCalificative,
		},
		{
			Path:    "/getExercitii/:student",
			Method:  http.MethodGet,
			Handler: eg.getExercitii,
		},
		{
			Path:    "/ping",
			Method:  http.MethodGet,
			Handler: eg.ping,
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

func (eg *evaluationGroup) getAllClasses(c *gin.Context) {
	if !eg.checkIfProfesor(c) {
		return
	}

	email := c.GetString(authentication.EmailKey)
	classes, err := eg.database.GetAllClasses(email)
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
			Data:  classes,
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)
}

func (eg *evaluationGroup) addCalificativ(c *gin.Context) {
	if !eg.checkIfProfesor(c) {
		return
	}

	var calificativ core.Calificativ
	err := json.NewDecoder(c.Request.Body).Decode(&calificativ)
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
	email := c.GetString(authentication.EmailKey)
	err = eg.database.AddCalificativ(email, &calificativ)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			elrondApiShared.GenericAPIResponse{
				Data:  calificativ,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		elrondApiShared.GenericAPIResponse{
			Data:  calificativ,
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)
}

// updateCalificativ
func (eg *evaluationGroup) updateCalificativ(c *gin.Context) {
	if !eg.checkIfProfesor(c) {
		return
	}

	var calificativ core.Calificativ
	err := json.NewDecoder(c.Request.Body).Decode(&calificativ)
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
	email := c.GetString(authentication.EmailKey)
	err = eg.database.UpdateCalificativ(email, &calificativ)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			elrondApiShared.GenericAPIResponse{
				Data:  calificativ,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		elrondApiShared.GenericAPIResponse{
			Data:  calificativ,
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)
}

func (eg *evaluationGroup) getExercitii(context *gin.Context) {
	if !eg.checkIfProfesor(context) {
		return
	}

	email := context.GetString(authentication.EmailKey)
	student := context.Param("student")
	exercitii, err := eg.database.GetExercitiiForProfesorAndStudent(email, student)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		elrondApiShared.GenericAPIResponse{
			Data:  exercitii,
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)

}

func (eg *evaluationGroup) getCalificative(context *gin.Context) {
	if !eg.checkIfProfesor(context) {
		return
	}

	email := context.GetString(authentication.EmailKey)
	student := context.Param("student")
	calificative, err := eg.database.GetCalificative(email, student)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			elrondApiShared.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  elrondApiShared.ReturnCodeInternalError,
			},
		)
		return
	}

	context.JSON(
		http.StatusOK,
		elrondApiShared.GenericAPIResponse{
			Data:  calificative,
			Error: "",
			Code:  elrondApiShared.ReturnCodeSuccess,
		},
	)

}

func (eg *evaluationGroup) ping(c *gin.Context) {
	if !eg.checkIfProfesor(c) {
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
