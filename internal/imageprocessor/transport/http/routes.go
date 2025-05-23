package http

import "github.com/gin-gonic/gin"

func (h *Handler) InitialRouteSetup() {
	h.Router.Use(gin.Recovery())
	h.Router.MaxMultipartMemory = 128 << 20
	h.Router.Use(LoggingMiddleware())
}

func (h *Handler) SetupRoutes() {

	h.Router.POST("/generate-preview", h.StartProcessingHandler, LoggingExitMiddleware())

}
