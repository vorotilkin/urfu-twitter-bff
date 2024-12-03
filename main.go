package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"twitter-bff/api"
	"twitter-bff/domain/services"
	"twitter-bff/infrastructure/posts"
	"twitter-bff/infrastructure/users"
	"twitter-bff/pkg/configuration"
	"twitter-bff/pkg/grpc"
	"twitter-bff/pkg/http"
	"twitter-bff/usecases"
)

type config struct {
	Http struct {
		Server http.Config
	}
	Grpc struct {
		Client struct {
			Users struct {
				Address string
			}
			Posts struct {
				Address string
			}
		}
	}
	Jwt struct {
		Secret string
	}
}

func newConfig(configuration *configuration.Configuration) (*config, error) {
	c := new(config)
	err := configuration.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	opts := []fx.Option{
		fx.Provide(zap.NewProduction),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(configuration.New),
		fx.Provide(http.NewServer),
		fx.Provide(newConfig),
		fx.Provide(validator.New),
		fx.Provide(func(c *config) http.Config {
			return http.Config{
				Addr:      c.Http.Server.Addr,
				SecretKey: c.Jwt.Secret,
			}
		}),
		fx.Provide(func(c *config) services.Config {
			return services.Config{
				SecretKey: c.Jwt.Secret,
			}
		}),
		fx.Provide(fx.Annotate(func(c *config) grpc.Config { return grpc.Config{Address: c.Grpc.Client.Users.Address} },
			fx.ResultTags(`name:"usersConfig"`))),
		fx.Provide(fx.Annotate(func(c *config) grpc.Config {
			return grpc.Config{Address: c.Grpc.Client.Posts.Address}
		},
			fx.ResultTags(`name:"postsConfig"`))),
		fx.Provide(fx.Annotate(grpc.NewClient, fx.ParamTags(`name:"usersConfig"`), fx.ResultTags(`name:"usersProvider"`))),
		fx.Provide(fx.Annotate(grpc.NewClient, fx.ParamTags(`name:"postsConfig"`), fx.ResultTags(`name:"postsProvider"`))),
		fx.Provide(fx.Annotate(
			users.NewRepository,
			fx.ParamTags(`name:"usersProvider"`),
			fx.As(new(services.LoginRepository)),
			fx.As(new(services.CreateRepository)),
			fx.As(new(services.UserByIDRepository)),
			fx.As(new(services.UpdateUserByIDRepository)),
		)),
		fx.Provide(fx.Annotate(
			posts.NewRepository,
			fx.ParamTags(`name:"postsProvider"`),
			fx.As(new(services.PostsRepository)),
		)),
		fx.Provide(services.NewCreateUserService),
		fx.Provide(services.NewLoginService),
		fx.Provide(services.NewUserByIDService),
		fx.Provide(services.NewUpdateUserByIDService),
		fx.Provide(services.NewPostsService),
		fx.Provide(usecases.NewEchoServer),
		fx.Invoke(func(lc fx.Lifecycle, server *http.Server) {
			lc.Append(fx.Hook{
				OnStart: server.OnStart,
				OnStop:  server.OnStop,
			})
		}),
		fx.Invoke(fx.Annotate(func(lc fx.Lifecycle, client *grpc.Client) {
			lc.Append(fx.Hook{
				OnStart: client.OnStart,
				OnStop:  client.OnStop,
			})
		}, fx.ParamTags("", `name:"usersProvider"`))),
		fx.Invoke(fx.Annotate(func(lc fx.Lifecycle, client *grpc.Client) {
			lc.Append(fx.Hook{
				OnStart: client.OnStart,
				OnStop:  client.OnStop,
			})
		}, fx.ParamTags("", `name:"postsProvider"`))),
		fx.Invoke(api.Registry),
	}

	app := fx.New(opts...)
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	<-app.Done()

	err = app.Stop(context.Background())
	if err != nil {
		panic(err)
	}
}
