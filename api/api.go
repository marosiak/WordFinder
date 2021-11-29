package api

import (
	"errors"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
	"time"
)

type API interface {
	Register(router *fasthttprouter.Router) error
}

type BaseAPI struct {
	listen  string
	apiList []API
	router  *fasthttprouter.Router
	server  *fasthttp.Server
}

func NewAPI(listen string, apiList ...API) (*BaseAPI, error) {
	a := &BaseAPI{
		listen: listen,
		router: fasthttprouter.New(),
	}

	for _, api := range apiList {
		if err := api.Register(a.router); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *BaseAPI) Listen() error {
	if a.server != nil {
		return errors.New("server already listen")
	}

	srv := &fasthttp.Server{
		Handler:     a.router.Handler,
		ReadTimeout: 1 * time.Minute,
		IdleTimeout: 5 * time.Second,
	}
	a.server = srv

	if strings.HasPrefix(a.listen, "unix:") {
		return srv.ListenAndServeUNIX(a.listen[5:], os.ModeSocket|os.ModePerm)
	}
	return srv.ListenAndServe(a.listen)
}
