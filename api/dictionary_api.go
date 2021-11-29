package api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/marosiak/WordFinder/config"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type DictionaryAPI interface {
	GetDictsList(ctx *fasthttp.RequestCtx)
}

var _ API = &InternalDictionaryAPI{}

type InternalDictionaryAPI struct {
	cfg    *config.Config
	logger *log.Entry
}

func (s *InternalDictionaryAPI) Register(r *fasthttprouter.Router) error {
	r.GET("/dictionaries", s.GetDictsList)
	r.GET("/dicts", s.GetDictsList)
	return nil
}

func NewDictionaryAPI(cfg *config.Config, logger *log.Entry) *InternalDictionaryAPI {
	return &InternalDictionaryAPI{cfg: cfg, logger: logger}
}

func (s *InternalDictionaryAPI) GetDictsList(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("TODO")
}
