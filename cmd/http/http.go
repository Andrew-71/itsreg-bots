package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	httpapi "github.com/bmstu-itstech/itsreg-bots/internal/api/http"
	"github.com/bmstu-itstech/itsreg-bots/internal/service"
	"github.com/bmstu-itstech/itsreg-bots/pkg/server"
)

func main() {
	app, _, _, closeFunc := service.NewApplication()
	defer func() {
		err := closeFunc()
		if err != nil {
			panic(err)
		}
	}()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return httpapi.HandlerFromMux(httpapi.NewHTTPServer(app), router)
	})
}
