package api

import (
	"twitter-bff/openapigen"
	"twitter-bff/pkg/http"
	"twitter-bff/usecases"
)

func Registry(provider *http.Server, serverImpl *usecases.EchoServer) {
	openapigen.RegisterHandlersWithBaseURL(provider.Echo(), serverImpl, "/api")
}
