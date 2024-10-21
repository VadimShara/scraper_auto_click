package service

import (
	"context"
	"net/http"
	"reflect"
	"scraper-first/pkg/logging"
	"strings"

	"github.com/gin-gonic/gin"
)

type handler struct {
	repository Repository
	logger *logging.Logger
}
 
func NewHandler(repository Repository, logger *logging.Logger) *handler {
	return &handler{
		repository: repository,
		logger : logger,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/search", h.GetList)
}

func (h *handler) GetList(c *gin.Context) {
	var e CreateEntity 
	var res *[]Entity
	val := reflect.ValueOf(&e).Elem()

	for i := 0; i < val.NumField(); i++ {
		if c.Request.URL.Query().Get(strings.ToLower(val.Type().Field(i).Name)) != "" {
			val.Field(i).SetString(c.Request.URL.Query().Get(strings.ToLower(val.Type().Field(i).Name)))
		}
	}

	res, err := h.repository.Find(context.Background(), &e)
	if err != nil {
		h.logger.Error(err)
	}

	for _, el := range *res {
		c.JSON(http.StatusOK, el)
	}
}