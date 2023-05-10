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
			Path:    "/createClass",
			Method:  http.MethodPost,
			Handler: ag.createClass,
		},
		{
			Path:    "/createProfesor",
			Method:  http.MethodPost,
			Handler: ag.registerProfesor,
		},
		{
			Path:    "/setAbsent",
			Method:  http.MethodPost,
			Handler: ag.setAbsent,
		},
		{
			Path:    "/delStudent",
			Method:  http.MethodPost,
			Handler: ag.delStudent,
		},
		{
			Path:    "/createExam",
			Method:  http.MethodPost,
			Handler: ag.createExam,
		},
	}
	ag.endpoints = endpoints

	return ag, nil
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

func (ag *adminGroup) registerProfesor(context *gin.Context) {
	if !ag.checkIfAdmin(context) {
		return
	}
	var prof authentication.Profesor
	if err := context.ShouldBindJSON(&prof); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	prof.IsAdmin = false
	err := ag.database.CreateProfesor(&prof)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"userId":   prof.ID,
		"email":    prof.Email,
		"username": prof.Username,
		"password": prof.Password,
	})
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

	students, err := ag.database.CreateClass(&class)
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
	for _, student := range students {
		c.JSON(http.StatusCreated, gin.H{
			"userId":   student.ID,
			"email":    student.Email,
			"username": student.Username,
			"password": student.Password,
		})
	}
}

func (ag *adminGroup) setAbsent(c *gin.Context) {
	if !ag.checkIfAdmin(c) {
		return
	}

	var mark core.AbsentStatus
	err := json.NewDecoder(c.Request.Body).Decode(&mark)
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

	err = ag.database.SetAbsent(&mark)
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

func (ag *adminGroup) delStudent(c *gin.Context) {
	if !ag.checkIfAdmin(c) {
		return
	}

	var student authentication.Student
	err := json.NewDecoder(c.Request.Body).Decode(&student)
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

	err = ag.database.DeleteStudent(&student.ID)
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

// createExam will create a new exam
func (ag *adminGroup) createExam(c *gin.Context) {
	if !ag.checkIfAdmin(c) {
		return
	}

	var exam core.Exam
	err := json.NewDecoder(c.Request.Body).Decode(&exam)
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

	err = ag.database.CreateExam(&exam)
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
