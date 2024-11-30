package api

import (
	"twitter-bff/pkg/http"
	"twitter-bff/server"
	"twitter-bff/usecases"
)

func Registry(provider *http.Server, serverImpl *usecases.EchoServer) {
	server.RegisterHandlersWithBaseURL(provider.Echo(), serverImpl, "/api")
}
