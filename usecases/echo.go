package usecases

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"time"
	"twitter-bff/domain/models"
	"twitter-bff/domain/services"
	"twitter-bff/server"
)

type EchoServer struct {
	createSvc      *services.CreateUserService
	loginSvc       *services.LoginService
	currentUserSvc *services.UserByIDService
}

func (s *EchoServer) Logout(echoCtx echo.Context) error {
	// Создаём cookie с пустым значением и временем истечения в прошлом
	cookie := &http.Cookie{
		Name:     models.JWTCookieName, // Имя куки, где хранится JWT токен
		Value:    "",                   // Очищаем значение
		Expires:  time.Unix(0, 0),      // Устанавливаем время истечения в прошлом
		HttpOnly: true,                 // Сохраняем HttpOnly, чтобы обезопасить куки
		Secure:   false,                // Убедитесь, что Secure соответствует вашим настройкам (true для HTTPS)
	}

	echoCtx.SetCookie(cookie)

	return echoCtx.JSON(http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

func (s *EchoServer) GetCurrentUser(echoCtx echo.Context) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	user, err := s.currentUserSvc.UserByID(context.Background(), jUser.UserID)
	if errors.Is(err, models.ErrNotFound) {
		return echoCtx.JSON(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, user)
}

func (s *EchoServer) GetUser(echoCtx echo.Context, id string) error {
	token, ok := echoCtx.Get("user").(jwt.MapClaims) // by default token is stored under `user` key
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}

	return echoCtx.JSON(http.StatusOK, token)
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
		return echoCtx.JSON(http.StatusUnprocessableEntity, errors.Wrap(err, "cant create token").Error())
	}

	cookie := new(http.Cookie)
	cookie.HttpOnly = true
	cookie.Name = models.JWTCookieName
	cookie.Value = token.Token

	echoCtx.SetCookie(cookie)

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

func jwtUser(echoCtx echo.Context) (models.JWTUser, error) {
	mapClaims, ok := echoCtx.Get("user").(jwt.MapClaims) // by default token is stored under `user` key
	if !ok {
		return models.JWTUser{}, errors.New("failed to cast claims as jwt.MapClaims")
	}

	expiredAt, err := mapClaims.GetExpirationTime()
	if err != nil {
		return models.JWTUser{}, errors.Wrap(err, "cant get token expiration time")
	}

	sub, err := mapClaims.GetSubject()
	if err != nil {
		return models.JWTUser{}, errors.Wrap(err, "cant get token subject")
	}

	userID, err := strconv.ParseInt(sub, 10, 32)
	if err != nil {
		return models.JWTUser{}, errors.Wrap(err, "cant parse user id")
	}

	return models.JWTUser{
		UserID:    int32(userID),
		ExpiredAt: expiredAt.Time,
	}, nil
}

func checkAuth(echoCtx echo.Context) (models.JWTUser, error) {
	jUser, err := jwtUser(echoCtx)
	if err != nil {
		return models.JWTUser{}, err
	}

	return jUser, nil
}

func NewEchoServer(
	createSvc *services.CreateUserService,
	loginSvc *services.LoginService,
	currentUserSvc *services.UserByIDService,
) *EchoServer {
	return &EchoServer{
		createSvc:      createSvc,
		loginSvc:       loginSvc,
		currentUserSvc: currentUserSvc,
	}
}
