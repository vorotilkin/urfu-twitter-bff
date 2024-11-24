package usecases

import (
	"context"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
	"net/http"
	"twitter-bff/domain/services"
	"twitter-bff/server"
)

type EchoServer struct {
	createSvc *services.CreateUserService
}

func NewEchoServer(createSvc *services.CreateUserService) *EchoServer {
	return &EchoServer{createSvc: createSvc}
}

func (s *EchoServer) CreateUser(echoCtx echo.Context) error {
	req := &server.UserCreateRequest{}

	err := echoCtx.Bind(req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = echoCtx.Validate(req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := context.Background()

	user, err := s.createSvc.Create(ctx, req.Name, req.Password, req.Name, string(req.Email))
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	response := &server.UserResponse{
		Email:    lo.ToPtr(openapi_types.Email(user.Email)),
		Id:       lo.ToPtr(user.ID),
		Name:     lo.ToPtr(user.Name),
		Username: lo.ToPtr(user.Username),
	}

	return echoCtx.JSON(http.StatusCreated, response)
}
