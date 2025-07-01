package web

import (
	"github.com/gin-gonic/gin"
)

type WebHandlers struct{}

func NewWebHandlers() *WebHandlers {
	return &WebHandlers{}
}

func (wh *WebHandlers) HandleHome(c *gin.Context) {
	c.File("./static/index.html")
}
