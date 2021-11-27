package api

import (
	"encoding/json"
	_ "github.com/fasthttp/router"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type New struct {
	Data  interface{} `json:"data"`
	Error *string     `json:"error"`
}

func WriteError(ctx *fasthttp.RequestCtx, error ErrorResponse) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(error.StatusCode)

	by, err := json.Marshal(New{Error: &error.Name})
	if err != nil {
		log.Error(err)
	}

	_, err = ctx.Write(by)
	if err != nil {
		log.Error(err)
	}
}

func WriteJSON(ctx *fasthttp.RequestCtx, code int, object New) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(code)

	by, err := json.Marshal(object)
	if err != nil {
		WriteError(ctx, ErrorByName("internal_error"))
	}

	ctx.Write(by)
}
