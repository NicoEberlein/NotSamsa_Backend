package http

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type Handler struct {
	Router *gin.Engine
	S3     *minio.Client
}

func NewHandler(router *gin.Engine, s3 *minio.Client) *Handler {
	return &Handler{
		Router: router,
		S3:     s3,
	}
}
