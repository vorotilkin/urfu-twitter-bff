package usecases

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"net/http"
	"twitter-bff/domain/services"
	"twitter-bff/server"
)

type EchoServer struct {
	createSvc *services.CreateUserService
	loginSvc  *services.LoginService
}

func (s *EchoServer) GetUser(echoCtx echo.Context, id string) error {
	token, ok := echoCtx.Get("user").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}
	return echoCtx.JSON(http.StatusOK, claims)
}

func (s *EchoServer) Login(echoCtx echo.Context) error {
	req := &server.LoginJSONBody{}

	if err := echoCtx.Bind(req); err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if req.Email == nil || req.Password == nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, "email or password is empty")
	}

	token, err := s.loginSvc.Login(context.Background(), string(lo.FromPtr(req.Email)), lo.FromPtr(req.Password))
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, errors.Wrap(err, "cant create token"))
	}

	return echoCtx.JSON(http.StatusOK, server.JWTResponse{
		AccessToken: token.Token,
	})
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

func NewEchoServer(createSvc *services.CreateUserService, loginSvc *services.LoginService) *EchoServer {
	return &EchoServer{createSvc: createSvc, loginSvc: loginSvc}
}
