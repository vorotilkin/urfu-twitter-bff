package usecases

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"net/http"
	"strconv"
	"time"
	"twitter-bff/domain/models"
	"twitter-bff/domain/services"
	"twitter-bff/openapigen"
	"twitter-bff/usecases/decorators"
)

type EchoServer struct {
	createSvc         *services.CreateUserService
	loginSvc          *services.LoginService
	userByIDService   *services.UserByIDService
	updateByIDService *services.UpdateUserByIDService
	postSvc           *services.PostsService
	followSvc         *services.FollowService
	likeSvc           *services.LikeService
}

func (s *EchoServer) Dislike(echoCtx echo.Context, postID int32) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	_, err = s.likeSvc.Like(context.Background(), jUser.UserID, postID, models.Dislike)
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusNoContent, nil)
}

func (s *EchoServer) Like(echoCtx echo.Context, postID int32) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	_, err = s.likeSvc.Like(context.Background(), jUser.UserID, postID, models.Like)
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusCreated, nil)
}

func (s *EchoServer) Unfollow(echoCtx echo.Context) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	var req openapigen.UnfollowJSONBody

	err = echoCtx.Bind(&req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = echoCtx.Validate(&req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	userID, err := strconv.ParseInt(lo.FromPtr(req.UserId), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	user, err := s.followSvc.Unfollow(context.Background(), jUser.UserID, int32(userID))
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoUser(user))
}

func (s *EchoServer) Follow(echoCtx echo.Context) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	var req openapigen.FollowJSONBody

	err = echoCtx.Bind(&req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = echoCtx.Validate(&req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	userID, err := strconv.ParseInt(lo.FromPtr(req.UserId), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	user, err := s.followSvc.Follow(context.Background(), jUser.UserID, int32(userID))
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoUser(user))
}

func (s *EchoServer) Comments(echoCtx echo.Context, params openapigen.CommentsParams) error {
	comments, err := s.postSvc.CommentsByPostID(context.Background(), lo.FromPtr(params.PostId))
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoComments(comments))
}

func (s *EchoServer) PostById(echoCtx echo.Context, id int32) error {
	jUser, _ := checkAuth(echoCtx)
	post, err := s.postSvc.PostByID(context.Background(), id, jUser.UserID)
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoPost(post))
}

func (s *EchoServer) Posts(echoCtx echo.Context, queryParams openapigen.PostsParams) error {
	ctx := context.Background()

	var (
		posts []models.Post
		err   error
	)

	if queryParams.UserId != nil {
		posts, err = s.postSvc.PostsByUserID(ctx, *queryParams.UserId)
		if err != nil {
			return echoCtx.JSON(ErrorHandler(err))
		}
	}

	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusOK, decorators.EchoPosts(posts))
	}

	posts, err = s.postSvc.FeedPosts(ctx, jUser.UserID)
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoPosts(posts))
}

func (s *EchoServer) CreatePost(echoCtx echo.Context) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	var req openapigen.CreatePostJSONBody

	err = echoCtx.Bind(&req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = echoCtx.Validate(&req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := context.Background()

	post, err := s.postSvc.Create(ctx, jUser.UserID, req.Body)
	if err != nil {
		if errors.Is(err, models.ErrInvalidArgument) {
			return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
		}
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusCreated, decorators.EchoPost(post))
}

func (s *EchoServer) UpdateUser(echoCtx echo.Context) error {
	jUser, err := checkAuth(echoCtx)
	if err != nil {
		return echoCtx.JSON(http.StatusUnauthorized, err.Error())
	}

	req := &openapigen.UserUpdateRequest{}

	err = echoCtx.Bind(req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = echoCtx.Validate(req)
	if err != nil {
		return echoCtx.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	userToUpdate := models.UserOption{
		ID:           jUser.UserID,
		Name:         mo.Some(req.Name),
		Username:     mo.Some(req.Username),
		Bio:          mo.PointerToOption(req.Bio),
		ProfileImage: mo.PointerToOption(req.ProfileImage),
		CoverImage:   mo.PointerToOption(req.CoverImage),
	}

	user, err := s.updateByIDService.UpdateUserByID(context.Background(), userToUpdate)
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoUser(user))
}

func (s *EchoServer) ListUsers(echoCtx echo.Context) error {
	users, err := s.userByIDService.NewUsers(context.Background())
	if err != nil {
		return echoCtx.JSON(ErrorHandler(err))
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoUsers(users))
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

	user, err := s.userByIDService.UserByID(context.Background(), jUser.UserID)
	if errors.Is(err, models.ErrNotFound) {
		return echoCtx.JSON(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoUser(user))
}

func (s *EchoServer) GetUser(echoCtx echo.Context, id int32) error {
	user, err := s.userByIDService.UserByID(context.Background(), id)
	if errors.Is(err, models.ErrNotFound) {
		return echoCtx.JSON(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, decorators.EchoUser(user))
}

func (s *EchoServer) Login(echoCtx echo.Context) error {
	req := &openapigen.LoginJSONBody{}

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

	return echoCtx.JSON(http.StatusOK, openapigen.JWTResponse{
		AccessToken: token.Token,
	})
}

func (s *EchoServer) CreateUser(echoCtx echo.Context) error {
	req := &openapigen.UserCreateRequest{}

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

	return echoCtx.JSON(http.StatusCreated, decorators.EchoUser(user))
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
	updateUserSvc *services.UpdateUserByIDService,
	postSvc *services.PostsService,
	followSvc *services.FollowService,
	likeSvc *services.LikeService,
) *EchoServer {
	return &EchoServer{
		createSvc:         createSvc,
		loginSvc:          loginSvc,
		userByIDService:   currentUserSvc,
		updateByIDService: updateUserSvc,
		postSvc:           postSvc,
		followSvc:         followSvc,
		likeSvc:           likeSvc,
	}
}
