package main

import (
	"github.com/NicoEberlein/NotSamsa_Backend/internal/imageprocessor/transport/http"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/storage/s3"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
)

type GlobalConfig struct {
	S3      *minio.Client
	Router  *gin.Engine
	Handler *http.Handler
}

var Config GlobalConfig

func main() {
	sss := s3.ConnectToS3()

	Config = GlobalConfig{
		S3: sss,
	}

	Config.Router = gin.Default()
	Config.Handler = http.NewHandler(Config.Router, Config.S3)

	Config.Handler.InitialRouteSetup()
	Config.Handler.SetupRoutes()

	if err := Config.Router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
